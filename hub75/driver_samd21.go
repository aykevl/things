package hub75

import (
	"device/arm"
	"device/sam"
	"machine"
	"time"
	"unsafe"
)

type chipSpecificSettings struct {
	bus *machine.SPI
}

const dmaDescriptors = 2

const dmacChannel = 0

type dmaDescriptor struct {
	btctrl   uint16
	btcnt    uint32
	srcaddr  unsafe.Pointer
	dstaddr  unsafe.Pointer
	descaddr unsafe.Pointer
}

var (
	dmaDescriptorSection          [dmaDescriptors]dmaDescriptor
	dmaDescriptorWritebackSection [dmaDescriptors]dmaDescriptor
)

func (d *Device) configureChip(dataPin, clockPin machine.Pin) {
	for i := 0; i < 5; i++ {
		println("sleep", i)
		time.Sleep(1 * time.Second)
	}

	// enable generic clock
	//sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_ID_SERCOM0_CORE)     // configure generic clock for SERCOM0
	//sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_GEN_GCLK0)          // source is generic clock generator 0
	//sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_CLKEN)                   // enable generic clock

	d.bus = &machine.SPI0
	const triggerSource = 0 // 0x02 // SERCOM0_DMAC_ID_TX
	d.bus.Configure(machine.SPIConfig{
		Frequency: 24000000,
		Mode:      0,
	})

	// Init DMAC
	// http://www.lucadavidian.com/2018/03/08/wifi-controlled-neo-pixels-strips/
	// https://svn.larosterna.com/oss/trunk/arduino/zerotimer/zerodma.cpp
	sam.PM.AHBMASK.SetBits(sam.PM_AHBMASK_DMAC_)
	sam.PM.APBBMASK.SetBits(sam.PM_APBBMASK_DMAC_)

	// Enable DMA IRQ.
	arm.EnableIRQ(sam.IRQ_DMAC)

	sam.DMAC.BASEADDR.Set(uint32(uintptr(unsafe.Pointer(&dmaDescriptorSection))))
	sam.DMAC.WRBADDR.Set(uint32(uintptr(unsafe.Pointer(&dmaDescriptorWritebackSection))))

	println("baseaddr:", sam.DMAC.BASEADDR.Get())

	// Enable peripheral with all priorities.
	sam.DMAC.CTRL.SetBits(sam.DMAC_CTRL_DMAENABLE | sam.DMAC_CTRL_LVLEN0 | sam.DMAC_CTRL_LVLEN1 | sam.DMAC_CTRL_LVLEN2 | sam.DMAC_CTRL_LVLEN3)

	// Add channel.
	// First disable the DMAC peripheral.
	//sam.DMAC.CTRL.ClearBits(sam.DMAC_CTRL_DMAENABLE)
	//for sam.DMAC.CTRL.HasBits(sam.DMAC_CTRL_DMAENABLE) {
	//}
	sam.DMAC.CHID.Set(0) // channel 0
	sam.DMAC.CHCTRLA.ClearBits(sam.DMAC_CHCTRLA_ENABLE)
	sam.DMAC.CHCTRLA.SetBits(sam.DMAC_CHCTRLA_SWRST)
	sam.DMAC.SWTRIGCTRL.ClearBits(1 << dmacChannel)

	sam.DMAC.CHCTRLB.Set((sam.DMAC_CHCTRLB_LVL_LVL0 << sam.DMAC_CHCTRLB_LVL_Pos) | (sam.DMAC_CHCTRLB_TRIGACT_BEAT << sam.DMAC_CHCTRLB_TRIGACT_Pos) | (triggerSource << sam.DMAC_CHCTRLB_TRIGSRC_Pos))
	println("CHCTRLB:", sam.DMAC.CHCTRLB.Get())

	// Enable all interrupts.
	sam.DMAC.CHINTENSET.SetBits(sam.DMAC_CHINTENSET_TCMPL | sam.DMAC_CHINTENSET_SUSP)

	dmaDescriptorSection[0] = dmaDescriptor{
		btctrl: (1 << 0) | // VALID: Descriptor Valid
			(0 << 3) | // BLOCKACT=NOACT: Block Action
			(1 << 10) | // SRCINC: Source Address Increment Enable
			(0 << 11) | // DSTINC: Destination Address Increment Enable
			(0 << 12) | // STEPSEL=SRC: Step Selection
			(0 << 13), // STEPSIZE=X1: Address Increment Step Size
		btcnt:    24, // beat count
		dstaddr:  unsafe.Pointer(&d.displayBitstrings[0][3][0]),//unsafe.Pointer(&d.bus.Bus.DATA.Reg),
		srcaddr: unsafe.Pointer(&d.displayBitstrings[0][1][0]),
	}

	// Enable channel.
	sam.DMAC.CHID.Set(0) // channel 0
	sam.DMAC.CHCTRLA.SetBits(sam.DMAC_CHCTRLA_ENABLE)

	// Trigger!
	sam.DMAC.SWTRIGCTRL.Set(1 << dmacChannel)

	// Enable DMA.
	//sam.DMAC.CTRL.SetBits(sam.DMAC_CTRL_DMAENABLE)
}

// startTransfer starts the SPI transaction to send the next row to the screen.
func (d *Device) startTransfer() {
	bitstring := d.displayBitstrings[d.row][d.colorBit]
	//d.bus.Tx(bitstring, nil)

	//println(uintptr(dmaDescriptorSection[0].btctrl))
	//dmaDescriptorSection[0].srcaddr = unsafe.Pointer(&bitstring[0])

	//sam.DMAC.CHCTRLA.SetBits(sam.DMAC_CHCTRLA_ENABLE)
	//sam.DMAC.SWTRIGCTRL.Set(1 << dmacChannel)
	d.bus.Bus.DATA.Set(uint32(bitstring[0]))
	//for _, w := range bitstring[1:] {
		// write data
		//d.bus.Bus.DATA.Set(uint32(w))

		// wait for receive
		//for !d.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_RXC) {
		//}
		//d.bus.Bus.DATA.Get()
	//}

	//// wait for receive
	//for !d.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_RXC) {
	//}
	//d.bus.Bus.DATA.Get()
}

// startOutputEnableTimer will enable and disable the screen for a very short
// time, depending on which bit is currently enabled.
func (d *Device) startOutputEnableTimer() {
	count := d.brightness << d.colorBit
	for i := uint32(0); i < count; i++ {
		d.oe.Low()
	}
	d.oe.High()
}

//go:export DMAC_IRQHandler
func dmacHandler() {
	println("DMA handler")
}
