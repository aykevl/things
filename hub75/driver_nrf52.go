package hub75

import (
	"device/arm"
	"device/nrf"
	"machine"
	"unsafe"
)

type chipSpecificSettings struct {
	bus           *nrf.SPIM_Type
	timer         *nrf.TIMER_Type
	gpioteChannel uint8 // GPIOTE channel used for the Output Enable pin.
	ppiChannel0   uint8 // PPI channel used to turn the screen on.
	ppiChannel1   uint8 // PPI channel used to turn the screen off.
}

func (d *Device) configureChip(dataPin, clockPin machine.Pin) {
	// Set some properties that are chip-specific.
	d.bus = nrf.SPIM0
	d.timer = nrf.TIMER0
	d.gpioteChannel = 0
	d.ppiChannel0 = 0
	d.ppiChannel1 = 1

	// Configure the SPI bus.
	d.bus.ENABLE.Set(nrf.SPIM_ENABLE_ENABLE_Disabled)
	d.bus.FREQUENCY.Set(nrf.SPIM_FREQUENCY_FREQUENCY_M8)
	d.bus.CONFIG.Set(0) // default config
	d.bus.PSEL.MISO.Set(0xffffffff)
	d.bus.PSEL.MOSI.Set(uint32(dataPin))
	d.bus.PSEL.SCK.Set(uint32(clockPin))
	d.bus.INTENSET.Set(nrf.SPIM_INTENSET_ENDTX_Msk)
	d.bus.ENABLE.Set(nrf.SPIM_ENABLE_ENABLE_Enabled)

	// Configure the SPI interrupt handler to the highest priority.
	// It doesn't need to be the highest, but it should be the same priority as
	// the timer.
	arm.SetPriority(nrf.IRQ_SPIM0_SPIS0_TWIM0_TWIS0_SPI0_TWI0, 0x00)
	arm.EnableIRQ(nrf.IRQ_SPIM0_SPIS0_TWIM0_TWIS0_SPI0_TWI0)
	d.bus.TXD.LIST.Set(nrf.SPIM_TXD_LIST_LIST_ArrayList)
	d.bus.RXD.MAXCNT.Set(0)

	// Configure the timer interrupt handler at the highest priority.
	// Also, make sure that when the timer fires, it automatically stops itself.
	arm.SetPriority(nrf.IRQ_TIMER0, 0x00)
	arm.EnableIRQ(nrf.IRQ_TIMER0)
	d.timer.PRESCALER.Set(0)
	d.timer.SHORTS.SetBits(nrf.TIMER_SHORTS_COMPARE1_CLEAR | nrf.TIMER_SHORTS_COMPARE1_STOP)
	d.timer.CC[0].Set(1)
	d.timer.INTENSET.Set(nrf.TIMER_INTENSET_COMPARE1)

	// Configure a GPIOTE channel.
	nrf.GPIOTE.CONFIG[d.gpioteChannel].Set(
		(nrf.GPIOTE_CONFIG_MODE_Task << nrf.GPIOTE_CONFIG_MODE_Pos) |
			(uint32(d.oe) << nrf.GPIOTE_CONFIG_PSEL_Pos) |
			(nrf.GPIOTE_CONFIG_POLARITY_None << nrf.GPIOTE_CONFIG_POLARITY_Pos) |
			(nrf.GPIOTE_CONFIG_OUTINIT_High << nrf.GPIOTE_CONFIG_OUTINIT_Pos))

	// Set up the PPI channels.
	// We use one channel (ppiChannel0 with CC[0]) to turn the screen on, and
	// another (ppiChannel1 with CC[1]) to turn it off again.
	nrf.PPI.CHENSET.Set(1 << d.ppiChannel0)
	nrf.PPI.CH[d.ppiChannel0].EEP.Set(uint32(uintptr(unsafe.Pointer(&d.timer.EVENTS_COMPARE[0]))))
	nrf.PPI.CH[d.ppiChannel0].TEP.Set(uint32(uintptr(unsafe.Pointer(&nrf.GPIOTE.TASKS_CLR[d.gpioteChannel]))))
	nrf.PPI.CHENSET.Set(1 << d.ppiChannel1)
	nrf.PPI.CH[d.ppiChannel1].EEP.Set(uint32(uintptr(unsafe.Pointer(&d.timer.EVENTS_COMPARE[1]))))
	nrf.PPI.CH[d.ppiChannel1].TEP.Set(uint32(uintptr(unsafe.Pointer(&nrf.GPIOTE.TASKS_SET[d.gpioteChannel]))))
}

// startTransfer starts the SPI transaction to send the next row to the screen.
func (d *Device) startTransfer() {
	bitstring := d.displayBitstrings[d.row][d.colorBit]
	d.bus.TXD.MAXCNT.Set(uint32(len(bitstring)))
	d.bus.TXD.PTR.Set(uint32(uintptr(unsafe.Pointer(&bitstring[0]))))
	d.bus.TASKS_START.Set(1)
}

// startOutputEnableTimer will enable and disable the screen for a very short
// time, depending on which bit is currently enabled.
func (d *Device) startOutputEnableTimer() {
	// Turn the screen on for the specified time.
	// Note that it is actually turned on at T+1 (by a PPI on CC[0]) so we also
	// have to add one to the time to get the correct turn-on time.
	d.timer.CC[1].Set(d.brightness<<d.colorBit + 1)
	d.timer.TASKS_START.Set(1)
}

//export SPIM0_SPIS0_TWIM0_TWIS0_SPI0_TWI0_IRQHandler
func spim0Handler() {
	if nrf.SPIM0.EVENTS_ENDTX.Get() != 0 {
		nrf.SPIM0.EVENTS_ENDTX.Set(0)
		display.handleSPIEvent()
	}
}

//export TIMER0_IRQHandler
func timer0Handler() {
	if nrf.TIMER0.EVENTS_COMPARE[1].Get() != 0 {
		nrf.TIMER0.EVENTS_COMPARE[1].Set(0)
		display.handleTimerEvent()
	}
}
