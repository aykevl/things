package main

// Microphone on the earring.
// Some power measurements:
//
//  |  90µA | enabling the mic
//  |  28µA | enabling ADC with auto-off (not used)
//  |  34µA | enabling ADC with auto-off (read at 400Hz or so)
//  | 124µA | combined power consumption while using the mic

import (
	"device/stm32"
	"machine"
)

var (
	sampleBuffer      [64]uint16
	sampleBufferIndex uint8
	sampleBufferSum   uint32
)

func addSample(sample uint16) {
	sampleBufferSum -= uint32(sampleBuffer[sampleBufferIndex])
	sampleBufferSum += uint32(sample)
	sampleBuffer[sampleBufferIndex] = sample
	sampleBufferIndex = (sampleBufferIndex + 1) % uint8(len(sampleBuffer))
}

func avgSample() uint16 {
	return uint16(sampleBufferSum / uint32(len(sampleBuffer)))
}

var (
	powerBuffer      [16]uint16
	powerBufferIndex uint8
	powerBufferSum   uint32
)

func addPower(power uint16) {
	powerBufferSum -= uint32(powerBuffer[powerBufferIndex])
	powerBufferSum += uint32(power)
	powerBuffer[powerBufferIndex] = power
	powerBufferIndex = (powerBufferIndex + 1) % uint8(len(powerBuffer))
}

//go:noinline
func enableMic() {
	// Power on microphone.
	machine.PB6.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PB6.High()

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

	// Read initial samples to know the average output of the microphone.
	// The first few samples seem to be bad (microphone is still starting up?),
	// so read a bit more than what would be needed to initialize the sample
	// buffer.
	for range len(sampleBuffer) + 8 {
		// Start conversion.
		stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)

		// Wait until ready.
		for (stm32.ADC.ISR.Get() & stm32.ADC_ISR_EOC) == 0 {
		}

		// Read sample, and add it to the buffer.
		addSample(uint16(stm32.ADC.DR.Get()))
	}
}

//go:noinline
func disableMic() {
	// Disable power to the microphone.
	machine.PB6.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	// Shut down ADC.
	stm32.RCC.APB2RSTR.Set(stm32.RCC_APB2RSTR_ADCRST)
	stm32.RCC.SetAPB2ENR_ADCEN(0)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)
}

func updateMic() {
	// Read the most recent sample.
	sample := uint16(stm32.ADC.DR.Get())

	// Remove the DC offset from the sample.
	sampleDiff := int(sample) - int(avgSample())

	// Calculate the power, which is the absolute of the sample value.
	if sampleDiff < 0 {
		sampleDiff = -sampleDiff
	}

	// Add this sample to the moving average of recent samples (for DC offset
	// removal).
	addSample(sample)

	// Add this sample to the moving average of power values.
	addPower(uint16(sampleDiff))

	// Start the next conversion.
	stm32.ADC.CR.Set(stm32.ADC_CR_ADEN | stm32.ADC_CR_ADSTART)
}
