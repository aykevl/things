// Connects to an APA102 SPI RGB LED strip with 30 LEDS.
package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"github.com/spaolacci/murmur3"
	"tinygo.org/x/drivers/apa102"
)

var leds = make([]color.RGBA, 60)

const (
	height = 14
)

// Parameters that are controlled with Bluetooth.
var (
	animationIndex uint8 = 0
	speed          uint8 = 10
)

// List of animations that can be selected over BLE.
var animations = []func(time.Time){
	solid,
	noise,
	fire,
	iris,
	gear,
	halfcircles,
	arrows,
	glitter,
	black,
}

func main() {
	println("starting")
	initHardware()

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: spiFrequency,
		Mode:      0,
		SCK:       spiClockPin,
		SDO:       spiDataPin,
		SDI:       machine.NoPin,
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

// Colorful noise.
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

// Looks a bit like spikes from inside to the outside.
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

// Looks like a typical blocky gear, with square gear teeth.
func gear(now time.Time) {
	long := int16((now.UnixNano()>>(32-speed))%8) == 0
	for y := int16(0); y < height; y++ {
		c := color.RGBA{}
		if long || y < 3 {
			c = baseColor
		}
		leds[y] = c
	}
}

// Somewhat improperly named. When two poi are spinning in opposite direction,
// it has a somewhat palm tree like appearance.
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

// Simple > arrows pointing in the direction the poi is moving.
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

// Random colored specles. Looks great in the dark because the poi itself is
// (nearly) invisible showing only these speckles.
func glitter(now time.Time) {
	// Make all LEDs black.
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}

	// Get a random number based on the time.
	t := uint32(now.UnixNano() >> (32 - speed))
	hash := murmur3.Sum32([]byte{byte(t), byte(t >> 8), byte(t >> 16), byte(t >> 24)})

	// Use this number to get an index.
	index := hash % (height * 2)
	if index >= height {
		return // don't sparkle all the time
	}

	// Get a random color on the color wheel.
	c := ledsgo.Color{uint16((hash >> 7)), 0xff, 0xff}.Spectrum()
	c.A = baseColor.A

	leds[index] = c
}

// Solid color. Useful to reduce distraction, for testing and as a not too
// distracting startup color.
func solid(now time.Time) {
	for y := int16(0); y < height; y++ {
		leds[y] = baseColor
	}
}

// Disable all LEDs. LEDs will still consume power as they use around 1mA even
// when they're dark.
func black(now time.Time) {
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}
}

// Fire animation. The flame is the configured color, so you can not only have a
// red flame, but also a green or blue flame. Red looks the best IMHO, though.
func fire(now time.Time) {
	var cooling = (14 * 16) / height // higher means faster cooling
	const detail = 400               // higher means more detailed flames
	for y := int16(0); y < height; y++ {
		heat := ledsgo.Noise2(int32((now.UnixNano()>>(23-speed))), int32(y)*detail)/256 + 128
		heat -= y * int16(cooling)
		if heat < 0 {
			heat = 0
		}
		c := coloredFlame(uint8(heat))
		c.A = baseColor.A
		leds[y] = c
	}
}

// Colored flame. Like a heat map, but the lowest temperatures are not fixed red
// but instead use the configured color.
func coloredFlame(index uint8) color.RGBA {
	if index < 128 {
		// <color>
		c := ledsgo.ApplyAlpha(baseColor, index*2)
		c.A = 255
		return c
	}
	if index < 224 {
		// <color>-yellow
		c1 := ledsgo.ApplyAlpha(baseColor, 255-uint8(uint32(index-128)*8/3))
		c2 := ledsgo.ApplyAlpha(color.RGBA{255, 255, 0, 255}, uint8(uint32(index-128)*8/3))
		return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 255}
		//return color.RGBA{255, uint8(uint32(index-128) * 8 / 3), 0, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}

// Heat map, where lower numbers are colder. Previously used in the fire
// animation.
func heatMap(index uint8) color.RGBA {
	if index < 128 {
		return color.RGBA{index * 2, 0, 0, 255}
	}
	if index < 224 {
		// red-yellow
		return color.RGBA{255, uint8(uint32(index-128) * 8 / 3), 0, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}
