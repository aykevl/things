// Package hub75 implements the hub75 protocol, as used in various "LED matrix"
// displays.
package hub75

// This is a driver for the common hub75 panels you can find on AliExpress etc.
//
// Background:
//   https://bikerglen.com/projects/lighting/led-panel-1up/
//   http://www.rayslogic.com/propeller/Programming/AdafruitRGB/AdafruitRGB.htm
// Datasheets:
//   http://www.dlfushi.com/uploads/201801/5a5427980fd8d.pdf
// Pins:
//   A, B, C, D:
//     Address line for rows, 0b0000 means topmost row is enabled and 0b1111
//     means bottom row is enabled.
//     It appears that these rows must be switched frequently or the display
//     won't show anything.
//   OE:
//     Output Enable. High means the display is completely dark, low means the
//     current row lights up with the data latched from the shift register.
//   Latch:
//     Also known as strobe. Copy the contents of the shift registers to the
//     output. It is normally low and it needs to be pulled high to update.
//
// This driver tries to use all available hardware on a chip to make the screen
// update as smoothly as possible while using as little CPU power in the process
// as possible. This means that it heavily relies on DMA and interrupts.
//
// This driver uses the following tricks to get this high performance:
//   * It uses binary coded modulation for the individual PWM levels. This means
//     that for 8 bits of color depth, only 8 brightness levels need to be sent
//     instead of all 255. Also, it scales much better so that getting to 11
//     bits isn't really difficult: the difficult part is only turning on the
//     screen for a very short time.
//   * The screen is turned on and off (PWM) while the next data buffer is being
//     sent to the display.
//   * Everything is interrupt driven, so that driving the screen takes up only
//     part of the CPU (thanks to DMA and timers). The rest can be used for
//     rendering the next frame in the animation.
//
// More precisely, this driver performs the following steps for each row of
// data:
//  1. The first iteration, it just starts sending the first buffer.
//  2. Once the buffer is sent, it triggers the latch, configures the correct
//     pin mux (ABCD) for this row, and starts a timer to enable/disable the
//     screen using the OE pin.
//  3. Once the timer is started (and in parallel of enabling/disabling the
//     screen), it queues up the next DMA buffer for SPI.
//  4. Once the timer has finished and the DMA buffer has been sent over SPI
//     (either of them can be the last), repeat from step 2.

import (
	"image/color"
	"machine"
	"runtime/volatile"
	"unsafe"
)

type Device struct {
	chipSpecificSettings
	a                 machine.Pin
	b                 machine.Pin
	c                 machine.Pin
	d                 machine.Pin
	oe                machine.Pin      // output enable pin
	lat               machine.Pin      // latch pin
	colorBit          uint8            // 0..7: which bit is currently drawn using binary coded modulation
	row               uint8            // 0..15: the row that is currently selected in the pin mux
	running           bool             // true if the driver is currently running (handling interrupts etc.)
	fullRefreshes     uint             // counter for the number of full refreshes, useful for statistics
	spiReady          uint8            // 1 when the SPI is finished sending, 0 otherwise (signaled from the interrupt)
	timerReady        uint8            // 1 when the timer has expried, 0 otherwise (signaled from the interrupt)
	framebuf          [3][32][32]uint8 // contains RGB data to be sent to the screen with the next call to Display()
	displayBitstrings [16][8][]uint8   // data that can be directly sent over SPI using DMA
	brightness        uint32           // at least 1, higher means brighter screen but slower updates
}

var display *Device

// New returns a new HUB75 driver. This is a singleton, don't attempt to use
// more than one.
func New(dataPin, clockPin, latPin, oePin, aPin, bPin, cPin, dPin machine.Pin) *Device {
	d := &Device{
		a:          aPin,
		b:          bPin,
		c:          cPin,
		d:          dPin,
		oe:         oePin,
		lat:        latPin,
		brightness: 1, // must be at least 1
	}

	if display != nil {
		panic("trying to instantiate more than one hub75 driver")
	}
	display = d

	d.a.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.b.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.c.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.d.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.oe.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.lat.Configure(machine.PinConfig{Mode: machine.PinOutput})

	d.a.Low()
	d.b.Low()
	d.c.Low()
	d.d.Low()
	d.lat.High()

	d.configureChip(dataPin, clockPin)

	return d
}

