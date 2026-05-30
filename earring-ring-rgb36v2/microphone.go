package main

// Microphone on the earring.
// Some power measurements:
//
//  |  22µA | increase clock speed
//  | 131µA | enabling the mic
//  |  26µA | enabling the ADC and reading at 800Hz or so
//  | 179µA | combined power consumption while using the mic

import (
	"device/stm32"
	"machine"
	"math/bits"
)

var (
	powerBuffer      [4]uint16
	powerBufferIndex uint8
	powerBufferSum   uint32
)

var powerNormalizer uint32 = 100 // ~silence will result in about 500

// Take current power, smooth it out a little, and auto-normalize it so it
// doesn't change too strongly and animations can be written for it somewhat
// reasonably.
func addPower(power uint16) {
	// Find a number that when multiplied with powerBufferSum will be
	// roughly in range 16384..32768.
	factor := 1 << max(bits.LeadingZeros16(power)-3, 2)

	// Slowly adjust the normalizer.
	if factor > int(powerNormalizer) && powerNormalizer < 4096 {
		powerNormalizer++
	}
	if factor < int(powerNormalizer) && powerNormalizer > 16 {
		powerNormalizer--
	}

	powerBufferSum -= uint32(powerBuffer[powerBufferIndex])
	powerBufferSum += uint32(power)
	powerBuffer[powerBufferIndex] = power
	powerBufferIndex = (powerBufferIndex + 1) % uint8(len(powerBuffer))
}

// Return the current (auto-normalized) volume. The returned value will on
// average fall in the range 16384..32768 but will also frequently go outside
// that range.
func currentVolume() uint32 {
	return powerBufferSum * powerNormalizer
}

//go:noinline
func enableMic() {
	// Increase system clock speed: animations with the microphone need slightly
	// higher framerate to look good.
	stm32.RCC.SetICSCR_MSIRANGE(stm32.RCC_ICSCR_MSIRANGE_Range3)

	// Power on microphone.
	machine.PB1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PB1.High()

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
}

//go:noinline
func disableMic() {
	// Disable power to the microphone.
	machine.PB1.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	// Shut down ADC.
	stm32.RCC.APB2RSTR.Set(stm32.RCC_APB2RSTR_ADCRST)
	stm32.RCC.APB2RSTR.Set(0)
	stm32.RCC.SetAPB2ENR_ADCEN(0)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Restore system clock speed.
	stm32.RCC.SetICSCR_MSIRANGE(defaultMSIRANGE)
}

var micSamples [12]uint16

func updateMic(index int) {
	// Read sample, and add it to the buffer.
	micSamples[index] = uint16(stm32.ADC.DR.Get())

	// Start the next conversion.
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)
}

// Estimated microphone DC offset, continuously updated.
var micDCOffset int = 1250 // ~1250 seems to be a good start on at least some of these

// Process the samples collected in the last frame.
func processSamples() int {
	offset := micDCOffset
	sampleSum := 0
	powerSum := 0
	for _, sample := range micSamples[:] {
		sampleSum += int(sample)
		sampleDiff := int(sample) - offset
		if sampleDiff < 0 {
			sampleDiff = -sampleDiff
		}
		powerSum += sampleDiff
	}
	micDCOffset = sampleSum / len(micSamples)
	return powerSum
}
