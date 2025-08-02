package main

import (
	"github.com/aykevl/ledsgo"
)

const initialMode = modeTestRotate // first animation (0 is off)

const (
	modeOff = iota
	modeTestRGB
	modeNoise
	modeLast

	modeTestPulse
	modeTestRotate
)

func animate(mode, led, frame int) Color {
	switch mode {
	case modeOff:
		return 0
	case modeNoise:
		return noise(led, frame)
	case modeTestPulse:
		return testPulse(led, frame)
	case modeTestRGB:
		return testRGB(led, frame)
	case modeTestRotate:
		return rotateSingleColor(led, frame)
	default:
		// bug
		return errorPattern(led, frame)
	}
}

func noise(led, frame int) Color {
	x := uint32(frame) << 3
	y := uint32(led)
	c := ledsgo.PartyColors.ColorAt(ledsgo.Noise2(x, uint32(y)<<5) * 2)
	return NewColor(c.R, c.G, c.B)
}

// Blink the first LED, roughly 0.5s on, 0.5s off.
func errorPattern(led, frame int) Color {
	if led == 0 {
		// Roughly 500ms on, 500ms off (assuming 30fps).
		return NewColor(uint8(frame%32)/16*128, 0, 0)
	}
	// Other LEDs are dark.
	return NewColor(0, 0, 0)
}

func rotateSingleColor(led, frame int) Color {
	idx := int(frame / 16 % 36)
	value := uint8(0)
	if led == idx {
		value = 128
	}
	if (led+1)%numLEDs == idx {
		value = 128 / 2
	}
	if (led+2)%numLEDs == idx {
		value = 128 / 4
	}
	if (led+3)%numLEDs == idx {
		value = 128 / 8
	}
	if (led+4)%numLEDs == idx {
		value = 128 / 16
	}
	if (led+5)%numLEDs == idx {
		value = 128 / 32
	}
	return NewColor(value, 0, 16)
}

// Pulse red LEDs around once per second, for testing.
func testPulse(led, frame int) Color {
	return NewColor(uint8((frame%64)<<2), 0, 0)
}

func testRGB(led, frame int) Color {
	switch (uint(frame) / 64) % 4 {
	case 0:
		return NewColor(255, 0, 0)
	case 1:
		return NewColor(0, 255, 0)
	case 2:
		return NewColor(0, 0, 255)
	default:
		return NewColor(64, 64, 64)
	}
}
