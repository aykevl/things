package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/style"
	"github.com/aykevl/tinygl/style/basic"
	"tinygo.org/x/drivers/pixel"
)

func main() {
	println("start")
	if board.Name == "simulator" {
		// Use the configuration for the Gopher Badge.
		board.Simulator.WindowWidth = 320
		board.Simulator.WindowHeight = 240
		board.Simulator.WindowPPI = 166
	}

	board.Buttons.Configure()
	run(board.Display.Configure(), board.Display.ConfigureTouch())
}

func run[T pixel.Color](display board.Displayer[T], touchInput board.TouchInput) {
	// Determine size and scale of the screen.
	width, height := display.Size()
	scalePercent := board.Display.PPI() * 100 / 120

	// Initialize the screen.
	buf := pixel.NewImage[T](int(width), int(height)/4)
	screen := tinygl.NewScreen[T](display, buf, board.Display.PPI())
	theme := basic.NewTheme(style.NewScale(scalePercent), screen)
	println("scale:", scalePercent, "=>", theme.Scale.Percent())

	// Create badge homescreen.
	header := theme.NewText("Hello world!")
	header.SetBackground(pixel.NewColor[T](255, 0, 0))
	header.SetColor(pixel.NewColor[T](255, 255, 255))
	listbox := theme.NewListBox([]string{
		"Noise",
		"Mandelbrot",
		"Display test colors",
		"Touch test",
		"Tearing test",
		"Sensors",
		"LEDs",
		"Images",
	})
	listbox.SetGrowable(0, 1) // listbox fills the rest of the screen
	listbox.Select(0)         // focus the first element
	home := theme.NewVBox(header, listbox)

	// Handle touch events in the listbox.
	listbox.SetEventHandler(func(event tinygl.Event, index int) {
		if event == tinygl.TouchTap {
			runApp(index, display, screen, home, touchInput)
		}
	})

	// Show screen.
	screen.SetChild(home)
	screen.Update()
	board.Display.SetBrightness(board.Display.MaxBrightness())

	for {
		// TODO: wait for input instead of polling
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
			case board.KeyUp:
				index := listbox.Selected() - 1
				if index < 0 {
					index = listbox.Len() - 1
				}
				listbox.Select(index)
			case board.KeyDown:
				index := listbox.Selected() + 1
				if index >= listbox.Len() {
					index = 0
				}
				listbox.Select(index)
			case board.KeyEnter, board.KeyA:
				runApp(listbox.Selected(), display, screen, home, touchInput)
			}
		}

		// Handle touch inputs.
		touches := touchInput.ReadTouch()
		if len(touches) > 0 {
			screen.SetTouchState(touches[0].X, touches[0].Y)
		} else {
			screen.SetTouchState(-1, -1)
		}

		screen.Update()
		time.Sleep(time.Second / 30)
	}
}

func runApp[T pixel.Color](index int, display board.Displayer[T], screen *tinygl.Screen[T], home *tinygl.VBox[T], touchInput board.TouchInput) {
	switch index {
	case 0:
		println("starting noise")
		noise(display, screen)
	case 1:
		println("starting Mandelbrot")
		mandelbrot(screen)
	case 2:
		println("starting display test colors")
		testColors(screen)
	case 3:
		println("starting touch test")
		testTouch(screen, touchInput)
	case 4:
		println("starting tearing test")
		testTearing(screen, touchInput)
	case 5:
		println("starting sensors")
		showSensors(screen)
	case 6:
		println("toggle LEDs")
		toggleLEDs()
	case 7:
		println("starting images")
		showImages(screen)
	}

	// The app used a different root element. Restore the homescreen.
	screen.SetChild(home)
}
