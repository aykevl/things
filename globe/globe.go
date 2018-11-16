package main

import (
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygo-drivers/ws2812"
)

const ledPin = 21

const (
	whole = 0x10000 / 256 // scaled 1
	grinv = 0x09e37 / 256 // 0x10000 / golden ratio
)

// The X, Y and Z parts of the LED vectors from the center of the globe.
var (
	vectorx = []int32{    0,     0 , -whole, -grinv,  grinv,  whole,  grinv, -grinv, -whole,      0,  whole,      0}
	vectory = []int32{grinv, -grinv,      0,  whole,  whole,      0, -whole, -whole,      0,  grinv,      0, -grinv}
	vectorz = []int32{whole,  whole,  grinv,      0,      0,  grinv,      0,      0, -grinv, -whole, -grinv, -whole}
)

func main() {
	pin := machine.GPIO{ledPin}
	pin.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})
	strip := ws2812.New(pin)
	colors := ledsgo.Strip(make([]uint32, 12))
	for {
		now := time.Now().UnixNano()
		for i := range colors {
			x := vectorx[i] + int32(now >> 22)
			y := vectory[i]
			z := vectorz[i]
			hue := uint16(ledsgo.Noise3(x, y, z))
			colors[i] = ledsgo.Color{hue, 0xff, 0xff}.Spectrum()
		}
		strip.WriteColors(colors)
		time.Sleep(time.Millisecond)
	}
}
