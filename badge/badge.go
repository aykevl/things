package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/pixel"
	"github.com/aykevl/tinygl/style"
	"github.com/aykevl/tinygl/style/basic"
)

func main() {
	println("start")
	board.Buttons.Configure()
	run(board.Display.Configure(), board.Display.ConfigureTouch())
}

func run[T pixel.Color](display board.Displayer[T], touchInput board.TouchInput) {
	// Determine size and scale of the screen.
	width, height := display.Size()
	scalePercent := board.Display.PPI() * 100 / 120

	// Initialize the screen.
	buf := pixel.NewImage[T](int(width), int(height)/4)
	screen := tinygl.NewScreen(display, buf, board.Display.PPI())
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
		// TODO: I have a local implementation but it relies on DMA
		// support in the display for good performance.
		//noise(display, buf)
	case 1:
		println("starting Mandelbrot")
		mandelbrot(display, screen.Buffer())
	case 2:
		println("starting display test colors")
		testColors(display, screen.Buffer())
	case 3:
		println("starting touch test")
		testTouch(screen, touchInput)
	case 4:
		println("starting tearing test")
		testTearing(display, screen, touchInput)
	case 5:
		println("starting sensors")
		showSensors(screen)
	}

	// Some apps use the same screen and set a different root
	// element. Restore the previous root.
	screen.SetChild(home)

	// Some apps draw directly on the screen. In that case, we need
	// to repaint the entire screen.
	home.RequestUpdate()
}
