//go:build baremetal

package main

import (
	"device/arm"
	"device/stm32"
	"machine"
	"runtime/volatile"
	"unsafe"
)

const (
	A1  = machine.PA15
	A2  = machine.PA10
	A3  = machine.PA9
	A4  = machine.PA8
	A5  = machine.PA7
	A6  = machine.PA6
	A7  = machine.PA5
	A8  = machine.PA4
	A9  = machine.PA3
	A10 = machine.PA2
	A11 = machine.PA1
	A12 = machine.PA0

	button1 = machine.PB7
	button2 = machine.PB3

	micPin = machine.PB1
)

// Set to one from the interrupt to indicate the button was pressed.
var buttonWake volatile.Register8

func initHardware() {
	// Configure buttons
	button1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	button2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure pins
	configureLEDs()

	setClockSpeed()

	// Load last animation state from flash.
	loadState()
}

func button1Pressed() bool {
	return !button1.Get() // low means pressed
}

func button2Pressed() bool {
	return !button2.Get()
}

// Range 2: around 262kHz (~55µA with all LEDs off)
// Range 3: around 524kHz (~83µA with all LEDs off)
const defaultMSIRANGE = stm32.RCC_ICSCR_MSIRANGE_Range2

func setClockSpeed() {
	// Switch to MSI.
	stm32.RCC.CFGR.ReplaceBits(stm32.RCC_CFGR_SWS_MSI, stm32.RCC_CFGR_SW_Msk, 0)
	for stm32.RCC.CFGR.Get()&stm32.RCC_CFGR_SW_Msk != stm32.RCC_CFGR_SWS_MSI {
	}

	// Disable PLL.
	stm32.RCC.CR.ClearBits(stm32.RCC_CR_PLLON)

	// Configure power:
	// - change to 1.2V range.
	// - enable LPDS or LPSDSR bit
	// - disable Vrefint (ultra low power)
	// - enter standby mode when entering deepsleep (PDDS=1)
	// - clear wakeup flag
	stm32.RCC.APB1ENR.SetBits(stm32.RCC_APB1ENR_PWREN)
	stm32.PWR.CR.Set(stm32.PWR_CR_VOS_V1_2<<stm32.PWR_CR_VOS_Pos |
		stm32.PWR_CR_LPSDSR |
		stm32.PWR_CR_ULP |
		//stm32.PWR_CR_PDDS | // set this bit to enter standby mode (instead of stop mode)
		stm32.PWR_CR_CWUF |
		0)
	stm32.PWR.CR.Get() // make sure the values have been written
	stm32.RCC.APB1ENR.ClearBits(stm32.RCC_APB1ENR_PWREN)

	// Disable HSI16 clock.
	stm32.RCC.CR.ClearBits(stm32.RCC_CR_HSI16ON)

	// Change flash latency to zero wait states. Saves 0.4µA or so.
	stm32.FLASH.SetACR_LATENCY(stm32.Flash_ACR_LATENCY_WS0)

	// Disable flash during sleep.
	stm32.FLASH.ACR.Set(stm32.Flash_ACR_SLEEP_PD)

	// Disable TIM21. Doesn't use much current when the clock is disabled, but
	// it's a (very) small win. Saves 0.2µA or so.
	stm32.RCC.SetAPB2ENR_TIM21EN(0)

	// Set MSI clock speed.
	stm32.RCC.SetICSCR_MSIRANGE(defaultMSIRANGE)

	// Reduce PCLK2/PCLK1 clocks since we don't need those peripherals (GPIO is
	// directly connected to the CPU). This saves around ~1.2µA.
	stm32.RCC.SetCFGR_PPRE1(stm32.RCC_CFGR_PPRE1_Div16)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Divide SYSCLK, for testing.
	//stm32.RCC.SetCFGR_HPRE(stm32.RCC_CFGR_HPRE_Div512)
}

