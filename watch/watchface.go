package main

import (
	_ "embed"
	"time"

	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/image"
	"tinygo.org/x/drivers/pixel"
	"tinygo.org/x/tinyfont/freesans"
)

//go:embed assets/watchface-0.raw
var char0 string

//go:embed assets/watchface-1.raw
var char1 string

//go:embed assets/watchface-2.raw
var char2 string

//go:embed assets/watchface-3.raw
var char3 string

//go:embed assets/watchface-4.raw
var char4 string

//go:embed assets/watchface-5.raw
var char5 string

//go:embed assets/watchface-6.raw
var char6 string

//go:embed assets/watchface-7.raw
var char7 string

//go:embed assets/watchface-8.raw
var char8 string

//go:embed assets/watchface-9.raw
var char9 string

//go:embed assets/watchface-colon.raw
var charColon string

// Create a simple digital watch face as the homescreen.
func (w *Watch[T]) createWatchFace(views *ViewManager[T]) View[T] {
	// TODO: make this configurable somehow.
	//return w.createTextWatchface(views)
	return w.createDigitalWatchface(views)
}

func (w *Watch[T]) createTextWatchface(views *ViewManager[T]) View[T] {
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

func (w *Watch[T]) createDigitalWatchface(views *ViewManager[T]) View[T] {
	const charWidth, charHeight = 48, 107
	var (
		black = pixel.NewColor[T](0, 0, 0)
		white = pixel.NewColor[T](255, 255, 255)

		// TODO: compress the images in some way (for example, using run-length
		// encoding). Right now the images together consume 7062 bytes of flash,
		// which is quite a lot.
		numbers = [10]image.Mono[T]{
			image.MakeMono(white, black, charWidth, charHeight, char0),
			image.MakeMono(white, black, charWidth, charHeight, char1),
			image.MakeMono(white, black, charWidth, charHeight, char2),
			image.MakeMono(white, black, charWidth, charHeight, char3),
			image.MakeMono(white, black, charWidth, charHeight, char4),
			image.MakeMono(white, black, charWidth, charHeight, char5),
			image.MakeMono(white, black, charWidth, charHeight, char6),
			image.MakeMono(white, black, charWidth, charHeight, char7),
			image.MakeMono(white, black, charWidth, charHeight, char8),
			image.MakeMono(white, black, charWidth, charHeight, char9),
		}
	)

	_, displayHeight := w.display.Size()
	canvas := gfx.NewCanvas(black, 96, 96)
	colonImage := image.MakeMono(white, black, charWidth, charHeight, charColon)

	// Create hh:mm images.
	h0 := gfx.NewImage[T](numbers[0], charWidth*0, int(displayHeight)/2-charHeight/2)
	h1 := gfx.NewImage[T](numbers[0], charWidth*1, int(displayHeight)/2-charHeight/2)
	colon := gfx.NewImage[T](colonImage, charWidth*2, int(displayHeight)/2-charHeight/2)
	m0 := gfx.NewImage[T](numbers[0], charWidth*3, int(displayHeight)/2-charHeight/2)
	m1 := gfx.NewImage[T](numbers[0], charWidth*4, int(displayHeight)/2-charHeight/2)
	canvas.Add(h0)
	canvas.Add(h1)
	canvas.Add(colon)
	canvas.Add(m0)
	canvas.Add(m1)

	eventWrapper := tinygl.NewEventBox[T](canvas)
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

	updateTime := func(hour, minute int) {
		h0.SetImage(numbers[hour/10])
		h1.SetImage(numbers[hour%10])
		m0.SetImage(numbers[minute/10])
		m1.SetImage(numbers[minute%10])
	}
	now := watchTime()
	hour := now.Hour()
	minute := now.Minute()
	updateTime(hour, minute)

	return NewView[T](eventWrapper, func(now time.Time) {
		// Update the watchface.
		if backlight > 0 {
			// Watch face is visible.
			newHour := now.Hour()
			newMinute := now.Minute()
			if hour != newHour || minute != newMinute {
				hour = newHour
				minute = newMinute
				updateTime(hour, minute)
			}
		}
	})
}
