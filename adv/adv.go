package main

import (
	"time"

	"github.com/aykevl/go-ble/s132v6"
)

func main() {
	println("starting...")
	time.Sleep(1 * time.Second)

	println("sd enabled:", sd.IsEnabled())
	println("sd enable:", sd.Enable(sd.DefaultClockSource))
	println("sd enabled:", sd.IsEnabled())

	var app_ram_base uintptr = 0x200039c0
	ram, err := sd.EnableBLE(app_ram_base)
	println("ble enabled:", ram, err)

	adv := sd.NewAdvertisement()

	params := &sd.AdvParams{
		Properties: sd.AdvProperties{
			Type:   sd.AdvTypeConnectableScannableUndirected,
			Fields: 0,
		},
		Interval:     100,
		Duration:     0, // unlimited
		MaxAdvEvts:   0,
		FilterPolicy: 0,
		PrimaryPhi:   0,
		SecondaryPhi: 0,
		Fields:       0,
	}
	adv.Configure("\x02\x01\x06"+"\x07\x09TinyGo", "", params) // flags + local name
	println("ble adv configure:", err)

	println("ble adv start:", adv.Start())

	for {
		time.Sleep(10 * time.Second)
	}
}
