package main

import (
	"time"
	"image/color"

	"github.com/aykevl/ledsgo"
)

const (
	brightness = 0xff
	spread     = 6  // higher means the noise gets more detailed
	speed      = 20 // higher means slower
)

func main() {
	for i := 0; i < 2; i++ {
		println("sleep", i)
		time.Sleep(1 * time.Second)
	}
	fullRefreshes := uint(0)
	previousSecond := int64(0)
	for {
		start := time.Now()
		//noise(start)
		fire()

		second := (time.Now().UnixNano() / int64(time.Second))
		if second != previousSecond {
			previousSecond = second
			newFullRefreshes := display.FullRefreshes()
			print("#", second, " screen=", newFullRefreshes-fullRefreshes, "fps animation=", time.Since(start).String(), "\r\n")
			fullRefreshes = newFullRefreshes
		}
		display.Display()
	}
}

func noise(now time.Time) {
	for x := int16(0); x < 32; x++ {
		for y := int16(0); y < 32; y++ {
			hue := uint16(ledsgo.Noise3(int32(now.UnixNano()>>speed), int32(x<<spread), int32(y<<spread))) * 2
			display.SetPixel(x, y, ledsgo.Color{hue, 0xff, brightness}.Spectrum())
		}
	}
}

func fire() {
	const pointsPerCircle = 12 // how many LEDs there are per turn of the torch
	const cooling = 8          // higher means faster cooling
	const detail = 400         // higher means more detailed flames
	const speed = 12           // higher means faster
	now := time.Now().UnixNano()
	for x := int16(0); x < 32; x++ {
		for y := int16(0); y < 32; y++ {
			heat := ledsgo.Noise2(int32(y*detail)-int32((now>>20)*speed), int32(x*detail))/256 + 128
			heat -= int16(y) * cooling
			if heat < 0 {
				heat = 0
			}
			display.SetPixel(x, y, heatMap(uint8(heat)))
		}
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
