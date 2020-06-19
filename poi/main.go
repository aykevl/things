// Connects to an APA102 SPI RGB LED strip with 30 LEDS.
package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/apa102"
)

var leds = make([]color.RGBA, 60)

const (
	height = 14
)

// Parameters that are controlled with Bluetooth.
var (
	animationIndex int   = 0
	speed          uint8 = 10
)

var animations = []func(time.Time){
	solid,
	noise,
	fire,
	iris,
	gear,
	halfcircles,
	arrows,
	black,
}

func main() {
	println("starting")
	initHardware()

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: spiFrequency,
		Mode:      0,
		SCK:       spiClockPin,
		MOSI:      spiDataPin,
		MISO:      machine.NoPin,
	})

	a := apa102.New(machine.SPI0)

	for i := uint(0); ; i++ {
		now := time.Now()
		animation := animations[animationIndex]
		animation(now)

		for i := 0; i < 16; i++ {
			leds[29-i] = leds[i]
		}
		a.WriteColors(leds)

		// print speed
		if i%100 == 0 {
			//duration := time.Since(now)
			//println("duration:", duration.String())
		}
	}
}

func noise(now time.Time) {
	const x = 0
	const spread = 7
	for y := int16(0); y < height; y++ {
		hue := uint16(ledsgo.Noise2(int32(now.UnixNano()>>(26-speed)), int32(y<<spread))) * 2
		c := ledsgo.Color{hue, 0xff, 0xff}.Spectrum()
		c.A = baseColor.A
		leds[y] = c
	}
}

func iris(now time.Time) {
	expansion := (ledsgo.Noise1(int32(now.UnixNano()>>(21-speed))) / 256) + 128 - 50
	for y := int16(0); y < height; y++ {
		intensity := expansion - y*16
		if intensity < 0 {
			intensity = 0
		}
		c := ledsgo.ApplyAlpha(baseColor, uint8(intensity))
		c.A = baseColor.A
		leds[y] = c
	}
}

func gear(now time.Time) {
	long := int16((now.UnixNano()>>(33-speed))%8) == 0
	for y := int16(0); y < height; y++ {
		c := color.RGBA{}
		if long || y < 3 {
			c = baseColor
		}
		leds[y] = c
	}
}

func halfcircles(now time.Time) {
	chosenOne := int16((now.UnixNano() >> (32 - speed)) % height)
	for y := int16(0); y < height; y++ {
		c := color.RGBA{}
		if y == chosenOne || y == chosenOne+1 {
			c = baseColor
		}
		leds[y] = c
	}
}

func arrows(now time.Time) {
	// First make them all black.
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}

	// Turn the two LEDs on that are part of the arrow.
	index := int16((now.UnixNano() >> (32 - speed)) % (height / 2))
	leds[index] = baseColor
	leds[height-1-index] = baseColor
}

func solid(now time.Time) {
	for y := int16(0); y < height; y++ {
		leds[y] = baseColor
	}
}

func black(now time.Time) {
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}
}

func fire(now time.Time) {
	var cooling = (14 * 16) / height // higher means faster cooling
	const detail = 400               // higher means more detailed flames
	for y := int16(0); y < height; y++ {
		heat := ledsgo.Noise1(int32((now.UnixNano()>>(23-speed)))-int32(y)*detail)/256 + 128
		heat -= y * int16(cooling)
		if heat < 0 {
			heat = 0
		}
		c := heatMap(uint8(heat))
		c.A = baseColor.A
		leds[y] = c
	}
}

func heatMap(index uint8) color.RGBA {
	if index < 128 {
		return color.RGBA{index * 2, 0, 0, 255}
	}
	if index < 224 {
		return color.RGBA{255, uint8(uint32(index-128) * 8 / 3), 0, 255}
	}
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}
