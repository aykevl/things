package main

import (
	"time"

	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/pixel"
	"tinygo.org/x/tinyfont/freesans"
)

// Create a simple digital watch face as the homescreen.
func (w *Watch[T]) createWatchFace(views *ViewManager[T]) View[T] {
	var (
		black = pixel.NewColor[T](0, 0, 0)
		white = pixel.NewColor[T](255, 255, 255)
	)

	now := watchTime()
	hour := now.Hour()
	minute := now.Minute()
	timeText := tinygl.NewText(&freesans.Regular24pt7b, white, black, formatTime(hour, minute))
	eventWrapper := tinygl.NewEventBox[T](timeText)
	eventWrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event == tinygl.TouchTap {
			if backlight == 0 {
				// Tapped on a sleeping watch.
				// Awake the screen.
				w.exitSleep()
			} else {
				// Regular tap on the clock.
				// TODO: detect gesture (for example, swipe upwards) to make it
				// harder to accidentally get in the settings menu.
				views.Push(w.createAppsView(views))
			}
		}
	})

	return NewView[T](eventWrapper, func(now time.Time) {
		// Update the watchface.
		if backlight > 0 {
			// Watch face is visible.
			newHour := now.Hour()
			newMinute := now.Minute()
			if hour != newHour || minute != newMinute {
				hour = newHour
				minute = newMinute
				timeText.SetText(formatTime(hour, minute))
			}
		}
	})
}
