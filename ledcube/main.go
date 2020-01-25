package main

import (
	"image/color"
	"time"

	"github.com/aykevl/ledsgo"
)

func main() {
	fullRefreshes := uint(0)
	previousSecond := int64(0)
	//demo := colorCoordinateAt
	demo := noiseAt
	//demo := fireAt
	for {
		start := time.Now()
		drawPixels(start, demo)
		display.Display()

		second := (start.UnixNano() / int64(time.Second))
		if second != previousSecond {
			previousSecond = second
			newFullRefreshes := display.FullRefreshes()
			animationTime := time.Since(start)
			animationFPS := int64(10 * time.Second / animationTime)
			print("#", second, " screen=", newFullRefreshes-fullRefreshes, "fps animation=", animationTime.String(), "/", (animationFPS / 10), ".", animationFPS%10, "fps\r\n")
			fullRefreshes = newFullRefreshes
		}
	}
}

// drawPixels updates every pixel on the cube by calling getColor for each pixel
// and drawing it to the screen. It maps virtual (3D) pixels to physical (2D)
// pixels in the process.
func drawPixels(t time.Time, getColor func(x, y, z int, t time.Time) color.RGBA) {
	// Somewhat arbitrarily picking the top left of the topmost panel as the (0,
	// 0, 31) of the 3D cube.
	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			display.SetPixel(int16(x+32*5), int16(y), getColor(x+1, y+1, 0, t))
			display.SetPixel(int16(x+32*4), int16(y), getColor(0, x+1, y+1, t))
			display.SetPixel(int16(x+32*3), int16(y), getColor(32-x, 0, y+1, t))
			display.SetPixel(int16(x+32*2), int16(y), getColor(33, 32-x, y+1, t))
			display.SetPixel(int16(x+32*1), int16(y), getColor(x+1, 33, y+1, t))
			display.SetPixel(int16(x+32*0), int16(y), getColor(x+1, 32-y, 33, t))
		}
	}
}

// noiseAt returns noise at the specified location.
func noiseAt(x, y, z int, t time.Time) color.RGBA {
	const (
		spread = 7  // higher means the noise gets more detailed
		speed  = 20 // higher means slower
	)
	hue := uint16(ledsgo.Noise4(int32(t.UnixNano()>>speed), int32(x<<spread), int32(y<<spread), int32(z<<spread))) * 2
	return ledsgo.Color{hue, 0xff, 0xff}.Spectrum()
}

// fireAt returns fire at the specified location.
func fireAt(x, y, z int, t time.Time) color.RGBA {
	const pointsPerCircle = 12 // how many LEDs there are per turn of the torch
	const cooling = 56         // higher means faster cooling
	const detail = 400         // higher means more detailed flames
	const speed = 12           // higher means faster
	const screenHeight = 33
	if z == 0 {
		return color.RGBA{}
	}
	heat := ledsgo.Noise3(int32((31-z)*detail)-int32((t.UnixNano()>>20)*speed), int32(x*detail), int32(y*detail))/32 + (128 * 8)
	heat -= int16(screenHeight-z) * cooling
	if heat < 0 {
		heat = 0
	}
	return heatMap(heat)
}

// heatMap maps a color in the range 0..2047 to a color in a heat index. Useful
// for making flames.
func heatMap(index int16) color.RGBA {
	if index < 128*8 {
		// red only
		return color.RGBA{uint8(index / 4), 0, 0, 255}
	}
	if index < 224*8 {
		// red-yellow
		return color.RGBA{255, uint8(uint32(index-128*8) / 3), 0, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (uint8(index - 224*8)), 255}
}

// colorCoordinateAt returns a color based on the 3 coordinates given. Useful
// for getting the virtual->physical pixel mapping right.
func colorCoordinateAt(x, y, z int, t time.Time) color.RGBA {
	// X represents red (more red to the right)
	// Y represents green (more green to the bottom)
	// Z represents blue (more blue to the bottom)
	return color.RGBA{uint8(x * 255 / 33), uint8(y * 255 / 33), uint8(z * 255 / 33), 0xff}
}
