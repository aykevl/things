package main

import (
	"time"
	"unsafe"

	"github.com/aykevl/board"
	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl/pixel"
)

const (
	maxIterations = 50
	frac          = 12 // fractional bits
)

// Render a mandelbrot. Allow the user to move around the screen and zoom in/out
// the fractal.
func mandelbrot[T pixel.Color](display board.Displayer[T], buffer []T) {
	w, h := display.Size()
	width := int(w)
	height := int(h)
	bufferLines := len(buffer) / width
	stepY := int(2 << (frac * 2) / int64(height))
	stepX := int((3 << (frac * 2)) / int64(width))
	centerX := stepX * width / -6
	centerY := 0
	needsRender := true
	for {
		board.Buttons.ReadInput()
		for {
			event := board.Buttons.NextEvent()
			if event == board.NoKeyEvent {
				break
			}
			if !event.Pressed() {
				continue
			}
			switch event.Key() {
			case board.KeyA:
				stepX = stepX * 2 / 3
				stepY = stepY * 2 / 3
				needsRender = true
			case board.KeyB:
				stepX = stepX * 3 / 2
				stepY = stepY * 3 / 2
				needsRender = true
			case board.KeyLeft:
				centerX -= (width * stepX) / 8
				needsRender = true
			case board.KeyRight:
				centerX += (width * stepX) / 8
				needsRender = true
			case board.KeyUp:
				centerY -= (height * stepY) / 8
				needsRender = true
			case board.KeyDown:
				centerY += (height * stepY) / 8
				needsRender = true
			case board.KeyEscape:
				return
			}
		}

		if needsRender {
			needsRender = false
			start := time.Now()
			i := centerY - (height/2)*stepY
			for startY := 0; startY < height; startY += bufferLines {
				chunkHeight := bufferLines
				if startY+chunkHeight >= height {
					chunkHeight = height - startY
				}
				for chunkY := 0; chunkY < chunkHeight; chunkY++ {
					r := centerX - (width/2)*stepX
					i += stepY
					for x := 0; x < width; x++ {
						r += stepX
						iterations := mandelbrotAt(r>>frac, i>>frac)
						rawColor := pixel.NewColor[T](0, 0, 0)
						if iterations != 255 {
							c := ledsgo.RainbowColors.ColorAt(uint16(iterations * 2048))
							rawColor = pixel.NewColor[T](c.R, c.G, c.B)
						}
						buffer[chunkY*width+x] = rawColor
					}
				}
				var zeroPixel T
				buf8 := unsafe.Slice((*uint8)(unsafe.Pointer(unsafe.SliceData(buffer))), int(width)*chunkHeight*int(unsafe.Sizeof(zeroPixel)))
				display.DrawRGBBitmap8(0, int16(startY), buf8, int16(width), int16(chunkHeight))
			}
			duration := time.Since(start)
			println("rendering took:", duration.String())
		}
		display.Display()
		time.Sleep(time.Second / 30)
	}
}

func mandelbrotAt(x0, y0 int) int {
	// This check is expensive, so don't do it.
	//if x0 < -2<<frac || x0 > 2<<frac || y0 < -2<<frac || y0 > 2<<frac {
	//	// Avoid integer overflow by not calculating values so far away from the center.
	//	return 2
	//}

	// This is the final optimized algorithm from Wikipedia:
	// https://en.wikipedia.org/wiki/Plotting_algorithms_for_the_Mandelbrot_set#Optimized_escape_time_algorithms
	x := 0
	y := 0
	iteration := 1
	x2 := 0 // .frac*2
	y2 := 0 // .frac*2
	for x2+y2 <= 4<<(frac*2) {
		y = (x*y)>>(frac-1) + y0
		x = (x2-y2)>>frac + x0
		x2 = x * x
		y2 = y * y
		iteration++
		if iteration == maxIterations {
			return 255
		}
	}
	return iteration
}

// Improved precision version of the mandelbrot function.
// With this, you can increase frac to a larger value.
func mandelbrotPreciseAt(x0, y0 int) int {
	x := 0
	y := 0
	iteration := 1
	x2 := 0 // .frac*2
	y2 := 0 // .frac*2
	for x2+y2 <= 4<<(frac) {
		// 2*x*y is replaced by square(x+y) - square(x) - square(y).
		// See: https://www.reenigne.org/blog/arbitrary-precision-mandelbrot-sets/
		y = square(x+y) - x2 - y2 + y0
		x = (x2 - y2) + x0
		x2 = square(x)
		y2 = square(y)
		iteration++
		if iteration == maxIterations {
			return 255
		}
	}
	return iteration
}

// Return the square of x, at the frac fixed point level.
func square(x int) int {
	// TODO: I don't think this is correct with frac > 16.
	top := x >> frac
	bottom := x & (1<<frac - 1)
	return (top*bottom)*2 + (top*top)<<frac + (bottom*bottom)>>frac
}
