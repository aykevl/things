package main

import (
	"time"

	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl/pixel"
)

var leds [18]pixel.LinearGRB888

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
						if animation >= 3 {
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

			switch animation {
			case 1:
				rainbowTrace(ledIndex, traceIndex)
			case 2:
				purpleCircles(ledIndex, traceIndex)
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
