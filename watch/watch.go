package main

// Watch firmware, mainly intended for the PineTime.

import (
	"strconv"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/pixel"
	"github.com/aykevl/tinygl/style"
	"github.com/aykevl/tinygl/style/basic"
	"tinygo.org/x/tinyfont/freesans"
)

var backlight = -1
var lastEvent time.Time

const screenTimeout = 5 * time.Second

func main() {
	if board.Name == "simulator" {
		// Watch dimensions:
		// diagonal: 33mm, side: 23.3mm or 0.91 inch
		board.Simulator.WindowWidth = 240
		board.Simulator.WindowHeight = 240
		board.Simulator.WindowPPI = 261
		board.Simulator.WindowDrawSpeed = time.Second * 16 / 8e6 // 8MHz, 16bpp
	}

	println("start")
	board.Power.Configure()
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
	scalePercent := board.Display.PPI() * 100 / 120
	scale := style.NewScale(scalePercent)
	println("scale:", board.Display.PPI(), "->", scale.Percent())
	screen := tinygl.NewScreen(display, buf, board.Display.PPI())
	views := &ViewManager[T]{
		screen: screen,
		Basic:  basic.NewTheme(scale, screen),
	}
	views.Background = black
	views.Foreground = white
	lastEvent = time.Now()

	// Helpers for sleep modes.
	enterSleep := func() {
		// Shut down the backlight, which is of course a huge battery drain.
		setBacklight(0)
		display.Sleep(true)
	}
	exitSleep := func() {
		backlight = -1
		lastEvent = time.Now()
		display.Sleep(false)
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
			// There is only one button on the watch, so just check whether it's
			// pressed.
			if event.Pressed() {
				lastEvent = time.Now()
				if backlight == 0 {
					// Sleeping, so wake up the screen.
					exitSleep()
				} else {
					if views.Len() > 1 {
						// Not sleeping, so go back to the home screen.
						for views.Len() > 1 {
							views.Pop()
						}
					} else {
						// Already on the home screen, so turn off the screen.
						enterSleep()
					}
				}
			}
		}
	})

	// Set up a UI.
	hello := tinygl.NewText(&freesans.Regular24pt7b, white, black, "00:00")
	eventWrapper := tinygl.NewEventBox[T](hello)
	var minute int = -1
	views.Push(NewView[T](eventWrapper, func(now time.Time) {
		// Update the watchface.
		if backlight > 0 {
			// Watch face is visible.
			newMinute := now.Minute()
			if minute != newMinute {
				minute = newMinute
				hello.SetText(formatTime(now.Hour(), minute))
			}
		}
	}))
	eventWrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event == tinygl.TouchTap {
			if backlight == 0 {
				// Tapped on a sleeping watch.
				// Awake the screen.
				exitSleep()
			} else {
				// Regular tap on the clock.
				// TODO: detect gesture (for example, swipe upwards) to make it
				// harder to accidentally get in the settings menu.
				views.Push(createAppsView(views))
			}
		}
	})

	// Run the default watch face.
	for {
		now := watchTime()

		// Check whether we need to disable the screen.
		if backlight > 0 && time.Now().Sub(lastEvent) > screenTimeout {
			// Going to enter sleep state.
			// First, clear all the views that might be running. Go back to the
			// homescreen (because that is what we'll show when awaking).
			for views.Len() > 1 {
				views.Pop()
			}
			enterSleep()
		}

		bl := backlight // backlight value _before_ calling Update()
		views.Update(now)
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
		backlight = level
		board.Display.SetBrightness(level)
	}
}

func createAppsView[T pixel.Color](views *ViewManager[T]) View[T] {
	// Constants used in this function.
	var (
		lightblue = pixel.NewColor[T](64, 64, 255)
	)

	// Create the settings UI.
	header := views.NewText("Settings")
	header.SetBackground(lightblue)
	list := views.NewListBox([]string{
		"Back",
		"Set time",
		"Touch test",
		"Sensors",
	})
	list.SetGrowable(1, 1)
	list.SetEventHandler(func(event tinygl.Event, index int) {
		if event != tinygl.TouchTap {
			return
		}
		views.Pop() // go back to the homescreen after closing the view
		switch index {
		case 0:
			// Nothing to do, just go back to the homescreen.
		case 1:
			views.Push(createClockAdjustView(views))
		case 2:
			views.Push(createTouchTestView(views))
		case 3:
			views.Push(createSensorsView(views))
		}
	})
	return NewView[T](views.NewVBox(header, list), nil)
}

// Create view to adjust the time on the watch.
func createClockAdjustView[T pixel.Color](views *ViewManager[T]) View[T] {
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
		if event != tinygl.TouchTap {
			return
		}
		if x < width/2 {
			hour = (hour + 1) % 24
		} else {
			minute = (minute + 1) % 60
		}
		text.SetText(formatTime(hour, minute))
	})
	sub.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event != tinygl.TouchTap {
			return
		}
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
		if event != tinygl.TouchTap {
			return
		}
		// Update time and close this view.
		oldTime := watchTime()
		diff := time.Duration(hour-oldTime.Hour()) * time.Hour
		diff += time.Duration(minute-oldTime.Minute()) * time.Minute
		diff -= time.Duration(oldTime.Nanosecond())
		adjustTime(diff)
		views.Pop()
	})

	return NewView[T](box, nil)
}

// view the values of a few sensors.
func createSensorsView[T pixel.Color](views *ViewManager[T]) View[T] {
	// Create the UI.
	header := views.NewText("Sensors")
	header.SetBackground(pixel.NewColor[T](0, 0, 255))
	header.SetColor(pixel.NewColor[T](255, 255, 255))
	battery := views.NewText("battery: ...")
	battery.SetAlign(tinygl.AlignLeft)
	voltage := views.NewText("voltage: ...")
	voltage.SetAlign(tinygl.AlignLeft)
	vbox := views.NewVBox(header, battery, voltage)
	wrapper := tinygl.NewEventBox[T](vbox)
	wrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event != tinygl.TouchTap {
			return
		}
		views.Pop()
	})

	// Create the view with the update callback.
	var lastTime time.Time
	return NewView[T](wrapper, func(now time.Time) {
		if now.Sub(lastTime) > time.Second/5 {
			lastTime = now
			// Update the UI with values.
			state, microvolts := board.Power.Status()
			battery.SetText("battery: " + state.String())
			voltage.SetText("voltage: " + formatVoltage(microvolts))
		}
	})
}

func formatVoltage(microvolts uint32) string {
	volts := strconv.Itoa(int(microvolts / 1000_000))
	decimals := strconv.Itoa(int(microvolts % 1000_000 / 10_000))
	for len(decimals) < 2 {
		decimals = "0" + decimals
	}
	return volts + "." + decimals + "V"
}
