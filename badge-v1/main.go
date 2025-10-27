package main

import (
	"device/stm32"
	"machine"
)

const (
	A1  = machine.PB11
	A2  = machine.PB10
	A3  = machine.PB9
	A4  = machine.PB8
	A5  = machine.PB7 // also PA3
	A6  = machine.PB6
	A7  = machine.PB5
	A8  = machine.PB4 // also PA6
	A9  = machine.PB3 // also PA7
	A10 = machine.PB2
	A11 = machine.PB1
	A12 = machine.PB0

	LED_VCC_ON  = machine.PB12
	PHOTO_SENSE = machine.PB14
	PHOTO_VCC   = machine.PB15
)

func main() {
	configureLEDs()

	setClockSpeed()

	mode := modeNoise

	for i := 0; i < 12; i++ {
		setLEDs(i, 0, 0, 0)
	}

	frame := 0
	index := 0
	for {
		index++
		if index == 12 {
			index = 0
			frame++
		}
		switch index {
		case 0:
			led0 := animate(mode, frame, 2, 1)
			led1 := animate(mode, frame, 3, 0)
			led2 := animate(mode, frame, 6, 3)
			setLEDs(0, led0, led1, led2)
		case 1:
			led0 := animate(mode, frame, 1, 1)
			led1 := animate(mode, frame, 4, 0)
			led2 := animate(mode, frame, 7, 3)
			setLEDs(1, led0, led1, led2)
		case 2:
			led0 := animate(mode, frame, 0, 1)
			led1 := animate(mode, frame, 5, 0)
			led2 := animate(mode, frame, 8, 3)
			setLEDs(2, led0, led1, led2)
		case 9:
			led0 := animate(mode, frame, 6, 2)
			led1 := animate(mode, frame, 6, 1)
			led2 := animate(mode, frame, 6, 0)
			setLEDs(9, led0, led1, led2)
		case 10:
			led0 := animate(mode, frame, 7, 2)
			led1 := animate(mode, frame, 7, 1)
			led2 := animate(mode, frame, 7, 0)
			setLEDs(10, led0, led1, led2)
		case 11:
			led0 := animate(mode, frame, 8, 2)
			led1 := animate(mode, frame, 8, 1)
			led2 := animate(mode, frame, 8, 0)
			setLEDs(11, led0, led1, led2)
		}
		updateLEDs()
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
