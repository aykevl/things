package main

import (
	"time"

	"github.com/aykevl/ledsgo"
	"github.com/aykevl/things/hub75"
)

var display *hub75.Device

const (
	brightness = 0xff
	spread     = 6  // higher means the noise gets more detailed
	speed      = 20 // higher means slower
)

func main() {
	display = hub75.New(22, 23, 24, 25, 20, 19, 18, 17)

	fullRefreshes := uint(0)
	previousSecond := int64(0)
	for {
		start := time.Now()
		for x := int16(0); x < 32; x++ {
			for y := int16(0); y < 32; y++ {
				hue := uint16(ledsgo.Noise3(int32(start.UnixNano()>>speed), int32(x<<spread), int32(y<<spread))) * 2
				display.SetPixel(x, y, ledsgo.Color{hue, 0xff, brightness}.Spectrum())
			}
		}
		display.Display()

		second := (time.Now().UnixNano() / int64(time.Second))
		if second != previousSecond {
			previousSecond = second
			newFullRefreshes := display.FullRefreshes()
			print("#", second, " screen=", newFullRefreshes-fullRefreshes, "fps animation=", time.Since(start).String(), "\r\n")
			fullRefreshes = newFullRefreshes
		}
	}
}
