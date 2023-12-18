package main

// Watch firmware, mainly intended for the PineTime.

import (
	"strconv"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/style"
	"github.com/aykevl/tinygl/style/basic"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pixel"
	"tinygo.org/x/tinyfont/freesans"
)

var backlight = -1
var lastEvent time.Time

var screenTimeoutIndex uint8 = 1
var screenTimeouts = [...]time.Duration{
	3 * time.Second,
	5 * time.Second,
	10 * time.Second,
	15 * time.Second,
	30 * time.Second,
	60 * time.Second,
}

func main() {
	if board.Name == "simulator" {
		// Watch dimensions:
		// diagonal: 33mm, side: 23.3mm or 0.91 inch
		board.Simulator.WindowTitle = "GopherWatch"
		board.Simulator.WindowWidth = 240
		board.Simulator.WindowHeight = 240
		board.Simulator.WindowPPI = 261
		board.Simulator.WindowDrawSpeed = time.Second * 12 / 8e6 // 8MHz, 12bpp
	}

	println("start")
	board.Power.Configure()
	board.Buttons.Configure()
	err := InitBluetooth()
	if err != nil {
		println("could not configure Bluetooth:", err)
	}
	watch := MakeWatch(board.Display.Configure(), board.Display.ConfigureTouch())
	watch.run()
}

type Watch[T pixel.Color] struct {
	display    board.Displayer[T]
	touchInput board.TouchInput
	screen     *tinygl.Screen[T]
	views      *ViewManager[T]
}

func MakeWatch[T pixel.Color](display board.Displayer[T], touchInput board.TouchInput) Watch[T] {
	return Watch[T]{
		display:    display,
		touchInput: touchInput,
	}
}

func (w *Watch[T]) run() {
	var (
		black = pixel.NewColor[T](0, 0, 0)
		white = pixel.NewColor[T](255, 255, 255)
	)

	// Configure the screen.
	width, _ := w.display.Size()
	buf := pixel.NewImage[T](int(width), 32)
	scalePercent := board.Display.PPI() * 100 / 120
	scale := style.NewScale(scalePercent)
	println("scale:", board.Display.PPI(), "->", scale.Percent())
	w.screen = tinygl.NewScreen[T](w.display, buf, board.Display.PPI())
	basicTheme := basic.NewTheme(scale, w.screen)
	basicTheme.Tint = pixel.NewColor[T](144, 144, 144)
	w.views = &ViewManager[T]{
		screen: w.screen,
		Basic:  basicTheme,
	}
	w.views.Background = black
	w.views.Foreground = white
	lastEvent = time.Now()

	// Configure the accelerometer.
	var lastSensorUpdate time.Time
	board.Sensors.Configure(drivers.Acceleration | drivers.Temperature)
	updateSensors := func(now time.Time) {
		board.Sensors.Update(drivers.Temperature | drivers.Acceleration)
		lastSensorUpdate = now
	}
	updateSensors(watchTime())

	// Set up a UI.
	w.views.Push(w.createWatchFace(w.views))

	// Run the default watch face.
	for {
		now := watchTime()

		// Check whether we need to disable the screen.
		if backlight > 0 && time.Now().Sub(lastEvent) > screenTimeouts[screenTimeoutIndex] {
			// Going to enter sleep state.
			// First, clear all the views that might be running. Go back to the
			// homescreen (because that is what we'll show when awaking).
			for w.views.Len() > 1 {
				w.views.Pop()
			}
			w.enterSleep()
		}

		bl := backlight // backlight value _before_ calling Update()
		w.views.Update(now)
		w.screen.Update()
		w.readInputs()
		if bl < 0 {
			// Either we just started up, or we came out of sleep.
			setBacklight(board.Display.MaxBrightness())
		}
		if backlight == 0 {
			// Sleeping, so don't refresh so often.
			// TODO: use interrupts instead (both the button and the touchscreen
			// can be triggered via interrupts).
			time.Sleep(time.Second / 10)

			// Update infrequently when the screen is off.
			if now.Sub(lastSensorUpdate) > 5*time.Second {
				updateSensors(now)
			}
		} else {
			// Screen is on, so be faster.
			board.Display.WaitForVBlank(time.Second / 60)

			// Update frequently when the screen is on.
			if now.Sub(lastSensorUpdate) > time.Second/4 {
				updateSensors(now)
			}
		}
	}
}

