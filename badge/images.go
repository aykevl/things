package main

import (
	_ "embed"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/image"
	"tinygo.org/x/drivers/pixel"
)

// Image created using ImageMagic:
//
//	convert tinygo-logo.png -background white -alpha remove -alpha off -resize 320x240 -depth 5 tinygo-logo-small.qoi
//
// (Remove alpha channel, resize to fit the screen, remove unnecessary bits,
// convert to QOI format).
//
//go:embed images/tinygo-logo-small.qoi
var tinygoLogo string

func showImages[T pixel.Color](screen *tinygl.Screen[T]) {
	// Load image.
	img, err := image.NewQOI[T](tinygoLogo)
	if err != nil {
		println("could not load image:", err)
		return
	}
	mainImage := tinygl.NewImage[T](img)
	mainImage.SetBackground(pixel.NewColor[T](255, 255, 255))
	screen.SetChild(mainImage)

	// Show screen.
	screen.Update()

	for {
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
				case board.KeyB, board.KeyEscape:
					return
				}
			}
		}

		time.Sleep(time.Second / 10)
	}
}
