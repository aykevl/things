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
	modeRainbowTrace
	modeNoise
	modeFireRed
	modeFireGreen
	modeFireBlue
	modeFlagLGBT
	modeFlagTrans
	modeLast

	modeTest
)

func animate(mode, led, frame int) Color {
	switch mode {
	case modeOff:
		return 0
	case modeRainbowTrace:
		return rainbowTrace(led, frame)
	case modeNoise:
		return noise(led, frame)
	case modeFireRed:
		return fire(led, frame, color.RGBA{R: 255})
	case modeFireGreen:
		return fire(led, frame, color.RGBA{G: 255})
	case modeFireBlue:
		return fire(led, frame, color.RGBA{B: 255})
	case modeFlagLGBT:
		return showPalette(led, frame, &flagLGBT)
	case modeFlagTrans:
		return showPalette(led, frame, &flagTrans)
	case modeTest:
		return testPulse(led, frame)
	default:
		// bug
		return errorPattern(led, frame)
	}
}

func noise(led, frame int) Color {
	x := uint32(frame) << 5
	y := uint32(led)
	c := ledsgo.PartyColors.ColorAt(ledsgo.Noise2(x, uint32(y)<<5) * 2)
	return NewColor(c.R, c.G, c.B)
}

var traceIndex int

func updateTraceIndex(led, frame int) int {
	// The animations have two tracers.
	if led == 0 && frame%2 == 0 {
		traceIndex++
		if traceIndex >= numLEDs {
			traceIndex = 0
		}
	}
	traceIndex2 := traceIndex + numLEDs/2
	if traceIndex2 >= numLEDs {
		traceIndex2 -= numLEDs
	}
	return traceIndex2
}

func rainbowTrace(led, frame int) Color {
	traceIndex2 := updateTraceIndex(led, frame)

	const div = 4
	if led == traceIndex {
		// First tracer.
		c := ledsgo.Color{H: uint16(frame * 128), S: 255, V: 255}.Rainbow()
		leds[led].R += c.R / div
		leds[led].G += c.G / div
		leds[led].B += c.B / div
	} else if led == traceIndex2 {
		// Second tracer, offset 180° on the color wheel and on the actual LED ring.
		c := ledsgo.Color{H: uint16(frame*128) + 0x8000, S: 255, V: 255}.Rainbow()
		leds[led].R += c.R / div
		leds[led].G += c.G / div
		leds[led].B += c.B / div
	} else {
		// dim LED
		c := leds[led]
		r := uint8(uint16(c.R) * 242 / 256)
		g := uint8(uint16(c.G) * 242 / 256)
		b := uint8(uint16(c.B) * 242 / 256)
		leds[led] = color.RGBA{r, g, b, 0}
	}
	c := leds[led]
	return NewColor(c.R, c.G, c.B)
}

// Fire animation in various colors.
// It is essential that this function is inlined, otherwise the fireColor isn't
// const-propagated and the whole animation is just way too slow to be usable.
//
//go:inline
func fire(led, frame int, fireColor color.RGBA) Color {
	intensityIndex := indexFromBottom(led)
	noiseIndex := uint32(frame) - uint32(intensityIndex)
	if led > numLEDs/2 {
		// Use a different flame on the otherside of the earring.
		// Without this, it would simply mirror the flame on both sides.
		// The constant is just an arbitrary number to give it enough distance
		// from the left side.
		noiseIndex += 0x1234
	}

	// Calculate the amount of heat on this particular pixel.
	heat := uint8(ledsgo.Noise1(noiseIndex<<10) >> 8)
	cooling := uint8(intensityIndex * 11)
	if heat < uint8(cooling) {
		heat = 0
	} else {
		heat -= uint8(cooling)
	}

	// Turn it into a flame, based on the given palette.
	// Perhaps we could use an actual 0-255 (or 0-64) palette instead? That
	// might be faster.
	c := coloredFlame(heat, fireColor)
	return NewColor(c.R, c.G, c.B)
}

// Colored flame. Like a heat map, but the lowest temperatures are not fixed red
// but instead use the configured color.
//
//go:inline
func coloredFlame(index uint8, fireColor color.RGBA) color.RGBA {
	if index < 128 {
		// <color>
		c := ledsgo.ApplyAlpha(fireColor, index*2)
		c.A = 255
		return c
	}
	if index < 224 {
		// <color>-yellow
		c1 := ledsgo.ApplyAlpha(fireColor, 255-uint8(uint32(index-128)*8/3))
		c2 := ledsgo.ApplyAlpha(color.RGBA{255, 255, 0, 255}, uint8(uint32(index-128)*8/3))
		return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}

func indexFromBottom(index int) int {
	// Start at the top with 18, move along the right size down to 0 at the
	// bottom, and then resume counting upwards again:
	// 18, 17, ..., 1, 0, 1, 2, ..., 17, 18
	newIndex := 18 - index
	if index >= 18 {
		newIndex = index - 18
	}
	return newIndex
}

type Palette [numLEDs]Color

var (
	flagRed    = NewColor(0xff, 0x00, 0x00)
	flagOrange = NewColor(0xff, 0x55, 0x00)
	flagYellow = NewColor(0x88, 0xff, 0x00)
	flagGreen  = NewColor(0x00, 0xff, 0x00)
	flagBlue   = NewColor(0x00, 0x00, 0xff)
	flagPurple = NewColor(0x80, 0x00, 0x80)
	flagLGBT   = Palette{
		flagRed, flagRed, flagRed, flagRed, flagRed, flagRed,
		flagOrange, flagOrange, flagOrange, flagOrange, flagOrange, flagOrange,
		flagYellow, flagYellow, flagYellow, flagYellow, flagYellow, flagYellow,
		flagGreen, flagGreen, flagGreen, flagGreen, flagGreen, flagGreen,
		flagBlue, flagBlue, flagBlue, flagBlue, flagBlue, flagBlue,
		flagPurple, flagPurple, flagPurple, flagPurple, flagPurple, flagPurple,
	}

	flagPastelBlue  = NewColor(0x11, 0x33, 0x88)
	flagPastelPink  = NewColor(0x80, 0x22, 0x22)
	flagPastelWhite = NewColor(0x88, 0xaa, 0xaa)
	flagTrans       = Palette{
		flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelWhite, flagPastelWhite, flagPastelWhite, flagPastelWhite,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelWhite, flagPastelWhite, flagPastelWhite, flagPastelWhite,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelBlue, flagPastelBlue,
	}
)

func showPalette(led, frame int, palette *Palette) Color {
	return palette[led]
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
