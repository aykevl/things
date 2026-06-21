package main

import (
	"device/stm32"
	"machine"
	"unsafe"

	"codeberg.org/maaike328p/datatrans"
	"codeberg.org/maaike328p/datatrans/dsp"
)

const (
	// possible values for 4.2MHz: 12000 14000 15000 25000 28000 30000 35000 42000 50000 etc
	// possible values for 4.0MHz: 10000 16000 20000 25000 40000 50000
	// possible values for 8.0MHz: 10000 16000 20000 25000 32000 40000 50000
	// but since we need calibration anyway, really any value would be fine
	sampleRate = 10000

	// number of samples per window we do frequency analysis on
	windowSize = sampleRate / datatrans.SymbolsPerSecond
)

// Size of a program slot in flash memory.
const slotSize = 1536 // 1.5kB

// Maximum size of a program (smaller than slotSize because we also need to save
// the program size).
const maxProgramSize = slotSize - 4

// Size (in bytes) of the buffer used during reception.
// Should be bigger than the slot, so that the fountain code has some space to
// work with.
const animationBufSize = maxProgramSize + 1024

var animationBuf [animationBufSize / 4]uint32

//go:align 4
var adcData [windowSize * 2]uint16

// ADC data with DC offset removed.
// Two extra elements: the last one to serve as a special "stop" value and the
// one after that because the Goertzel pair code reads two values at once.
var adcDataNormalized [windowSize + 2]int32

func adcDataWriteIndex() uint32 {
	return uint32(len(adcData)) - stm32.DMA1.CH[0].NDTR.Get()
}

var decoder datatrans.Decoder

var adcWriteWindow uint8 = 0 // 0 or 1 to indicate current window

var dcOffset uint32

func adcDelaySamples(samples uint32) {
	if samples != 0 {
		stm32.DMA1.CH[0].CR.ClearBits(stm32.DMA_CH_CR_EN)
		stm32.DMA1.CH[0].CR.Get()
		for range samples {
			for stm32.TIM22.CNT.Get() != 0 {
			}
			for stm32.TIM22.CNT.Get() == 0 {
			}
		}

		stm32.DMA1.CH[0].CR.SetBits(stm32.DMA_CH_CR_EN)

		// This apparently helps to avoid a stuck DMA??
		// Though it doesn't seem to be needed anymore.
		//if stm32.DMA1.CH[0].CR.Get()&stm32.DMA_CH_CR_EN == 0 {
		//	for {
		//	}
		//}
	}
}

// Wait until the write pointer is in the next window, and update
// adcDataNormalized so it contains the data for the window that's ready.
func adcDataNextWindow(delayBefore uint32) {
	var rawSamples []uint16
	switch adcWriteWindow {
	case 0:
		if adcDataWriteIndex() >= windowSize {
			// Shouldn't be in the next window already!
			println("overflow!")
		}
		for adcDataWriteIndex() < windowSize {
			// Make sure it's in the next window.
		}
		adcDelaySamples(delayBefore)
		adcWriteWindow = 1 // next window
		rawSamples = adcData[:windowSize]
	default:
		if adcDataWriteIndex() < windowSize {
			// Shouldn't be in the first window already!
			println("overflow!")
		}
		for adcDataWriteIndex() >= windowSize {
			// Make sure it's in the first window.
		}
		adcDelaySamples(delayBefore)
		adcWriteWindow = 0 // next window
		rawSamples = adcData[windowSize : windowSize*2]
	}

	// Adjust samples by DC offset.
	// At the same time, also sum up all the samples it reads so we can
	// calculate the DC offset for the next window.
	dcOffset = dsp.CopyWithOffset(rawSamples, adcDataNormalized[:], dcOffset)
}

