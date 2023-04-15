package main

// Watch firmware, mainly intended for the PineTime.

import (
	"strconv"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/pixel"
	"tinygo.org/x/tinyfont/freesans"
)

var backlight = -1
var lastEvent time.Time

const screenTimeout = 5 * time.Second

func main() {
	// Watch dimensions:
	// diagonal: 33mm, side: 23.3mm or 0.91 inch
	board.Simulator.WindowWidth = 240
	board.Simulator.WindowHeight = 240
	board.Simulator.WindowPPI = 261

	println("start")
	board.Buttons.Configure()
	run(board.Display.Configure(), board.Display.ConfigureTouch())
}

func run[T pixel.Color](display board.Displayer[T], touchInput board.TouchInput) {
	var (
		black = pixel.NewColor[T](0, 0, 0)
		white = pixel.NewColor[T](255, 255, 255)
	)

	// Configure the screen.
	width, _ := display.Size()
	buf := make([]T, width*32)
	screen := tinygl.NewScreen(display, buf, board.Display.PPI())
	lastEvent = time.Now()

	// This is a stack of views that can be added on top and popped from when
	// going back to the previous view.
	var views []tinygl.Object[T]
	pushView := func(view tinygl.Object[T]) {
		views = append(views, view)
		screen.SetChild(view)
	}
	popView := func() {
		views = views[:len(views)-1]
		screen.SetChild(views[len(views)-1])
	}

	// Helper to get out of sleep mode (turn on the display etc).
	exitSleep := func() {
		backlight = -1
		lastEvent = time.Now()
	}

	// Configure input events.
	screen.SetUpdateCallback(func(screen *tinygl.Screen[T]) {
		touchPoints := touchInput.ReadTouch()
		if len(touchPoints) != 0 {
			lastEvent = time.Now()
			screen.SetTouchState(touchPoints[0].X, touchPoints[0].Y)
		} else {
			screen.SetTouchState(-1, -1)
		}

		board.Buttons.ReadInput()
		for {
			event := board.Buttons.NextEvent()
			if event == board.NoKeyEvent {
				break
			}
			if event.Pressed() {
				// There is only one button on the watch.
				lastEvent = time.Now()
				if backlight == 0 {
					exitSleep()
				}
			}
		}
	})

	// Set up a UI.
	hello := tinygl.NewText(&freesans.Regular24pt7b, white, black, "00:00")
	eventWrapper := tinygl.NewEventBox[T](hello)
	pushView(eventWrapper)
	eventWrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event == tinygl.TouchTap {
			if backlight == 0 {
				// Tapped on a sleeping watch.
				// Awake the screen.
				exitSleep()
			} else {
				// Regular tap on the clock.
				pushView(runSettings[T](popView))
			}
		}
	})

	// Run the default watch face.
	var minute int = -1
	for {
		now := watchTime()

		// Check whether we need to disable the screen.
		if backlight > 0 && time.Now().Sub(lastEvent) > screenTimeout {
			// Going to enter sleep state.
			// First, clear all the views that might be running. Go back to the
			// homescreen (because that is what we'll show when awaking).
			for len(views) > 1 {
				popView()
			}
			// Shut down the backlight, which is of course a huge battery drain.
			setBacklight(0)
		}

		// Update the watchface.
		if len(views) == 1 {
			// Watch face is visible.
			newMinute := now.Minute()
			if minute != newMinute {
				minute = newMinute
				hello.SetText(formatTime(now.Hour(), minute))
			}
		} else {
			// Some other view is laid over the watch face.
			minute = -1
		}

		bl := backlight // backlight value _before_ calling Update()
		screen.Update()
		if bl < 0 {
			// Either we just started up, or we came out of sleep.
			setBacklight(board.Display.MaxBrightness())
		}
		if backlight == 0 {
			// Sleeping, so don't refresh so often.
			// TODO: use interrupts instead (both the button and the touchscreen
			// can be triggered via interrupts).
			time.Sleep(time.Second / 10)
		} else {
			// Not sleeping, so be faster.
			board.Display.WaitForVBlank(time.Second / 60)
		}
	}
}

// Format a time without using time.Format.
func formatTime(hour, minute int) string {
	h := strconv.Itoa(hour)
	if len(h) == 1 {
		h = "0" + h
	}
	m := strconv.Itoa(minute)
	if len(m) == 1 {
		m = "0" + m
	}
	return h + ":" + m
}

// Set the backlight to the given level. This is a no-op if it wouldn't change
// the backlight level.
func setBacklight(level int) {
	if backlight != level {
		println("change backlight level:", level)
		if backlight <= 0 {
			// Allow the LCD to refresh the display completely before showing
			// it. This is especially important when coming out of reset to
			// avoid a white flash of the default memory contents.
			time.Sleep(20 * time.Millisecond)
		}
		backlight = level
		board.Display.SetBrightness(level)
	}
}

// Settings view. Actually, this only allows setting the time.
func runSettings[T pixel.Color](closeView func()) tinygl.Object[T] {
	// Constants used in this function.
	var (
		green = pixel.NewColor[T](32, 255, 0)
		red   = pixel.NewColor[T](255, 0, 0)
		black = pixel.NewColor[T](0, 0, 0)
		white = pixel.NewColor[T](255, 255, 255)
	)
	w, _ := board.Display.Size()
	width := int(w)

	// Configure UI.
	start := watchTime()
	hour := start.Hour()
	minute := start.Minute()
	addText := tinygl.NewText(&freesans.Regular24pt7b, green, black, "+   +")
	subText := tinygl.NewText(&freesans.Regular24pt7b, red, black, "-   -")
	add := tinygl.NewEventBox[T](addText)
	sub := tinygl.NewEventBox[T](subText)
	text := tinygl.NewText(&freesans.Regular24pt7b, white, black, formatTime(hour, minute))
	textWrapper := tinygl.NewEventBox[T](text)
	add.SetGrowable(1, 1)
	sub.SetGrowable(1, 1)
	box := tinygl.NewVBox[T](black, add, textWrapper, sub)

	// Add event handlers.
	add.SetEventHandler(func(event tinygl.Event, x, y int) {
		if x < width/2 {
			hour = (hour + 1) % 24
		} else {
			minute = (minute + 1) % 60
		}
		text.SetText(formatTime(hour, minute))
	})
	sub.SetEventHandler(func(event tinygl.Event, x, y int) {
		if x < width/2 {
			hour--
			if hour < 0 {
				hour = 23
			}
		} else {
			minute--
			if minute < 0 {
				minute = 59
			}
		}
		text.SetText(formatTime(hour, minute))
	})
	textWrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		// Update time and close this view.
		oldTime := watchTime()
		diff := time.Duration(hour-oldTime.Hour()) * time.Hour
		diff += time.Duration(minute-oldTime.Minute()) * time.Minute
		diff -= time.Duration(oldTime.Nanosecond())
		adjustTime(diff)
		closeView()
	})

	return box
}
