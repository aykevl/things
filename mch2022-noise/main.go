package main

import (
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ili9341"
)

const (
	width  = 320
	height = 240
)

func main() {
	println("starting...")

	machine.LCD_MODE.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LCD_MODE.Low()

	spi := machine.SPI2
	err := spi.Configure(machine.SPIConfig{
		Frequency: 40_000_000, // the ESP32 supports up to 40MHz SPI
		SCK:       18,
		SDO:       23,
		SDI:       35,
	})
	if err != nil {
		println("couldn't configure SPI:", err)
	}
	display := ili9341.NewSPI(spi, machine.LCD_DC, machine.SPI0_CS_LCD_PIN, machine.LCD_RESET)
	display.Configure(ili9341.Config{
		Rotation: ili9341.Rotation90,
	})

	const size = 4
	buffer := make([]uint16, width*height)
	for cycle := 0; ; cycle++ {
		start := time.Now()
		for y := 0; y < height; y += size {
			for x := 0; x < width; x += size {
				value := ledsgo.Noise3(
					uint32(x)<<4,
					uint32(y)<<4,
					uint32(start.UnixNano()>>21))
				c := ledsgo.PartyColors.ColorAt(value * 2)
				raw := makeColor(c.R, c.G, c.B)
				for x2 := x; x2 < x+size; x2++ {
					for y2 := y; y2 < y+size; y2++ {
						buffer[y2*width+x2] = raw
					}
				}
			}
		}
		draw := time.Now()
		display.DrawRGBBitmap(0, 0, buffer, width, height)
		if cycle%16 == 0 {
			end := time.Now()
			println("frame update:", draw.Sub(start).String())
			println("frame draw:  ", end.Sub(draw).String())
			println("frame:", end.Sub(start).String(), time.Second/end.Sub(start))
		}
	}
}

func handleError(err error) {
	if err != nil {
		for {
			println("error:", err)
			time.Sleep(time.Second)
		}
	}
}

func makeColor(r, g, b uint8) uint16 {
	r = gamma8[r]
	g = gamma8[g]
	b = gamma8[b]
	c := uint16(r&0xF8)<<8 +
		uint16(g&0xFC)<<3 +
		uint16(b&0xF8)>>3
	c = c>>8 | c<<8 // swap endianness
	return c
}

// Gamma brightness lookup table <https://victornpb.github.io/gamma-table-generator>
// gamma = 0.40 steps = 256 range = 0-255
var gamma8 = [...]uint8{
	0, 28, 37, 43, 48, 53, 57, 61, 64, 67, 70, 73, 75, 78, 80, 82,
	84, 86, 88, 90, 92, 94, 96, 97, 99, 101, 102, 104, 105, 107, 108, 110,
	111, 113, 114, 115, 117, 118, 119, 120, 122, 123, 124, 125, 126, 127, 129, 130,
	131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146,
	147, 148, 149, 149, 150, 151, 152, 153, 154, 155, 155, 156, 157, 158, 159, 160,
	160, 161, 162, 163, 164, 164, 165, 166, 167, 167, 168, 169, 170, 170, 171, 172,
	173, 173, 174, 175, 175, 176, 177, 177, 178, 179, 179, 180, 181, 182, 182, 183,
	183, 184, 185, 185, 186, 187, 187, 188, 189, 189, 190, 190, 191, 192, 192, 193,
	194, 194, 195, 195, 196, 197, 197, 198, 198, 199, 199, 200, 201, 201, 202, 202,
	203, 203, 204, 205, 205, 206, 206, 207, 207, 208, 208, 209, 209, 210, 211, 211,
	212, 212, 213, 213, 214, 214, 215, 215, 216, 216, 217, 217, 218, 218, 219, 219,
	220, 220, 221, 221, 222, 222, 223, 223, 224, 224, 225, 225, 226, 226, 227, 227,
	228, 228, 229, 229, 230, 230, 230, 231, 231, 232, 232, 233, 233, 234, 234, 235,
	235, 235, 236, 236, 237, 237, 238, 238, 239, 239, 240, 240, 240, 241, 241, 242,
	242, 243, 243, 243, 244, 244, 245, 245, 246, 246, 246, 247, 247, 248, 248, 248,
	249, 249, 250, 250, 251, 251, 251, 252, 252, 253, 253, 253, 254, 254, 255, 255,
}
