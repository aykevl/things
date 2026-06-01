package main

// Microphone on the earring.
// Some power measurements:
//
//  |  22µA | increase clock speed
//  | 131µA | enabling the mic
//  |  26µA | enabling the ADC and reading at 800Hz or so
//  | 179µA | combined power consumption while using the mic

import (
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

// Estimated microphone DC offset, continuously updated.
var micDCOffset int

var micSamples [12]uint16

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

	// Slowly adjust the DC offset each update.
	newOffset := sampleSum / len(micSamples)
	if newOffset > micDCOffset {
		micDCOffset++
	} else {
		micDCOffset--
	}

	return powerSum
}
