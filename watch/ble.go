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
		LocalName: "InfiniTime",
		Interval:  bluetooth.NewDuration(1285 * time.Millisecond),
	})
	if err != nil {
		return err
	}
	err = adv.Start()
	if err != nil {
		return err
	}

	// Add Device Information Service. This is necessary for Gadgetbridge,
	// otherwise it keeps showing an error ("the bind value at index 2 is
	// null").
	err = adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDDeviceInformation,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDManufacturerNameString,
				Value: []byte("Pine64"),
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFirmwareRevisionString,
				Value: []byte("GopherWatch-dev"), // unspecified version
				Flags: bluetooth.CharacteristicReadPermission,
			},
		},
	})
	if err != nil {
		return err
	}

	// Add battery service.
	err = adapter.AddService(&bluetooth.Service{
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
	if err != nil {
		return err
	}

	// Current Time Service. This enables Gadgetbridge to sync the time.
	err = adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDCurrentTime,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDCurrentTime,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 10 {
						return // unexpected value
					}
					year := int(value[0]) | int(value[1])<<8
					month := time.Month(value[2])
					day := int(value[3])
					hour := int(value[4])
					minute := int(value[5])
					second := int(value[6])
					nanosecond := int(value[7]) * (1e9 / 256)
					newTime := time.Date(year, month, day, hour, minute, second, nanosecond, time.UTC)
					oldTime := time.Now()
					diff := newTime.Sub(oldTime)
					adjustTime(diff)
				},
			},
		},
	})
	if err != nil {
		return err
	}

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
