package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/pixel"
)

const (
	maxIterations = 50
	frac          = 12 // fractional bits
)

// Render a mandelbrot. Allow the user to move around the screen and zoom in/out
// the fractal.
func mandelbrot[T pixel.Color](screen *tinygl.Screen[T]) {
	width, height := screen.Size()
	stepY := int(2 << (frac * 2) / int64(height))
	stepX := int((3 << (frac * 2)) / int64(width))
	centerX := stepX * width / -6
	centerY := 0

	black := pixel.NewColor[T](0, 0, 0)
	canvas := gfx.NewCustomCanvas(black, 0, 0, func(screen *tinygl.Screen[T], displayX, displayY, displayWidth, displayHeight, x, y int) {
		buffer := screen.Buffer()
		bufferLines := buffer.Len() / displayWidth
		start := time.Now()
		i := centerY - (displayHeight/2)*stepY
		for startY := 0; startY < displayHeight; startY += bufferLines {
			chunkHeight := bufferLines
			if startY+chunkHeight >= displayHeight {
				chunkHeight = displayHeight - startY
			}
			img := buffer.Rescale(displayWidth, chunkHeight)
			for chunkY := 0; chunkY < chunkHeight; chunkY++ {
				r := centerX - (displayWidth/2)*stepX
				i += stepY
				for x := 0; x < displayWidth; x++ {
					r += stepX
					iterations := mandelbrotAt(r>>frac, i>>frac)
					//iterations := mandelbrotPreciseAt(r, i)
					rawColor := pixel.NewColor[T](0, 0, 0)
					if iterations != 255 {
						c := ledsgo.RainbowColors.ColorAt(uint16(iterations * 2048))
						rawColor = pixel.NewColor[T](c.R, c.G, c.B)
					}
					img.Set(x, chunkY, rawColor)
				}
			}
			screen.Send(displayX, displayY+startY, img)
		}
		duration := time.Since(start)
		println("rendering took:", duration.String())
	})
	screen.SetChild(canvas)

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
				canvas.RequestUpdate()
			case board.KeyB:
				stepX = stepX * 3 / 2
				stepY = stepY * 3 / 2
				canvas.RequestUpdate()
			case board.KeyLeft:
				centerX -= (width * stepX) / 8
				canvas.RequestUpdate()
			case board.KeyRight:
				centerX += (width * stepX) / 8
				canvas.RequestUpdate()
			case board.KeyUp:
				centerY -= (height * stepY) / 8
				canvas.RequestUpdate()
			case board.KeyDown:
				centerY += (height * stepY) / 8
				canvas.RequestUpdate()
			case board.KeyEscape:
				return
			}
		}

		screen.Update()
		board.Display.WaitForVBlank(time.Second / 30)
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
// This is a bit slower (~40%) on chips with a 32x32=64 multiply instruction
// (like smull on ARM, available in Cortex-M3 and above). It is _much_ slower on
// chips without such a multiply instruction.
func mandelbrotPreciseAt(_x0, _y0 int) int {
	const frac2 = 24
	x0 := int32(_x0) >> (frac*2 - frac2)
	y0 := int32(_y0) >> (frac*2 - frac2)
	x := int32(0)
	y := int32(0)
	iteration := 1
	x2 := int64(0) // .frac*2
	y2 := int64(0) // .frac*2
	for x2+y2 <= 4<<(frac2*2) {
		y = int32((int64(x)*int64(y))>>(frac2-1)) + y0
		x = int32((x2-y2)>>frac2) + x0
		x2 = int64(x) * int64(x)
		y2 = int64(y) * int64(y)
		iteration++
		if iteration == maxIterations {
			return 255
		}
	}
	return iteration
}
