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

var connectedDevice chan bluetooth.Device

func InitBluetooth() error {
	err := adapter.Enable()
	if err != nil {
		return err
	}

	adapter.SetConnectHandler(handleBLEConnection)
	connectedDevice = make(chan bluetooth.Device, 1)
	go connectionHandler()

	// TODO: use a shorter advertisement interval after start and after losing
	// connection. For example, a 20ms interval for 30 seconds as stated in the
	// Apple guidelines:
	// https://developer.apple.com/accessories/Accessory-Design-Guidelines.pdf
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

func handleBLEConnection(device bluetooth.Device, connected bool) {
	if connected {
		select {
		case connectedDevice <- device:
		default:
		}
	}
}

// Background goroutine that updates connection parameters as needed.
func connectionHandler() {
	for device := range connectedDevice {
		// Wait a bit after connecting so that initial negotiating can be
		// faster.
		time.Sleep(time.Second * 5)

		// Following the Apple accessory design guidelines, picking a connection
		// latency of around 500ms that is a multiple of 15ms (and giving the
		// device 15ms of space). My Android 13 phone picks 510ms as the
		// connection interval with these parameters.
		// For comparison, the Mi Band 3 negotiates 517.5ms as the connection
		// interval after a sync.
		device.RequestConnectionParams(bluetooth.ConnectionParams{
			MinInterval: bluetooth.NewDuration(495 * time.Millisecond),
			MaxInterval: bluetooth.NewDuration(510 * time.Millisecond),
			Timeout:     bluetooth.NewDuration(5 * time.Second),
		})
	}
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