func sleepUntilButtonPress() {
	// Goal: stop mode, without RTC, with 1 GPIO pin enabled for button interrupt

	// Enable pin interrupt.
	button1.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		buttonWake.Set(1)
	})

	// Disable GPIO pins during sleep.
	disableLEDs()

	// Clear LEDs to avoid flash on poweron.
	for i := 0; i < 12; i++ {
		setLEDs(i, 0, 0, 0)
	}

	// Enter stop mode, wake up on a button press.
	buttonWake.Set(0)
	arm.SCB.SCR.SetBits(arm.SCB_SCR_SLEEPDEEP)
	for {
		arm.Asm("wfe")
		if buttonWake.Get() != 0 {
			// Wait a second and check that the button doesn't get released
			// during this time.
			if waitForPoweron() {
				break
			}
		}
	}
	arm.SCB.SCR.ClearBits(arm.SCB_SCR_SLEEPDEEP)

	// Disable interrupt, so it won't cause flickering when switching modes.
	button1.SetInterrupt(machine.PinFalling, nil)
}

// Wait for a bit before turning on.
func waitForPoweron() bool {
	for range 12 * 30 {
		pressed := !button1.Get() // low means pressed
		if !pressed {
			// Button was released early, shut down device again.
			return false
		}

		// Do this just to delay stuff a little.
		updateLEDs()
	}

	// Turn on LEDs again after power down.
	configureLEDs()

	// Indicate that the chip should start up again.
	return true
}

var micEnabled bool

//go:noinline
func enableMic() {
	if micEnabled {
		return
	}
	micEnabled = true

	// Increase system clock speed: animations with the microphone need slightly
	// higher framerate to look good.
	stm32.RCC.SetICSCR_MSIRANGE(stm32.RCC_ICSCR_MSIRANGE_Range3)

	// Power on microphone.
	micPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	micPin.High()

	// Provide clock to ADC.
	// Note that PCLK2 is divided when the mic is disabled to save power (~1µA),
	// but the ADC doesn't work with such a low frequency so we need to set the
	// division value back to 1 (no division).
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div1)
	stm32.RCC.SetAPB2ENR_ADCEN(1)

	// Use the PCLK clock instead of HSI16.
	stm32.ADC.CFGR2.Set(stm32.ADC_CFGR2_CKMODE_PCLK << stm32.ADC_CFGR2_CKMODE_Pos)

	// Set to low frequency mode (necessary below 3.5MHz).
	stm32.ADC.CCR.Set(stm32.ADC_CCR_LFMEN)

	// Power on ADC.
	// - clear ADRDY bit
	// - ADEN set to 1
	// - wait for ADRDY flag to be 1
	// - set ADSTART to 1
	stm32.ADC.ISR.Set(stm32.ADC_ISR_ADRDY)             // clear ADRDY bit
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN)                // enable ADC
	for stm32.ADC.ISR.Get()&stm32.ADC_ISR_ADRDY == 0 { // wait until it is enabled
	}

	// Set "auto off" mode to save power.
	// Apparently this needs to be done after powering on.
	stm32.ADC.CFGR1.Set(stm32.ADC_CFGR1_AUTOFF)

	// Select the channel.
	stm32.ADC.CHSELR.Set(1 << 8) // PB0=ADC_IN8

	// Start first continuous conversion.
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)

	// Calibrate the DC offset.
	sum := 0
	const calibNum = 1024
	for range calibNum {
		// Wait until ready.
		for (stm32.ADC.ISR.Get() & stm32.ADC_ISR_EOC) == 0 {
		}

		// Read sample, and add it to the buffer.
		sum += int(stm32.ADC.DR.Get())

		// Start the next conversion.
		stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)
	}
	micDCOffset = sum / calibNum
}

//go:noinline
func disableMic() {
	if !micEnabled {
		return
	}
	micEnabled = false

	// Disable power to the microphone.
	micPin.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	// Shut down ADC.
	stm32.RCC.APB2RSTR.Set(stm32.RCC_APB2RSTR_ADCRST)
	stm32.RCC.APB2RSTR.Set(0)
	stm32.RCC.SetAPB2ENR_ADCEN(0)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Restore system clock speed.
	stm32.RCC.SetICSCR_MSIRANGE(defaultMSIRANGE)
}

