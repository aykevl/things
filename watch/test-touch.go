package main

import (
	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"tinygo.org/x/drivers/pixel"
)

func createTouchTestView[T pixel.Color](views *ViewManager[T]) View[T] {
	// Determine size and scale of the screen.
	scalePercent := board.Display.PPI() * 100 / 120

	// Create canvas.
	black := pixel.NewColor[T](0, 0, 0)
	canvas := gfx.NewCanvas(black, 8, 8)

	// Create touch point.
	touch := gfx.NewRect(pixel.NewColor[T](255, 255, 255), 0, 0, scalePercent/4, scalePercent/4)
	canvas.Add(touch)
	touch.SetHidden(true)

	// Listen for touch events.
	wrapper := tinygl.NewEventBox[T](canvas)
	wrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		switch event {
		case tinygl.TouchStart:
			// Finger touched down on the screen.
			touch.SetHidden(false)
			_, _, w, h := touch.Bounds()
			touch.Move(x-w/2, y-h/2)
		case tinygl.TouchMove:
			// Finger moved.
			_, _, w, h := touch.Bounds()
			touch.Move(x-w/2, y-h/2)
		case tinygl.TouchEnd:
			// Finger is now removed.
			touch.SetHidden(true)
		}
	})

	return NewView[T](wrapper, nil)
}
