package main

import (
	"machine"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/pixel"
	"tinygo.org/x/drivers/st7789"
)

func testTearing[T pixel.Color](display tinygl.Displayer, screen *tinygl.Screen[T]) {
	var (
		white = pixel.NewColor[T](255, 255, 255)
		black = pixel.NewColor[T](0, 0, 0)
		green = pixel.NewColor[T](0, 255, 0)
		red   = pixel.NewColor[T](255, 0, 0)
	)
	const (
		size  = 32
		speed = 4
	)
	width, height := screen.Size()
	canvas := gfx.NewCanvas(black, 32, 32)
	rect := canvas.CreateRect(0, 0, width, size, white)
	teIndicator := canvas.CreateRect(0, 0, 8, 8, red)
	screen.SetChild(canvas)

	d := display.(*st7789.Device)
	d.Command(0x35) // TEON
	d.Data(0)       // M=0

	avoidTearing := false
	te := machine.GPIO9
	te.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	verticalDirection := speed
	start := time.Now()
	for cycle := 0; ; cycle++ {
		// Read keyboard inputs.
		board.Buttons.ReadInput()
		for {
			// Read keyboard events.
			event := board.Buttons.NextEvent()
			if event == board.NoKeyEvent {
				break
			}
			if event.Pressed() {
				switch event.Key() {
				case board.KeyA:
					avoidTearing = !avoidTearing
					if avoidTearing {
						teIndicator.SetColor(green)
					} else {
						teIndicator.SetColor(red)
					}
				case board.KeyB, board.KeyEscape:
					return
				}
			}
		}

		x, y, _, h := rect.Bounds()
		if y <= 0 {
			verticalDirection = speed
		} else if y+h >= height {
			verticalDirection = -speed
		}
		rect.Move(x, y+verticalDirection)
		screen.Update()

		now := time.Now()
		if cycle%64 == 0 {
			println("rendering took:", now.Sub(start).String())
		}
		sleep := time.Second/60 - now.Sub(start)
		time.Sleep(sleep)
		start = now

		if avoidTearing {
			// Wait until the TE line goes high.
			// When it is high, the framebuffer is not currently used.
			for te.Get() == false {
			}
		}
	}
}
