package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"tinygo.org/x/drivers/pixel"
)

const (
	// Rectangle size for which a noise value is calculated. Must be a power of
	// two.
	// Lower values mean better quality, but increased computation needs.
	noisePixelSize = 8

	// Noise size. A large value means more detail (higher frequency noise).
	noiseSize = 16
)

func noise[T pixel.Color](display board.Displayer[T], screen *tinygl.Screen[T]) {
	var colorMap [256]T
	for i := range colorMap {
		c := ledsgo.Color{H: uint16(i * 256), S: 255, V: 255}.Rainbow()
		colorMap[i] = pixel.NewLinearColor[T](c.R, c.G, c.B)
	}

	black := pixel.NewColor[T](0, 0, 0)
	var drawTimeSum time.Duration
	canvas := gfx.NewCustomCanvas[T](black, 0, 0, func(screen *tinygl.Screen[T], displayX, displayY, displayWidth, displayHeight, x, y int) {
		// TODO: many parameters (displayX, displayY, x, y) are ignored. They
		// should be used to paint noise at the right location on the screen.

		noiseWidth := (int(displayWidth)+noisePixelSize-1)/noisePixelSize + 1
		noiseLine := make([]uint16, noiseWidth)

		now := time.Now()
		z := uint32(now.UnixNano() >> 20)
		for i := 0; i < noiseWidth; i++ {
			value := ledsgo.Noise3(0, uint32(i*noisePixelSize)*noiseSize, z)
			noiseLine[i] = value
		}
		for y := 0; y < int(displayHeight); {
			buf := screen.Buffer()
			lines := buf.Len() / int(displayWidth)
			if y+lines >= int(displayHeight) {
				lines = int(displayHeight) - y
			}
			lines -= lines % noisePixelSize // round down
			buf = buf.Rescale(displayWidth, lines)
			for chunkY := 0; chunkY < lines; chunkY += noisePixelSize {
				valueTopLeft := noiseLine[0]
				valueBottomLeft := ledsgo.Noise3(uint32(y+chunkY+noisePixelSize)*noiseSize, 0, z)
				noiseLine[0] = valueBottomLeft
				for x := 0; x < int(displayWidth); x += noisePixelSize {
					valueTopRight := noiseLine[x/noisePixelSize+1]
					valueBottomRight := ledsgo.Noise3(uint32(y+chunkY+noisePixelSize)*noiseSize, uint32(x+noisePixelSize)*noiseSize, z)
					noiseLine[x/noisePixelSize+1] = valueBottomRight
					interpolatedTop := int(valueTopLeft) * noisePixelSize
					interpolatedDiffTop := int(valueTopRight) - int(valueTopLeft)
					interpolatedBottom := int(valueBottomLeft) * noisePixelSize
					interpolatedDiffBottom := int(valueBottomRight) - int(valueBottomLeft)
					for posX := 0; posX < noisePixelSize; posX++ {
						interpolatedDiffY := (interpolatedBottom - interpolatedTop) / noisePixelSize
						interpolated := interpolatedTop
						for posY := 0; posY < noisePixelSize; posY++ {
							c := colorMap[uint(interpolated/128/noisePixelSize)%256]
							buf.Set(x+posX, chunkY+posY, c)
							interpolated += interpolatedDiffY
						}
						interpolatedTop += interpolatedDiffTop
						interpolatedBottom += interpolatedDiffBottom
					}
					valueTopLeft = valueTopRight
					valueBottomLeft = valueBottomRight
				}
			}
			drawStart := time.Now()
			screen.Send(0, y, buf)
			drawTimeSum += time.Since(drawStart)
			y += lines
		}

	})
	screen.SetChild(canvas)

	var frame uint32
	frameSumStart := time.Now()
	for {
		// Handle keyboard inputs, to exit from the noise demo.
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
			case board.KeyEscape, board.KeyB:
				if event.Pressed() {
					return
				}
			}
		}

		// The canvas needs to be updated each cycle, so request an update.
		canvas.RequestUpdate()
		screen.Update()

		// Don't wait for VBlank to improve rendering speed.

		// Print draw statistics.
		// It prints the time each frame takes (excluding the print line below),
		// and the time spent starting the next transmission to the display.
		const numFrames = 32
		frame++
		if frame%numFrames == 0 {
			frameTimeSum := time.Since(frameSumStart)
			println("time:", (frameTimeSum / numFrames).String(), (drawTimeSum / numFrames).String())
			frameSumStart = time.Now()
			drawTimeSum = 0
		}
	}
}
