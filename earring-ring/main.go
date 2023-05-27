package main

import (
	"time"

	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl/pixel"
)

const numLeds = 18

var leds [numLeds]pixel.LinearGRB888

var cycle uint16

var buttonPressed bool

func main() {
	initHardware()

	var traceIndex uint8
	animation := uint8(1)
	ledIndex := uint8(0)
	for {
		// Special handling for sleep mode.
		// We don't want to write to LEDs or anything like that, just let the
		// chip stay asleep.
		if animation == 0 {
			pressed := isButtonPressed()
			if pressed != buttonPressed {
				buttonPressed = pressed
				if pressed {
					// wake up, go to first animation
					animation = 1
					continue
				}
			}
			// This uses ~2µA sleep current.
			// This can be optimized further if needed using pin interrupts.
			time.Sleep(time.Second / 4)
			continue
		}

		updateLEDs()

		// Update 2 LEDs.
		for i := uint8(0); i < 2; i++ {
			const numAnimations = 6
			switch animation {
			case 1:
				rainbowTrace(ledIndex, traceIndex)
			case 2:
				fireAndIce(ledIndex, traceIndex)
			case 3:
				purpleCircles(ledIndex, traceIndex)
			case 4:
				sparkle(ledIndex)
			case 5:
				showPalette(ledIndex, &flagLGBT)
			case 6:
				showPalette(ledIndex, &flagTrans)
			}

			ledIndex++
			if ledIndex >= 18 {
				ledIndex = 0
				cycle++
				// Trace animation.
				if cycle%4 == 1 {
					traceIndex++
					if traceIndex >= 18 {
						traceIndex = 0
					}
				}
				// Respond to button presses.
				pressed := isButtonPressed()
				if pressed != buttonPressed {
					buttonPressed = pressed
					if pressed {
						animation++
						if animation > numAnimations {
							// Wrap around, and enter sleep mode.
							animation = 0
							for i := range leds {
								leds[i] = pixel.LinearGRB888{}
							}
							initHardware()
							continue
						}
					}
				}
			}
		}
	}
}

// Three purple tracers running in circles.
func purpleCircles(index, traceIndex uint8) {
	// This animation has three tracers.
	traceIndex2 := traceIndex + uint8(len(leds))/3
	if traceIndex2 >= uint8(len(leds)) {
		traceIndex2 -= uint8(len(leds))
	}
	traceIndex3 := traceIndex + uint8(len(leds))*2/3
	if traceIndex3 >= uint8(len(leds)) {
		traceIndex3 -= uint8(len(leds))
	}

	if index == traceIndex || index == traceIndex2 || index == traceIndex3 {
		// First tracer.
		leds[index].R += 16
		leds[index].B += 8
	} else {
		// dim LED
		c := any(leds[index]).(pixel.LinearGRB888)
		r := uint8(uint16(c.R) * 230 / 256)
		g := uint8(uint16(c.G) * 230 / 256)
		b := uint8(uint16(c.B) * 230 / 256)
		leds[index] = pixel.NewLinearGRB888(r, g, b)
	}
}

func rainbowTrace(index, traceIndex uint8) {
	// This animation has two tracers.
	traceIndex2 := traceIndex + uint8(len(leds))/2
	if traceIndex2 >= uint8(len(leds)) {
		traceIndex2 -= uint8(len(leds))
	}

	const div = 4
	if index == traceIndex {
		// First tracer.
		c1 := ledsgo.Color{H: cycle * 128, S: 255, V: 255}.Rainbow()
		leds[index].R += c1.R / div
		leds[index].G += c1.G / div
		leds[index].B += c1.B / div
	} else if index == traceIndex2 {
		// Second tracer, offset 180° on the color wheel and on the actual LED ring.
		c2 := ledsgo.Color{H: cycle*128 + 0x8000, S: 255, V: 255}.Rainbow()
		leds[index].R += c2.R / div
		leds[index].G += c2.G / div
		leds[index].B += c2.B / div
	} else {
		// dim LED
		c := any(leds[index]).(pixel.LinearGRB888)
		r := uint8(uint16(c.R) * 225 / 256)
		g := uint8(uint16(c.G) * 225 / 256)
		b := uint8(uint16(c.B) * 225 / 256)
		leds[index] = pixel.NewLinearGRB888(r, g, b)
	}
}

