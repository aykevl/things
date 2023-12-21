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

// Current watchface.
var watchFaceIndex uint8

// Create a simple digital watch face as the homescreen.
func (w *Watch[T]) createWatchFace(views *ViewManager[T]) View[T] {
	switch watchFaceIndex {
	case 0:
		return w.createTextWatchface(views)
	case 1:
		return w.createDigitalWatchface(views)
	case 2:
		return w.createAnalogWatchface(views)
	default:
		// should be unreachable
		return w.createTextWatchface(views)
	}
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

func (w *Watch[T]) createAnalogWatchface(views *ViewManager[T]) View[T] {
	var (
		black             = pixel.NewColor[T](0, 0, 0)
		hourMark          = pixel.NewColor[T](100, 100, 100)
		hourHandleColor   = pixel.NewColor[T](255, 255, 255)
		minuteHandleColor = pixel.NewColor[T](255, 255, 255)
		centerDotColor    = pixel.NewColor[T](255, 255, 255)
	)

	canvas := gfx.NewCanvas(black, 96, 96)
	displayWidth, displayHeight := w.display.Size()
	centerX, centerY := int(displayWidth)/2, int(displayHeight)/2
	r := int(displayHeight) / 2

	// Hour markers.
	for i := 0; i < 360; i += 360 / 12 {
		markX, markY := getAnalogWatchCoord(i)
		canvas.Add(gfx.NewLine(hourMark,
			centerX+markX*r/255,
			centerY+markY*r/255,
			centerX+markX*r/290,
			centerY+markY*r/290,
			r/16))
	}

	// Dot in the center.
	canvas.Add(gfx.NewCircle(centerDotColor, centerX, centerY, 6))

	// Handles, at an arbitrary position.
	hourHandle := gfx.NewLine(hourHandleColor, centerX, centerY, centerX, centerY, 10)
	canvas.Add(hourHandle)
	minuteHandle := gfx.NewLine(minuteHandleColor, centerX, centerY, centerX, centerY, 6)
	canvas.Add(minuteHandle)

	updateTime := func(hour, minute, second int) {
		hour = hour % 12
		hourX, hourY := getAnalogWatchCoord(hour*30 + minute/2)
		hourHandle.SetPosition(centerX+hourX*10/255, centerY+hourY*10/255, centerX+hourX*r/500, centerY+hourY*r/500)
		minuteX, minuteY := getAnalogWatchCoord(minute*6 + second/10)
		minuteHandle.SetPosition(centerX+minuteX*10/255, centerY+minuteY*10/255, centerX+minuteX*r/300, centerY+minuteY*r/300)
	}
	now := watchTime()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	updateTime(hour, minute, second)

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
	return NewView[T](eventWrapper, func(now time.Time) {
		// Update the watchface.
		if backlight > 0 {
			// Watch face is visible.
			newHour := now.Hour()
			newMinute := now.Minute()
			newSecond := now.Second()
			if hour != newHour || minute != newMinute || second != newSecond {
				hour = newHour
				minute = newMinute
				second = newSecond
				updateTime(hour, minute, second)
			}
		}
	})
}

// Return the -255..255 X and Y coordinates on a circle, starting at 12 o'clock
// going right around. The index is the number of degrees (0..359).
func getAnalogWatchCoord(index int) (x, y int) {
	switch {
	case index < 90:
		return int(watchCoord[index]), -int(watchCoord[89-index])
	case index < 180:
		return int(watchCoord[179-index]), int(watchCoord[index-90])
	case index < 270:
		return -int(watchCoord[index-180]), int(watchCoord[269-index])
	default: // index < 360
		return -int(watchCoord[359-index]), -int(watchCoord[index-270])
	}
}

// Table with precalculated sin/cos values.
// Python oneliner:
//
//	r=255; [round(math.sin(i/180*math.pi)*r) for i in range(90)]
var watchCoord = [...]uint8{
	0, 4, 9, 13, 18, 22, 27, 31, 35, 40, 44, 49, 53, 57, 62, 66, 70, 75, 79, 83, 87, 91, 96, 100, 104, 108, 112, 116, 120, 124, 127, 131, 135, 139, 143, 146, 150, 153, 157, 160, 164, 167, 171, 174, 177, 180, 183, 186, 190, 192, 195, 198, 201, 204, 206, 209, 211, 214, 216, 219, 221, 223, 225, 227, 229, 231, 233, 235, 236, 238, 240, 241, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 253, 254, 254, 254, 255, 255, 255,
}