func updateMic(index int) {
	// Read sample, and add it to the buffer.
	micSamples[index] = uint16(stm32.ADC.DR.Get())

	// Start the next conversion.
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)
}

func loadState() {
	// Load previous settings from flash.
	numPages := uint32(machine.Flash.Size()) / uint32(machine.Flash.EraseBlockSize())
	lastPage := numPages - 1
	machine.Flash.ReadAt(storedState[:], int64(lastPage*uint32(machine.Flash.EraseBlockSize())))

	// Check whether they seem to be correct.
	calculatedHash := hash32(storedState[:len(storedState)-4])
	storedHash := 0 |
		uint32(storedState[len(storedState)-4])<<0 |
		uint32(storedState[len(storedState)-3])<<8 |
		uint32(storedState[len(storedState)-2])<<16 |
		uint32(storedState[len(storedState)-1])<<24
	if calculatedHash != storedHash {
		if storedHash == 0 {
			// First use: show QA pattern.
			storedState = [len(storedState)]uint8{stateOffsetMode: modeQA}
		} else {
			// Reset the state to initial.
			storedState = [len(storedState)]uint8{stateOffsetMode: initialMode}
		}
	}
}

func saveState() {
	// Calculate the hash for this state.
	hash := hash32(storedState[:len(storedState)-4])
	storedState[len(storedState)-4] = uint8(hash >> 0)
	storedState[len(storedState)-3] = uint8(hash >> 8)
	storedState[len(storedState)-2] = uint8(hash >> 16)
	storedState[len(storedState)-1] = uint8(hash >> 24)

	// Store the state in the last page.
	numPages := uint32(machine.Flash.Size()) / uint32(machine.Flash.EraseBlockSize())
	lastPage := numPages - 1
	machine.Flash.EraseBlocks(int64(lastPage), 1)
	machine.Flash.WriteAt(storedState[:], int64(lastPage*uint32(machine.Flash.EraseBlockSize())))
}

// Get FNV-1a hash of the given memory buffer.
//
// https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV-1a_hash
func hash32(buf []byte) uint32 {
	var result uint32 = 2166136261 // FNV offset basis
	for _, c := range buf {
		result ^= uint32(c) // XOR with byte
		result *= 16777619  // FNV prime
	}
	return result
}

func loadPattern(slot int) []byte {
	// We also store config info in flash, make sure to leave it alone.
	patternsFlashEnd := machine.FlashDataEnd() - uintptr(machine.Flash.EraseBlockSize())

	slotAddr := patternsFlashEnd - slotSize*uintptr(slot+1)
	animationInFlash := unsafe.Slice((*uint32)(unsafe.Pointer(slotAddr)), slotSize/4)
	copy(animationBuf[:], animationInFlash)

	fileSize := int(animationInFlash[0])
	if fileSize > min(maxProgramSize, animationBufSize) {
		fileSize = 0
	}

	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(animationBuf[1:]))), fileSize)
}

// Store a pattern in flash.
func storePattern(slot int) (ok bool) {
	// We also store config info in flash, make sure to leave it alone.
	patternsFlashEnd := uintptr(machine.Flash.Size()) - uintptr(machine.Flash.EraseBlockSize())

	// Find the start address (in machine.Flash) for this slot.
	slotAddr := patternsFlashEnd - slotSize*uintptr(slot+1)

	// Clear flash at the destination.
	blockStart := int(slotAddr) / int(machine.Flash.EraseBlockSize())
	numBlocks := slotSize / int(machine.Flash.EraseBlockSize())
	if blockStart < 0 {
		return false
	}
	err := machine.Flash.EraseBlocks(int64(blockStart), int64(numBlocks))
	if err != nil {
		return false
	}

	// Write the pattern to flash.
	data8Len := min(slotSize, len(animationBuf)*4)
	data8 := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(animationBuf[:]))), data8Len)
	_, err = machine.Flash.WriteAt(data8, int64(slotAddr))
	if err != nil {
		return false
	}

	return true
}
