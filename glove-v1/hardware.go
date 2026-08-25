package main

import (
	"device/stm32"
	"machine"
)

const (
	LEDRed   = machine.PA3 // TIM2 channel 4
	LEDGreen = machine.PA2 // TIM2 channel 3
	LEDBlue  = machine.PA1 // TIM2 channel 2

	Button = machine.PA0
)

func main() {
	setClockSpeed()

	LEDRed.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LEDGreen.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LEDBlue.Configure(machine.PinConfig{Mode: machine.PinOutput})

	Button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	machine.TIM2.Configure(machine.PWMConfig{Period: 8000 * 4})
	machine.TIM2.Channel(LEDRed)
	machine.TIM2.Channel(LEDGreen)
	machine.TIM2.Channel(LEDBlue)
	println("top:", machine.TIM2.Top())

	// Main loop
	mode := firstMode
	buttonWasPressed := false
	for cycles := uint32(0); ; cycles++ {
		animation(cycles, mode)

		// Handle button presses.
		buttonPressed := !Button.Get()
		if buttonPressed && !buttonWasPressed {
			mode++
			if mode == modeLast {
				mode = 0
			}
		}
		buttonWasPressed = buttonPressed
	}
}

func setColor(c Color) {
	// Top is assumed to be a fixed number.
	const top = 1024
	machine.TIM2.Set(3, top-uint32(c.R())) // red
	machine.TIM2.Set(2, top-uint32(c.G())) // green
	machine.TIM2.Set(1, top-uint32(c.B())) // blue
}

// Clock speed configuration. This has been copied from earring-ring-rgb36v2 mostly.
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
	stm32.RCC.SetICSCR_MSIRANGE(stm32.RCC_ICSCR_MSIRANGE_Range3)

	// Reduce PCLK2/PCLK1 clocks since we don't need those peripherals (GPIO is
	// directly connected to the CPU). This saves around ~1.2µA.
	stm32.RCC.SetCFGR_PPRE1(stm32.RCC_CFGR_PPRE1_Div16)
	stm32.RCC.SetCFGR_PPRE2(stm32.RCC_CFGR_PPRE2_Div16)

	// Divide SYSCLK, for testing.
	//stm32.RCC.SetCFGR_HPRE(stm32.RCC_CFGR_HPRE_Div512)
}