func dataRecv(slot int) {
	// Reset LEDs.
	disableLEDs()

	// Enable indicator lights.
	A10.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	A12.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	A9.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A9.High()
	if slot >= 1 {
		A8.Configure(machine.PinConfig{Mode: machine.PinOutput})
		A8.High()
	}
	if slot >= 2 {
		A7.Configure(machine.PinConfig{Mode: machine.PinOutput})
		A7.High()
	}

	// Enable HSI16.
	stm32.RCC.CR.SetBits(stm32.RCC_CR_HSI16ON)

	// Switch to HSI16/4 as system clock.
	stm32.RCC.SetCFGR_HPRE(stm32.RCC_CFGR_HPRE_Div4)
	stm32.RCC.SetCFGR_SW(stm32.RCC_CFGR_SW_HSI16)

	// Power on microphone.
	micPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	micPin.High()

	// Provide clock to ADC.
	// Note that PCLK2 is divided when the mic is disabled to save power (~1µA),
	// but the ADC doesn't work with such a low frequency so we need to set the
	// division value back to 1 (no division).
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div1)
	stm32.RCC.SetAPB2ENR_ADCEN(1)

	// Power on ADC.
	// - clear ADRDY bit
	// - ADEN set to 1
	// - wait for ADRDY flag to be 1
	// - set ADSTART to 1
	stm32.ADC.ISR.Set(stm32.ADC_ISR_ADRDY)             // clear ADRDY bit
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN)                // enable ADC
	for stm32.ADC.ISR.Get()&stm32.ADC_ISR_ADRDY == 0 { // wait until it is enabled
	}

	// Select the channel.
	stm32.ADC.CHSELR.Set(1 << 8) // PB0=ADC_IN8

	// Enable a timer, for use with the DMA.
	const systemClock = 4e6
	const clockPrescaler = uint32(systemClock/sampleRate/2 - 1)
	stm32.RCC.SetAPB2ENR_TIM22EN(1)
	stm32.TIM22.PSC.Set(clockPrescaler) // prescale so it counts at twice the samplerate
	stm32.TIM22.ARR.Set(1)              // wraparound every other clock cycle
	stm32.TIM22.CR1.Set(stm32.TIM_CR1_CEN)
	stm32.TIM22.CR2.Set(stm32.TIM_CR2_MMS_Update << stm32.TIM_CR2_MMS_Pos)

	// Debug: slow down the ADC
	//stm32.TIM22.PSC.Set(0xffff)
	//stm32.TIM22.ARR.Set(1)

	// Configure DMA.
	stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_DMAEN)
	stm32.DMA1.CH[0].PAR.Set(uint32(uintptr(unsafe.Pointer(&stm32.ADC.DR))))
	stm32.DMA1.CH[0].MAR.Set(uint32(uintptr(unsafe.Pointer(&adcData))))
	stm32.DMA1.CH[0].NDTR.Set(uint32(len(adcData)))
	stm32.DMA1.CH[0].CR.Set(stm32.DMA_CH_CR_MINC | stm32.DMA_CH_CR_MSIZE_Bits16<<stm32.DMA_CH_CR_MSIZE_Pos | stm32.DMA_CH_CR_PSIZE_Bits16<<stm32.DMA_CH_CR_PSIZE_Pos | stm32.DMA_CH_CR_TEIE | stm32.DMA_CH_CR_CIRC)
	stm32.DMA1.CH[0].CR.SetBits(stm32.DMA_CH_CR_EN)

	// Set the ADC to trigger from DMA.
	stm32.ADC.CFGR1.Set(stm32.ADC_CFGR1_DMAEN | stm32.ADC_CFGR1_DMACFG | stm32.ADC_CFGR1_EXTSEL_TIM22_TRGO<<stm32.ADC_CFGR1_EXTSEL_Pos | stm32.ADC_CFGR1_EXTEN_RisingEdge<<stm32.ADC_CFGR1_EXTEN_Pos)
	stm32.ADC.CFGR2.SetBits(0 |
		stm32.ADC_CFGR2_OVSE |
		1<<stm32.ADC_CFGR2_OVSS_Pos |
		stm32.ADC_CFGR2_OVSR_Mul32<<stm32.ADC_CFGR2_OVSR_Pos)

	// Start the ADC!
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)

	// Special marker used by the Goertzel implementation to know this is the
	// end of the array.
	adcDataNormalized[windowSize] = 0x7fff_ffff

	// Reset the trim value (it will be calibrated during reception).
	// This might need to be done a long time after enabling HSI16 to avoid a
	// HardFault? Not sure.
	stm32.RCC.SetICSCR_HSI16TRIM(16)

	// Receive the data.
	decoder.Initialize(animationBuf[1:], maxProgramSize, sampleRate)
	receiveData(slot)

	// Disable everything again.

	// Shut down DMA.
	stm32.RCC.AHBENR.ClearBits(stm32.RCC_AHBENR_DMAEN)

	// Shut down ADC.
	stm32.RCC.APB2RSTR.Set(stm32.RCC_APB2RSTR_ADCRST)
	stm32.RCC.APB2RSTR.Set(0)
	stm32.RCC.SetAPB2ENR_ADCEN(0)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Shut down timer.
	stm32.RCC.SetAPB2ENR_TIM22EN(0)

	// Disable power to the microphone.
	micPin.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	// Switch back to MSI.
	stm32.RCC.CFGR.ReplaceBits(stm32.RCC_CFGR_SWS_MSI, stm32.RCC_CFGR_SW_Msk, 0)
	for stm32.RCC.CFGR.Get()&stm32.RCC_CFGR_SW_Msk != stm32.RCC_CFGR_SWS_MSI {
	}
	stm32.RCC.SetCFGR_HPRE(stm32.RCC_CFGR_HPRE_Div1)

	// Disable HSI16
	stm32.RCC.CR.ClearBits(stm32.RCC_CR_HSI16ON)

	// Restore the LEDs.
	// Note: have to set the below two as input due to a bug (pull mode is not
	// reset when switching from input to output mode).
	A10.Configure(machine.PinConfig{Mode: machine.PinInput})
	A11.Configure(machine.PinConfig{Mode: machine.PinInput})
	A12.Configure(machine.PinConfig{Mode: machine.PinInput})
	configureLEDs()
}

