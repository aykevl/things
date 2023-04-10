package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/pixel"
)

func testTouch[T pixel.Color](screen *tinygl.Screen[T], touchInput board.TouchInput) {
	// Determine size and scale of the screen.
	width, _ := board.Display.Size()
	physicalWidth, _ := board.Display.PhysicalSize()
	scalePercent := int(width) * 21 / physicalWidth

	// Create canvas.
	black := pixel.NewColor[T](0, 0, 0)
	canvas := gfx.NewCanvas(black, 8, 8)
	screen.SetChild(canvas)

	// Create touch point.
	touch := canvas.CreateRect(0, 0, scalePercent/4, scalePercent/4, pixel.NewColor[T](255, 255, 255))
	touch.SetHidden(true)

	// Show screen.
	screen.Update()

	for {
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
				case board.KeyB, board.KeyEscape:
					return
				}
			}
		}

		// Read touch inputs.
		// TODO: be able to show multiple touch points (multitouch isn't
		// implemented yet).
		touches := touchInput.ReadTouch()
		if len(touches) != 0 {
			_, _, width, height := touch.Bounds()
			touch.Move(int(touches[0].X)-width/2, int(touches[0].Y)-height/2)
			touch.SetHidden(false)
		} else {
			touch.SetHidden(true)
		}
		screen.Update()

		time.Sleep(time.Second / 30)
	}
}
