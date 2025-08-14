package main

import (
	"github.com/aykevl/ledsgo"
)

const initialMode = 1 // first animation (0 is off)

const (
	modeOff = iota
	modeSpirals
	modeScanner
	modeNoise
	modeLast

	modeTestPulse
	modeTestRotate
	modeTestRGB
)

func animate(mode, led, frame int) Color {
	switch mode {
	case modeOff:
		return 0
	case modeNoise:
		return noise(led, frame)
	case modeSpirals:
		return spirals(led, frame)
	case modeScanner:
		return scanner(led, frame)
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

var multiply uint32 = 65536

func noise(led, frame int) Color {
	if led == 0 && frame%16 == 0 {
		multiply = multiply * 33 / 32
		println("multiply:", multiply)
	}
	x := uint32(frame) << 3
	x = 0
	y := uint32(led) * multiply
	c := ledsgo.PartyColors.ColorAt(ledsgo.Noise2(x, y) * 2)
	return NewColor(c.R, c.G, c.B)
}

func scanner(led, frame int) Color {
	// angle over 360 degrees: 137.508
	// angle over 256 "degrees": 97.783
	// angle over 1024 "degrees": 391.134
	pos := (1023 - uint(led*391+frame*4)%1024) / 4
	if pos <= 128 {
		pos = 0
	} else {
		pos = pos - 128
	}
	return NewColor(0, uint8(pos), 0)
}

func spirals(led, frame int) Color {
	// angle over 360 degrees: 137.508
	// angle over 256 "degrees": 97.783
	// angle over 1024 "degrees": 391.134
	// angle over 32768 "degrees": 12516.3
	variation := (int(ledsgo.Noise1(uint32(frame)<<2)) - 32768) / 64
	pos := uint(led*(12516+variation)+frame*256) % 32768
	if pos >= 16384 {
		pos = 32767 - pos
	}
	pos = (pos * pos) >> 15
	c := ledsgo.PartyColors.ColorAt(ledsgo.Noise2(uint32(led)<<4, uint32(frame)))
	c = ledsgo.ApplyAlpha(c, uint8(pos>>7))
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