var (
	flagLGBT = [18]pixel.LinearGRB888{
		{R: 0xff / 3, G: 0x00 / 3, B: 0x00 / 3}, // red
		{R: 0xff / 3, G: 0x00 / 3, B: 0x00 / 3},
		{R: 0xff / 3, G: 0x00 / 3, B: 0x00 / 3},
		{R: 0xff / 3, G: 0x22 / 3, B: 0x00 / 3}, // orange
		{R: 0xff / 3, G: 0x22 / 3, B: 0x00 / 3},
		{R: 0xff / 3, G: 0x22 / 3, B: 0x00 / 3},
		{R: 0x88 / 3, G: 0xff / 3, B: 0x00 / 3}, // yellow
		{R: 0x88 / 3, G: 0xff / 3, B: 0x00 / 3},
		{R: 0x88 / 3, G: 0xff / 3, B: 0x00 / 3},
		{R: 0x00 / 3, G: 0xff / 3, B: 0x00 / 3}, // green
		{R: 0x00 / 3, G: 0xff / 3, B: 0x00 / 3},
		{R: 0x00 / 3, G: 0xff / 3, B: 0x00 / 3},
		{R: 0x00 / 3, G: 0x00 / 3, B: 0xff / 3}, // blue
		{R: 0x00 / 3, G: 0x00 / 3, B: 0xff / 3},
		{R: 0x00 / 3, G: 0x00 / 3, B: 0xff / 3},
		{R: 0x80 / 3, G: 0x00 / 3, B: 0x80 / 3}, // purple
		{R: 0x80 / 3, G: 0x00 / 3, B: 0x80 / 3},
		{R: 0x80 / 3, G: 0x00 / 3, B: 0x80 / 3},
	}

	flagTrans = [18]pixel.LinearGRB888{
		{R: 0x01, G: 0x11, B: 0x66},             // pastel blue
		{R: 0x01, G: 0x11, B: 0x66},             // pastel blue
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x33 / 3, G: 0xcc / 3, B: 0xff / 3}, // white
		{R: 0x33 / 3, G: 0xcc / 3, B: 0xff / 3}, // white
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x01, G: 0x11, B: 0x66},             // pastel blue
		{R: 0x01, G: 0x11, B: 0x66},             // pastel blue
		{R: 0x01, G: 0x11, B: 0x66},             // pastel blue
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x33 / 3, G: 0xcc / 3, B: 0xff / 3}, // white
		{R: 0x33 / 3, G: 0xcc / 3, B: 0xff / 3}, // white
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x66 / 3, G: 0x22 / 3, B: 0x80 / 3}, // pastel pink
		{R: 0x01, G: 0x11, B: 0x66},             // pastel blue
	}
)

// Show a simple color palette on the earring.
func showPalette(ledIndex uint8, palette *[18]pixel.LinearGRB888) {
	leds[ledIndex] = palette[ledIndex]
}

// Red (fire) and blue (ice) swirling around in circles.
func fireAndIce(index, traceIndex uint8) {
	// This animation has two tracers.
	traceIndex2 := traceIndex + uint8(len(leds))/2
	if traceIndex2 >= uint8(len(leds)) {
		traceIndex2 -= uint8(len(leds))
	}

	const div = 4
	if index == traceIndex {
		// First tracer.
		leds[index].R = leds[index].R/2 + 0xff/div
		leds[index].G = leds[index].G/2 + 0x33/div
		leds[index].B = leds[index].B / 2
	} else if index == traceIndex2 {
		// Second tracer, offset 180° on the LED ring.
		leds[index].R = leds[index].R/2 + 0x08/div
		leds[index].G = leds[index].G / 2
		leds[index].B += 0xff / div
	} else {
		// Tails, dim the LEDs.
		c := leds[index]
		if c.R > c.B {
			// Fire. Dim the red a bit.
			c.R = uint8(uint16(c.R) * 225 / 256)
			c.G = uint8(uint16(c.G) * 225 / 256)
			if c.R < 8 {
				c.R = 8
			}
		} else {
			// Ice. Dim the blue a bit.
			c.B = uint8(uint16(c.B) * 225 / 256)
		}
		leds[index] = c
	}
}

var sparkleIndex uint8
var sparkleIntensity uint8

// Sparkle! Glitter!
// Fade in some LEDs relatively quicly, and fade all LEDs out slowly.
func sparkle(index uint8) {
	const maxIntensity = 16 // must be a power of two
	if index == 0 {
		if sparkleIntensity == 0 {
			// Current LED is turned on. Move on to the next, at some random
			// time.
			r := random()
			if r%16 == 0 {
				// Pick the next LED to turn on.
				idx := uint8(r >> 16)
				for idx >= 18 {
					idx -= 18
				}
				sparkleIndex = idx
				if leds[sparkleIndex] == (pixel.LinearGRB888{}) {
					// Only sparkle dark pixels, to avoid mixing colors.
					sparkleIntensity = uint8(uint16(r)>>8) % maxIntensity
				}
			}
		}
	}
	if index == sparkleIndex && sparkleIntensity != 0 {
		// Turn on the given LED quickly.
		sparkleIntensity--
		sparkleColor := ledsgo.PartyColors.ColorAt(cycle * 64)
		c := leds[index]
		c.R += sparkleColor.R / maxIntensity
		c.G += sparkleColor.G / maxIntensity
		c.B += sparkleColor.B / maxIntensity
		leds[index] = c
	} else {
		// dim LED
		if cycle%8 == 0 {
			c := leds[index]
			c.R = uint8(uint16(c.R) * 253 / 256)
			c.G = uint8(uint16(c.G) * 253 / 256)
			c.B = uint8(uint16(c.B) * 253 / 256)
			leds[index] = c
		}
	}
}

var xorshift32State uint32 = 1

func random() uint32 {
	x := xorshift32State
	x = xorshift32(x)
	xorshift32State = x
	return x
}

// Xorshift32 RNG. The usual algorithm uses the shift amounts [13, 17, 5], but
// [7, 1, 9] as used below are a much better fit for AVR. It is "only" 37
// instructions (excluding return) when compiled with Clang 16 while the usual
// algorithm uses 57 instructions.
// On the other hand, avr-gcc (tested most versions starting with 5.4.0) is just
// terrible in both cases, using loops for these shifts.
//
// The [7, 1, 9] algorithm is mentioned on page 9 of this paper:
// http://www.iro.umontreal.ca/~lecuyer/myftp/papers/xorshift.pdf
//
// Godbolt reference (both algorithms in avr-gcc and Clang 16):
// https://godbolt.org/z/KdeKhP54d
func xorshift32(x uint32) uint32 {
	x ^= x << 7
	x ^= x >> 1
	x ^= x << 9
	return x
}
