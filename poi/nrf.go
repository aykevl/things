// +build nrf

package main

import (
	"github.com/tinygo-org/bluetooth"
)

var (
	serviceUUID    = bluetooth.NewUUID([16]byte{0xd8, 0x65, 0x00, 0x01, 0x90, 0xc8, 0x46, 0x7a, 0xb3, 0xd2, 0xe4, 0xc1, 0x8a, 0x61, 0x10, 0x7f})
	animationUUID  = bluetooth.NewUUID([16]byte{0xd8, 0x65, 0x00, 0x02, 0x90, 0xc8, 0x46, 0x7a, 0xb3, 0xd2, 0xe4, 0xc1, 0x8a, 0x61, 0x10, 0x7f})
	speedUUID      = bluetooth.NewUUID([16]byte{0xd8, 0x65, 0x00, 0x03, 0x90, 0xc8, 0x46, 0x7a, 0xb3, 0xd2, 0xe4, 0xc1, 0x8a, 0x61, 0x10, 0x7f})
	brightnessUUID = bluetooth.NewUUID([16]byte{0xd8, 0x65, 0x00, 0x04, 0x90, 0xc8, 0x46, 0x7a, 0xb3, 0xd2, 0xe4, 0xc1, 0x8a, 0x61, 0x10, 0x7f})
	colorUUID      = bluetooth.NewUUID([16]byte{0xd8, 0x65, 0x00, 0x05, 0x90, 0xc8, 0x46, 0x7a, 0xb3, 0xd2, 0xe4, 0xc1, 0x8a, 0x61, 0x10, 0x7f})
)

func initHardware() {
	println("init hardware")
	// Initialize the adapter.
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		println("could not enable:", err.Error())
		return
	}

	// Start advertising.
	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    bluetoothName,
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	})
	if err != nil {
		println("could not configure advertisement:", err.Error())
		return
	}
	err = adv.Start()
	if err != nil {
		println("could not start advertisement:", err.Error())
		return
	}

	err = adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  animationUUID,
				Value: []byte{animationIndex},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) < 1 {
						return
					}
					if int(value[0]) < len(animations) {
						animationIndex = value[0]
						println("animation is now:", animationIndex)
					}
				},
			},
			{
				UUID:  speedUUID,
				Value: []byte{speed},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) < 1 {
						return
					}
					speed = value[0]
					println("speed is now:", speed)
				},
			},
			{
				UUID:  brightnessUUID,
				Value: []byte{baseColor.A},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) < 1 {
						return
					}
					baseColor.A = value[0]
					println("brightness is now:", baseColor.A)
				},
			},
			{
				UUID:  colorUUID,
				Value: []byte{baseColor.R, baseColor.G, baseColor.B},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) < 3 {
						return
					}
					baseColor.R = value[0]
					baseColor.G = value[1]
					baseColor.B = value[2]
					println("updated color")
				},
			},
		},
	})
	if err != nil {
		println("could not add poi service:", err.Error())
		return
	}
}