func receiveData(slot int) {
	decoder.Reset()

	button2WasReleased := false // true once button 2 was released once (for exiting the receiver)
	whiteWasOn := false         // used for blinking LED

	// Wait for the first full symbol.
	adcDataNextWindow(0)
	for {
		result := decoder.ProcessSamples(adcDataNormalized[:])

		// Adjust clock frequency, to calibrate to the audio source.
		// This is also used to continuously adjust the frequency
		trim := stm32.RCC.GetICSCR_HSI16TRIM()
		if adj := decoder.SpeedAdjust(); adj < 0 {
			trim = max(13, trim-1)
			stm32.RCC.SetICSCR_HSI16TRIM(trim)
		} else if adj > 0 {
			trim = min(19, trim+1)
			stm32.RCC.SetICSCR_HSI16TRIM(trim)
		}

		// Decoder had to stop for some reason (error or fully received).
		if result < 0 {
			if result == datatrans.ErrCodeReceived {
				// Fully received! Show green LEDs (while saving the binary):
				A10.Configure(machine.PinConfig{Mode: machine.PinInput})
				A11.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
				A12.Configure(machine.PinConfig{Mode: machine.PinInput})

				// Save the binary.
				binary := decoder.Bytes()
				animationBuf[0] = uint32(len(binary))
				storePattern(slot)
				return
			}
			decoder.Reset()
			// Some error, restart.
			adcDataNextWindow(0)
			continue
		}

		if decoder.ReceivingValid() && !whiteWasOn {
			// Blink the LEDs white-ish.
			A11.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
			whiteWasOn = true
		} else {
			// Go back to the previous (purple) color.
			A11.Configure(machine.PinConfig{Mode: machine.PinInput})
			whiteWasOn = false
		}

		// Wait until the next symbol arrives.
		adcDataNextWindow(uint32(result))

		// Exit on button2 press (but don't exit immediately when button2 is
		// still pressed due to the long press to enter receiving mode).
		pressed := button2Pressed()
		if !pressed {
			button2WasReleased = true
		}
		if pressed && button2WasReleased {
			// stop receiving
			return
		}
	}
}
