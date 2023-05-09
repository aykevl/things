package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl/pixel"
)

func testColors[T pixel.Color](display board.Displayer[T], buf pixel.Image[T]) {
	width, height := display.Size()

	// Draw the test colors.
	img := buf.Rescale(int(width), 1)
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			gray := uint8(x * 255 / int(width))
			var c T
			switch y * 14 / int(height) {
			case 0, 1:
				c = pixel.NewColor[T](gray, 0, 0) // red
			case 3, 4:
				c = pixel.NewColor[T](0, gray, 0) // green
			case 6, 7:
				c = pixel.NewColor[T](0, 0, gray) // blue
			case 9, 10:
				c = pixel.NewColor[T](gray, gray, gray)
			case 12, 13:
				r := gamma_lut[uint8(255-gray)]
				g := gamma_lut[uint8(gray)]
				c = pixel.NewColor[T](r, g, 0)
			}
			img.Set(x, 0, c)
		}
		display.DrawRGBBitmap8(0, int16(y), img.RawBuffer(), width, 1)
	}

	// Wait for back button.
	for {
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

		// Make sure the display stays alive (in the simulator).
		display.Display()
		time.Sleep(time.Second / 30)
	}
}

// Gamma brightness lookup table <https://victornpb.github.io/gamma-table-generator>
// gamma = 0.45 steps = 256 range = 0-255
var gamma_lut = [256]uint8{
	0, 21, 29, 35, 39, 43, 47, 51, 54, 57, 59, 62, 64, 67, 69, 71,
	73, 75, 77, 79, 81, 83, 85, 86, 88, 90, 91, 93, 94, 96, 97, 99,
	100, 102, 103, 104, 106, 107, 108, 110, 111, 112, 113, 114, 116, 117, 118, 119,
	120, 121, 122, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136,
	137, 138, 139, 140, 141, 142, 143, 143, 144, 145, 146, 147, 148, 149, 150, 150,
	151, 152, 153, 154, 155, 156, 156, 157, 158, 159, 160, 160, 161, 162, 163, 164,
	164, 165, 166, 167, 167, 168, 169, 170, 170, 171, 172, 173, 173, 174, 175, 175,
	176, 177, 178, 178, 179, 180, 180, 181, 182, 182, 183, 184, 184, 185, 186, 186,
	187, 188, 188, 189, 190, 190, 191, 192, 192, 193, 193, 194, 195, 195, 196, 197,
	197, 198, 198, 199, 200, 200, 201, 201, 202, 203, 203, 204, 204, 205, 206, 206,
	207, 207, 208, 208, 209, 210, 210, 211, 211, 212, 212, 213, 214, 214, 215, 215,
	216, 216, 217, 217, 218, 219, 219, 220, 220, 221, 221, 222, 222, 223, 223, 224,
	224, 225, 225, 226, 227, 227, 228, 228, 229, 229, 230, 230, 231, 231, 232, 232,
	233, 233, 234, 234, 235, 235, 236, 236, 237, 237, 238, 238, 239, 239, 240, 240,
	241, 241, 242, 242, 242, 243, 243, 244, 244, 245, 245, 246, 246, 247, 247, 248,
	248, 249, 249, 250, 250, 250, 251, 251, 252, 252, 253, 253, 254, 254, 255, 255,
}
