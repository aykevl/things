//go:build softdevice || !baremetal

package main

// Bluetooth support for the watch.
// This should be kept reasonably portable, so that at least testing on Linux
// will continue to work.

import (
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var batteryLevel bluetooth.Characteristic

func InitBluetooth() error {
	err := adapter.Enable()
	if err != nil {
		return err
	}

	// TODO: use a shorter advertisement interval after start and after losing
	// connection. For example, a 20ms interval for 30 seconds as stated in the
	// Apple guidelines.
	// An interval of 1285 uses around 11µA according to the online power profiler:
	// https://devzone.nordicsemi.com/power/w/opp/2/online-power-profiler-for-bluetooth-le
	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "GopherWatch",
		Interval:  bluetooth.NewDuration(1285 * time.Millisecond),
	})
	if err != nil {
		return err
	}
	err = adv.Start()
	if err != nil {
		return err
	}

	// Add battery service.
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDBattery,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &batteryLevel,
				UUID:   bluetooth.CharacteristicUUIDBatteryLevel,
				Value:  []byte{0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})

	return nil
}

var updateBatteryLevelBuf [1]byte
var batteryLevelValue uint8

func updateBatteryLevel(level uint8) {
	if level == batteryLevelValue {
		return
	}
	updateBatteryLevelBuf[0] = level
	_, err := batteryLevel.Write(updateBatteryLevelBuf[:])
	if err != nil {
		return
	}
	batteryLevelValue = level
}
