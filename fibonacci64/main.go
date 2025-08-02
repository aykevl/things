package main

import (
	"device/arm"
	"device/stm32"
	"machine"
	"runtime/volatile"
)

const (
	A1  = machine.PA11
	A2  = machine.PA10
	A3  = machine.PA9
	A4  = machine.PA8
	A5  = machine.PA7
	A6  = machine.PA6
	A7  = machine.PA5
	A8  = machine.PA4
	A9  = machine.PA3
	A10 = machine.PA2
	A11 = machine.PA1
	A12 = machine.PA0
	A13 = machine.PB6
	A14 = machine.PB5
	A15 = machine.PB4
	A16 = machine.PB3
	A17 = machine.PB2
	A18 = machine.PB1
	A19 = machine.PB0

	button = machine.PB7
)

// Set to one from the interrupt to indicate the button was pressed.
var buttonWake volatile.Register8

func main() {
	// Configure button
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	button.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		buttonWake.Set(1)
	})

	// Configure pins
	configureLEDs()

	// Set A1-A12 as open drain (and importantly, skip SWDIO/SWCLK)
	stm32.GPIOA.OTYPER.Set(0b0000_1111_1111_1111)

	// Set A13-A19 as open drain.
	stm32.GPIOB.OTYPER.Set(0b0000_0000_0111_1111)

	setClockSpeed()

	// Zero all LEDs.
	for i := 0; i < 19; i++ {
		animateLEDs(modeOff, i, 0)
	}

	index := 0 // 0..18, group of 3 LEDs that will be updated together
	frame := 0
	mode := initialMode
	buttonPressed := false
	for {
		animateLEDs(mode, index, frame)

		// Bitbang the LEDs.
		updateLEDs()

		index++
		if index == 19 {
			index = 0
			frame++

			// Read the button every frame update.
			pressed := !button.Get() // low means pressed
			if pressed && pressed != buttonPressed {
				mode++
				if mode >= modeLast {
					// Last, so wrap around.
					mode = 0
				}
				if mode == modeOff {
					// Sleep until a button press, then go to the next mode.
					sleepUntilButtonPress()
					mode++
				}

				// Clear LEDs to be sure.
				for i := 0; i < 19; i++ {
					animateLEDs(modeOff, i, 0)
				}
			}
			buttonPressed = pressed
		}
	}
}

func setClockSpeed() {
	// Switch to MSI.
	stm32.RCC.CFGR.ReplaceBits(stm32.RCC_CFGR_SWS_MSI, stm32.RCC_CFGR_SW_Msk, 0)
	for stm32.RCC.CFGR.Get()&stm32.RCC_CFGR_SW_Msk != stm32.RCC_CFGR_SWS_MSI {
	}

	// Disable PLL.
	stm32.RCC.CR.ClearBits(stm32.RCC_CR_PLLON)

	// Configure power:
	// - change to 1.2V range.
	// - enable LPDS or LPSDSR bit
	// - disable Vrefint (ultra low power)
	// - enter standby mode when entering deepsleep (PDDS=1)
	// - clear wakeup flag
	stm32.RCC.APB1ENR.SetBits(stm32.RCC_APB1ENR_PWREN)
	stm32.PWR.CR.Set(stm32.PWR_CR_VOS_V1_2<<stm32.PWR_CR_VOS_Pos |
		stm32.PWR_CR_LPSDSR |
		stm32.PWR_CR_ULP |
		//stm32.PWR_CR_PDDS | // set this bit to enter standby mode (instead of stop mode)
		stm32.PWR_CR_CWUF |
		0)
	stm32.PWR.CR.Get() // make sure the values have been written
	stm32.RCC.APB1ENR.ClearBits(stm32.RCC_APB1ENR_PWREN)

	// Disable HSI16 clock.
	stm32.RCC.CR.ClearBits(stm32.RCC_CR_HSI16ON)

	// Change flash latency to zero wait states. Saves 0.4µA or so.
	stm32.FLASH.SetACR_LATENCY(stm32.Flash_ACR_LATENCY_WS0)

	// Disable flash during sleep.
	stm32.FLASH.ACR.Set(stm32.Flash_ACR_SLEEP_PD)

	// Disable TIM21. Doesn't use much current when the clock is disabled, but
	// it's a (very) small win. Saves 0.2µA or so.
	stm32.RCC.SetAPB2ENR_TIM21EN(0)

	// Set MSI clock speed.
	// Range 2: around 262kHz (~55µA with all LEDs off)
	// Range 3: around 524kHz (~83µA with all LEDs off)
	stm32.RCC.SetICSCR_MSIRANGE(stm32.RCC_ICSCR_MSIRANGE_Range4)

	// Reduce PCLK2/PCLK1 clocks since we don't need those peripherals (GPIO is
	// directly connected to the CPU). This saves around ~1.2µA.
	stm32.RCC.SetCFGR_PPRE1(stm32.RCC_CFGR_PPRE1_Div16)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Divide SYSCLK, for testing.
	//stm32.RCC.SetCFGR_HPRE(stm32.RCC_CFGR_HPRE_Div512)
}

func sleepUntilButtonPress() {
	// Goal: stop mode, without RTC, with 1 GPIO pin enabled for button interrupt

	// Disable GPIO pins during sleep.
	disableLEDs()

	// Enter stop mode, wake up on a button press.
	buttonWake.Set(0)
	arm.SCB.SCR.SetBits(arm.SCB_SCR_SLEEPDEEP)
	for {
		arm.Asm("wfe")
		if buttonWake.Get() != 0 {
			break
		}
	}
	arm.SCB.SCR.ClearBits(arm.SCB_SCR_SLEEPDEEP)

	// Restore GPIO pins to their previous configuration.
	configureLEDs()
}