// FullRefreshes returns the number of full screen refreshes (all rows + all
// brightness levels) since this driver was started.
func (d *Device) FullRefreshes() uint {
	return d.fullRefreshes
}

// SetPixel updates the pixel RGB values at index x, y.
func (d *Device) SetPixel(x int16, y int16, c color.RGBA) {
	d.framebuf[0][x][y] = c.R
	d.framebuf[1][x][y] = c.G
	d.framebuf[2][x][y] = c.B
}

// flush copies the data in the frame buffer to the output bit strings that can
// be sent over SPI.
//go:nobounds
func (d *Device) flush() {
	// Make sure all bitstrings are present.
	for row := 0; row < 16; row++ {
		for bit := 0; bit < 8; bit++ {
			if d.displayBitstrings[row][bit] == nil {
				d.displayBitstrings[row][bit] = make([]uint8, 24)
			}
		}
	}

	for row := uint(0); row < 32; row++ {
		for colorIndex := 2; colorIndex >= 0; colorIndex-- {
			for bit := uint(0); bit < 8; bit++ {
				for colByte := uint(0); colByte < 4; colByte++ {
					// Unroll this loop for slightly higher performance.
					c := uint32(0)
					word := *(*uint32)(unsafe.Pointer(&d.framebuf[colorIndex][row][colByte*8+0]))
					word >>= bit
					c |= (word & (1 << 0)) << 7
					c |= (word & (1 << 8)) >> 2
					c |= (word & (1 << 16)) >> 11
					c |= (word & (1 << 24)) >> 20
					word = *(*uint32)(unsafe.Pointer(&d.framebuf[colorIndex][row][colByte*8+4]))
					word >>= bit
					c |= (word & (1 << 0)) << 3
					c |= (word & (1 << 8)) >> 6
					c |= (word & (1 << 16)) >> 15
					c |= (word & (1 << 24)) >> 24
					bitstringIndex := colByte + uint(2-colorIndex)*8
					if (row % 32) < 16 {
						bitstringIndex += 4
					}
					d.displayBitstrings[row%16][bit][bitstringIndex] = uint8(c)
				}
			}
		}
	}
}

// Display sends the buffer (if any) to the screen.
func (d *Device) Display() error {
	// Update the bitstrings that are sent over SPI.
	// TODO: perhaps we need double buffering here?
	d.flush()

	// Check if the driver is already running, and start it if it is not.
	if !d.running {
		d.running = true
		// Start by sending the first bitstring over SPI, and pretend that the
		// timer of the previous buffer has already been finished.
		d.timerReady = 1
		d.startTransfer()
	}
	return nil
}

// sendNext sends the next update. It should be called when the previous timer
// and SPI transfer have both finished.
func (d *Device) sendNext() {
	// Send the latch signal.
	// This means that everything that was shifted into the shift register will
	// now be set as the output value of the shift register.
	d.lat.High()
	d.lat.Low()

	// Update the row selection to match the current row.
	d.a.Set(d.row&0x01 != 0)
	d.b.Set(d.row&0x02 != 0)
	d.c.Set(d.row&0x04 != 0)
	d.d.Set(d.row&0x08 != 0)

	// Start the 'output enable' timer.
	d.startOutputEnableTimer()

	// Switch to the next row and possibly next color bit level.
	d.row = (d.row + 1) % 16
	if d.row == 0 {
		d.colorBit++
		if d.colorBit >= 8 {
			d.colorBit = 0
			d.fullRefreshes++
		}
	}

	// Start the next SPI transaction.
	d.startTransfer()
}

// handleTimerEvent is called from the timer interrupt, once the screen has been
// enabled and disabled again.
func (d *Device) handleTimerEvent() {
	// Start the next cycle if the SPI buffer has been sent.
	if volatile.LoadUint8(&d.spiReady) != 0 {
		volatile.StoreUint8(&d.spiReady, 0)
		d.sendNext()
	} else {
		volatile.StoreUint8(&d.timerReady, 1)
	}
}

// handleSPIEvent is called from the SPI interrupt, once the DMA buffer has been
// successfully sent.
func (d *Device) handleSPIEvent() {
	// Start the next cycle if the timer has also finished.
	if volatile.LoadUint8(&d.timerReady) != 0 {
		volatile.StoreUint8(&d.timerReady, 0)
		d.sendNext()
	} else {
		volatile.StoreUint8(&d.spiReady, 1)
	}
}