func (w *Watch[T]) readInputs() {
	touchPoints := w.touchInput.ReadTouch()
	if len(touchPoints) != 0 {
		lastEvent = time.Now()
		w.screen.SetTouchState(touchPoints[0].X, touchPoints[0].Y)
	} else {
		w.screen.SetTouchState(-1, -1)
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
				w.exitSleep()
			} else {
				if w.views.Len() > 1 {
					// Not sleeping, so go back to the home screen.
					for w.views.Len() > 1 {
						w.views.Pop()
					}
				} else {
					// Already on the home screen, so turn off the screen.
					w.enterSleep()
				}
			}
		}
	}
}

func (w *Watch[T]) enterSleep() {
	// Shut down the backlight, which is of course a huge battery drain.
	setBacklight(0)
	w.display.Sleep(true)
}

func (w *Watch[T]) exitSleep() {
	backlight = -1
	lastEvent = time.Now()
	w.display.Sleep(false)
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

func (w *Watch[T]) createAppsView(views *ViewManager[T]) View[T] {
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
		"Rotate",
		"Screen timeout",
	})
	list.SetGrowable(1, 1)
	list.SetPadding(0, 8)
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
		case 4:
			// Rotate the screen by 180°.
			w.display.SetRotation((w.display.Rotation() + 2) % 4)
		case 5:
			views.Push(createScreenTimeoutView(views))
		}
	})
	return NewView[T](tinygl.NewVerticalScrollBox[T](header, list, nil), nil)
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
	width, _ := views.screen.Size()

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
	temperature := views.NewText("temp: ...")
	temperature.SetAlign(tinygl.AlignLeft)
	acceleration := views.NewText("accel: ...")
	acceleration.SetAlign(tinygl.AlignLeft)
	stepsText := views.NewText("steps: ...")
	stepsText.SetAlign(tinygl.AlignLeft)
	vbox := views.NewVBox(header, battery, voltage, temperature, acceleration, stepsText)
	wrapper := tinygl.NewEventBox[T](vbox)
	wrapper.SetEventHandler(func(event tinygl.Event, x, y int) {
		if event != tinygl.TouchTap {
			return
		}
		views.Pop()
	})

	var steps uint32 = 0xffff_ffff
	var ax, ay, az int32
	var temp int32 = 0x7fff_ffff

	// Create the view with the update callback.
	var lastTime time.Time
	return NewView[T](wrapper, func(now time.Time) {
		if now.Sub(lastTime) > time.Second/5 {
			lastTime = now

			// Update the UI with the battery status.
			state, microvolts, percent := board.Power.Status()
			batteryText := "battery: "
			if percent >= 0 {
				batteryText += strconv.Itoa(int(percent)) + "% "
			}
			if state != board.Discharging {
				batteryText += state.String()
			}
			battery.SetText(batteryText)
			voltage.SetText("voltage: " + formatVoltage(microvolts))
		}

		temp2 := board.Sensors.Temperature()
		if temp != temp2 {
			temp = temp2
			temperature.SetText("temp: " + strconv.Itoa(int(temp/1000)) + "C")
		}

		ax2, ay2, az2 := board.Sensors.Acceleration()
		if ax != ax2 || ay != ay2 || az != az2 {
			ax = ax2
			ay = ay2
			az = az2
			acceleration.SetText("accel: " + strconv.FormatInt(int64(ax/10000), 10) + " " + strconv.FormatInt(int64(ay/10000), 10) + " " + strconv.FormatInt(int64(az/10000), 10))
		}

		steps2 := board.Sensors.Steps()
		if steps != steps2 {
			steps = steps2
			stepsText.SetText("steps: " + strconv.Itoa(int(steps)))
		}
	})
}

func createScreenTimeoutView[T pixel.Color](views *ViewManager[T]) View[T] {
	var (
		lightblue = pixel.NewColor[T](64, 64, 255)
	)
	header := views.NewText("Screen timeout")
	header.SetBackground(lightblue)

	list := views.NewListBox([]string{
		"3s",
		"5s",
		"10s",
		"15s",
		"30s",
		"60s",
	})
	list.SetGrowable(1, 1)
	list.SetColumns(2)
	list.SetPadding(0, 16)
	list.Select(int(screenTimeoutIndex))
	list.SetAlign(tinygl.AlignCenter)
	list.SetEventHandler(func(event tinygl.Event, index int) {
		if event != tinygl.TouchTap {
			return
		}
		views.Pop() // go back to the homescreen after closing the view
		screenTimeoutIndex = uint8(index)
	})
	return NewView[T](views.NewVBox(header, list), nil)
}

func formatVoltage(microvolts uint32) string {
	microvolts += 500 // divided by 1000, so add 500 for correct rounding
	volts := strconv.Itoa(int(microvolts / 1000_000))
	decimals := strconv.Itoa(int(microvolts % 1000_000 / 1000))
	for len(decimals) < 2 {
		decimals = "0" + decimals
	}
	return volts + "." + decimals + "V"
}
