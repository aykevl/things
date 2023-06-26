package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/pixel"
)

func testTearing[T pixel.Color](display tinygl.Displayer, screen *tinygl.Screen[T], touchInput board.TouchInput) {
	var (
		white = pixel.NewColor[T](255, 255, 255)
		black = pixel.NewColor[T](0, 0, 0)
		green = pixel.NewColor[T](0, 255, 0)
		red   = pixel.NewColor[T](255, 0, 0)
	)
	const (
		size  = 16
		speed = 4
	)

	// Create the canvas.
	width, height := screen.Size()
	canvas := gfx.NewCanvas(black, 32, 32)
	verticalRect := canvas.CreateRect(0, 0, size, height, white)
	horizontalRect := canvas.CreateRect(0, 0, width, size, white)
	teIndicator := canvas.CreateRect(0, 0, 8, 8, red)
	screen.SetChild(canvas)

	// Change the mode on each touch on the screen.
	mode := 0
	canvas.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event == tinygl.TouchStart {
			mode = (mode + 1) % 4
		}
	})

	lastRenderStart := time.Now()
	for cycle := 0; ; cycle++ {
		// Read keyboard inputs.
		board.Buttons.ReadInput()
		for {
			// Read keyboard events.
			event := board.Buttons.NextEvent()
			if event == board.NoKeyEvent {
				break
			}
			if event.Pressed() {
				switch event.Key() {
				case board.KeyA, board.KeyEnter:
					mode = (mode + 1) % 4
				case board.KeyB, board.KeyEscape:
					return
				}
			}
		}

		// Handle touch inputs.
		touches := touchInput.ReadTouch()
		if len(touches) > 0 {
			screen.SetTouchState(touches[0].X, touches[0].Y)
		} else {
			screen.SetTouchState(-1, -1)
		}

		moveHorizontal := mode/2 == 1
		avoidTearing := mode%2 == 1

		// Update dot in the top left corner that indicates whether tearing is
		// avoided or not.
		horizontalRect.SetHidden(moveHorizontal)
		verticalRect.SetHidden(!moveHorizontal)
		if avoidTearing {
			teIndicator.SetColor(green)
		} else {
			teIndicator.SetColor(red)
		}

		// Update vertical/horizontal moving bars.
		pixels := int(lastRenderStart.UnixNano() / (int64(1e9) / (speed * 60)))
		if moveHorizontal {
			maxOffset := (width*2 - size*2)
			offset := pixels % maxOffset
			if offset >= maxOffset/2 {
				offset = maxOffset - offset
			}
			verticalRect.Move(offset, 0)
		} else {
			maxOffset := (height*2 - size*2)
			offset := pixels % maxOffset
			if offset >= maxOffset/2 {
				offset = maxOffset - offset
			}
			horizontalRect.Move(0, offset)
		}

		// Render next frame.
		if avoidTearing {
			board.Display.WaitForVBlank(time.Second / 60)
		} else {
			duration := lastRenderStart.Add(time.Second / 60).Sub(time.Now())
			time.Sleep(duration)
		}
		renderStart := time.Now()
		screen.Update()
		renderDuration := time.Since(renderStart)
		if renderDuration > time.Millisecond*7 {
			println("rendering took:", renderDuration.String())
		}
		lastRenderStart = renderStart
	}
}
