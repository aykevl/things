//go:build attiny1616

package main

import (
	"device/avr"
	"machine"
	"runtime/interrupt"
	"unsafe"
)

// #include "bitbang.h"
import "C"

func init() {
	// Use 20MHz/32 = 625kHz
	// This results in a current consumption of around 0.41mA with all LEDs off.
	// It *should* result in a current consumption of 187.5µA, but I think the
	// 20MHz oscillator uses a fair bit of current too.
	avr.CPU.CCP.Set(0xD8)                 // unlock protected registers
	avr.CLKCTRL.MCLKCTRLB.Set(0x4<<1 | 1) // prescaler of 32
}

const button = machine.PC3

func initHardware() {
	avr.PORTA.DIR.Set(0b0111_1110) // Configure PA1-PA6 as output.
	avr.PORTB.DIR.Set(0b0011_1111) // Configure PB0-PB5 as output (R2, G2, B2, R3, G3, B3)
	avr.PORTB.OUT.Set(0b0011_1111) // Set PB0-PB5 low
	avr.PORTC.DIR.Set(0b0000_0111) // Configure PC0-PC2 as output (R1, G1, B1)
	avr.PORTC.OUT.Set(0b0000_0111) // Set PC0-PC2 low

	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// This is the only remaining pin. Configure it as output just in case that
	// helps with standby current consumption.
	machine.PA7.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func disableLEDs() {
	avr.PORTA.OUTCLR.Set(0b0111_1110) // Set PA1-PA6 high
	avr.PORTB.OUTCLR.Set(0b0011_1111) // Set PB0-PB5 high
	avr.PORTC.OUTCLR.Set(0b0000_0111) // Set PC0-PC2 high
}

func isButtonPressed() bool {
	return !button.Get()
}

func updateLEDs() {
	// Pinout for the anodes:
	//   A1: PA2
	//   A2: PA1
	//   A3: PA6
	//   A4: PA5
	//   A5: PA4
	//   A6: PA3
	// Pinout for the cathodes:
	//   R1: PC0
	//   G1: PC1
	//   B1: PC2
	//   R2: PB0
	//   G2: PB1
	//   B2: PB2
	//   R3: PB3
	//   G3: PB4
	//   B3: PB5

	// order of the LEDs:
	// index   part #   Rx/Gx/Bx  anode
	// 0       LED16    1         A6
	// 1       LED17    2         A6
	// 2       LED18    3         A6
	// 3       LED13    1         A5
	// 4       LED14    2         A5
	// 5       LED15    3         A5
	// 6       LED11    2         A4
	// 7       LED12    3         A4
	// 8       LED10    1         A4
	// 9       LED8     2         A3
	// 10      LED7     1         A3
	// 11      LED9     3         A3
	// 12      LED6     3         A2
	// 13      LED5     2         A2
	// 14      LED4     1         A2
	// 15      LED3     3         A1
	// 16      LED1     1         A1
	// 17      LED2     2         A1

	state := interrupt.Disable()

	// R1
	avr.PORTC.OUTTGL.Set(1 << 0)
	showLEDs(
		leds[10].R, // PA6, A3
		leds[ 8].R, // PA5, A4
		leds[ 3].R, // PA4, A5
		leds[ 0].R, // PA3, A6
		leds[16].R, // PA2, A1
		leds[14].R, // PA1, A2
	)
	avr.PORTC.OUTTGL.Set(1 << 0)

	// R2
	avr.PORTB.OUTTGL.Set(1 << 0)
	showLEDs(
		leds[ 9].R, // PA6, A3
		leds[ 6].R, // PA5, A4
		leds[ 4].R, // PA4, A5
		leds[ 1].R, // PA3, A6
		leds[17].R, // PA2, A1
		leds[13].R, // PA1, A2
	)
	avr.PORTB.OUTTGL.Set(1 << 0)

	// R3
	avr.PORTB.OUTTGL.Set(1 << 3)
	showLEDs(
		leds[11].R, // PA6, A3
		leds[ 7].R, // PA5, A4
		leds[ 5].R, // PA4, A5
		leds[ 2].R, // PA3, A6
		leds[15].R, // PA2, A1
		leds[12].R, // PA1, A2
	)
	avr.PORTB.OUTTGL.Set(1 << 3)

	// G1
	avr.PORTC.OUTTGL.Set(1 << 1)
	showLEDs(
		leds[10].G, // PA6, A3
		leds[ 8].G, // PA4, A4
		leds[ 3].G, // PA4, A5
		leds[ 0].G, // PA3, A6
		leds[16].G, // PA2, A1
		leds[14].G, // PA1, A2
	)
	avr.PORTC.OUTTGL.Set(1 << 1)

	// G2
	avr.PORTB.OUTTGL.Set(1 << 1)
	showLEDs(
		leds[ 9].G, // PA6, A3
		leds[ 6].G, // PA5, A4
		leds[ 4].G, // PA4, A5
		leds[ 1].G, // PA3, A6
		leds[17].G, // PA2, A1
		leds[13].G, // PA1, A2
	)
	avr.PORTB.OUTTGL.Set(1 << 1)

	// G3
	avr.PORTB.OUTTGL.Set(1 << 4)
	showLEDs(
		leds[11].G, // PA6, A3
		leds[ 7].G, // PA5, A4
		leds[ 5].G, // PA4, A5
		leds[ 2].G, // PA3, A6
		leds[15].G, // PA2, A1
		leds[12].G, // PA1, A2
	)
	avr.PORTB.OUTTGL.Set(1 << 4)

	// B1
	avr.PORTC.OUTTGL.Set(1 << 2)
	showLEDs(
		leds[10].B, // PA6, A3
		leds[ 8].B, // PA4, A4
		leds[ 3].B, // PA4, A5
		leds[ 0].B, // PA3, A6
		leds[16].B, // PA2, A1
		leds[14].B, // PA1, A2
	)
	avr.PORTC.OUTTGL.Set(1 << 2)

	// B2
	avr.PORTB.OUTTGL.Set(1 << 2)
	showLEDs(
		leds[ 9].B, // PA6, A3
		leds[ 6].B, // PA5, A4
		leds[ 4].B, // PA4, A5
		leds[ 1].B, // PA3, A6
		leds[17].B, // PA2, A1
		leds[13].B, // PA1, A2
	)
	avr.PORTB.OUTTGL.Set(1 << 2)

	// B3
	avr.PORTB.OUTTGL.Set(1 << 5)
	showLEDs(
		leds[11].B, // PA6, A3
		leds[ 7].B, // PA5, A4
		leds[ 5].B, // PA4, A5
		leds[ 2].B, // PA3, A6
		leds[15].B, // PA2, A1
		leds[12].B, // PA1, A2
	)
	avr.PORTB.OUTTGL.Set(1 << 5)

	interrupt.Restore(state)
}

// The bitbang function is very large, but without the //go:noinline it would be
// inlined. I guess LLVM doesn't recognize the inline assembly is very large.
//
//go:noinline
func showLEDs(c1, c2, c3, c4, c5, c6 uint8) {
	port := (*uint8)(unsafe.Pointer(uintptr(0x0001))) // VPORTA.OUT (alias of PORTA.OUT in I/O space)
	C.bitbang_show_leds(c1, c2, c3, c4, c5, c6, port)
}
