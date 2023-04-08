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
	run(board.Display.Configure())
}

func run[T pixel.Color](display board.Displayer[T]) {
	// Determine size and scale of the screen.
	width, height := display.Size()
	physicalWidth, _ := board.Display.PhysicalSize()
	scalePercent := int(width) * 21 / physicalWidth

	// Initialize the screen.
	buf := make([]T, int(width)*int(height)/4)
	screen := tinygl.NewScreen(display, buf)
	theme := basic.NewTheme(style.NewScale(scalePercent), screen)
	println("scale:", scalePercent, "=>", theme.Scale.Percent())

	// Create badge homescreen.
	header := theme.NewText("Hello world!")
	header.SetBackground(pixel.NewColor[T](255, 0, 0))
	header.SetColor(pixel.NewColor[T](255, 255, 255))
	listbox := theme.NewListBox([]string{"Noise", "Mandelbrot", "Settings"})
	listbox.SetGrowable(0, 1) // listbox fills the rest of the screen
	listbox.Select(0)         // focus the first element
	home := theme.NewVBox(header, listbox)

	// Show screen.
	screen.SetChild(home)
	screen.Update()

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
				index := listbox.Selected()
				switch index {
				case 0:
					println("starting noise")
					// TODO: I have a local implementation but it relies on DMA
					// support in the display for good performance.
					//noise(display, buf)
				case 1:
					println("starting Mandelbrot")
					mandelbrot(display, buf)
				}
			}
		}
		screen.Update()

		time.Sleep(time.Second / 30)
	}
}
