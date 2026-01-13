package main

import (
	"device/arm"
	"device/stm32"
	"machine"
	"runtime/volatile"
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

	button = machine.PB7
)

// Set to one from the interrupt to indicate the button was pressed.
var buttonWake volatile.Register8

func main() {
	// Configure button
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure pins
	configureLEDs()

	// Set A1-A12 as open drain (and importantly, skip SWDIO/SWCLK)
	stm32.GPIOA.OTYPER.Set(0b1000_0111_1111_1111)

	setClockSpeed()

	// Zero all LEDs.
	for i := 0; i < 12; i++ {
		setLEDs(i, 0, 0, 0)
	}

	index := 0 // 0..11, group of 3 LEDs that will be updated together
	frame := 0
	mode := initialMode
	framesPressed := 0
	previousMode := 0
	for {
		// Update 3 LEDs at a time, since that's convenient for the
		// RGB-to-bitplane conversion.
		led0 := animate(mode, index+0, frame)
		led1 := animate(mode, index+12, frame)
		led2 := animate(mode, index+24, frame)
		setLEDs(index, uint32(led0), uint32(led1), uint32(led2))

		// Bitbang the LEDs.
		updateLEDs()

		index++
		if index == 12 {
			index = 0
			frame++

			// Read the button every frame update.
			pressed := !button.Get() // low means pressed
			if framesPressed == 30 {
				turnOffAnimation(mode, frame)

				// Sleep until a button press.
				sleepUntilButtonPress()

				// To continue the startup animation, set the mode to "power
				// on".
				previousMode = mode
				mode = modePowerOn
				frame = 0
				framesPressed = -0x8000_0000 // don't switch to the next animation on button release
			}
			if mode == modePowerOn {
				if frame == numLEDs/2 {
					mode = previousMode
				}
			} else {
				if !pressed && framesPressed > 0 {
					// Move to the next mode.
					mode++
					if mode >= modeLast {
						// Last, so wrap around.
						mode = 0
					}

					// Clear LEDs before moving on to the next mode.
					for i := 0; i < 12; i++ {
						setLEDs(i, 0, 0, 0)
					}
				}
			}
			if pressed {
				framesPressed++
			} else {
				framesPressed = 0
			}
		}
	}
}

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
	// Range 2: around 262kHz (~55µA with all LEDs off)
	// Range 3: around 524kHz (~83µA with all LEDs off)
	stm32.RCC.SetICSCR_MSIRANGE(stm32.RCC_ICSCR_MSIRANGE_Range2)

	// Reduce PCLK2/PCLK1 clocks since we don't need those peripherals (GPIO is
	// directly connected to the CPU). This saves around ~1.2µA.
	stm32.RCC.SetCFGR_PPRE1(stm32.RCC_CFGR_PPRE1_Div16)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Divide SYSCLK, for testing.
	//stm32.RCC.SetCFGR_HPRE(stm32.RCC_CFGR_HPRE_Div512)
}

// Show an animation during shutdown. It keeps the current animation, but
// freezes it in place, and turns off LEDs in sequence.
func turnOffAnimation(mode, frame int) {
	for i := range 36 / 3 {
		for index := range 12 {
			// Shut down LEDs in groups of 3.
			// We have to calculate the animation for two reasons:
			//  1. The animation may change in brightness if we don't call
			//     animate() to slow it down the same way as the normal
			//     animation.
			//  2. The 3 LEDs we can update at a time are spread over the ring
			//     of LEDs, not sequentially next to each other.
			led0 := animate(mode, index+0, frame)
			if index+0 < i*3 {
				led0 = 0
			}
			led1 := animate(mode, index+12, frame)
			if index+12 < i*3 {
				led1 = 0
			}
			led2 := animate(mode, index+24, frame)
			if index+24 < i*3 {
				led2 = 0
			}
			setLEDs(index, uint32(led0), uint32(led1), uint32(led2))

			// Bitbang the LEDs.
			updateLEDs()
		}
	}
}

func sleepUntilButtonPress() {
	// Goal: stop mode, without RTC, with 1 GPIO pin enabled for button interrupt

	// Enable pin interrupt.
	button.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
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
	button.SetInterrupt(machine.PinFalling, nil)
}

// Wait for a bit before turning on.
func waitForPoweron() bool {
	for range 12 * 30 {
		pressed := !button.Get() // low means pressed
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
