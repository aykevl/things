package main

import (
	"image/color"

	"github.com/aykevl/ledsgo"
)

const numLEDs = 36

// LED values, as a cache needed for some animations.
var leds [numLEDs]color.RGBA

const initialMode = 1 // first animation (0 is off)

const (
	modeOff = iota
	modeTestRGB
	modeNoise
	modeLast

	modeTest
)

func animate(mode, led, frame int) Color {
	switch mode {
	case modeOff:
		return 0
	case modeNoise:
		return noise(led, frame)
	case modeTest:
		return testPulse(led, frame)
	case modeTestRGB:
		return testRGB(led, frame)
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
	idx := int(frame / 8 % 36)
	value := uint8(0)
	if led == idx {
		value = 128
	}
	if (led+1)%36 == idx {
		value = 128 / 2
	}
	if (led+2)%36 == idx {
		value = 128 / 4
	}
	if (led+3)%36 == idx {
		value = 128 / 8
	}
	if (led+4)%36 == idx {
		value = 128 / 16
	}
	if (led+5)%36 == idx {
		value = 128 / 32
	}
	return NewColor(0, 0, value)
}

// Pulse red LEDs around once per second, for testing.
func testPulse(led, frame int) Color {
	return NewColor(uint8((frame%32)<<3), 0, 0)
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
