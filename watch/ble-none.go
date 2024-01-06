//go:build !(softdevice || !baremetal)

// Dummy Bluetooth implementation, to be able to run this watch on systems that
// don't support Bluetooth.

package main

func InitBluetooth() error {
	// nothing to do
	return nil
}

func updateBatteryLevel(level uint8) {
}

func updateStepCountValue(stepCount uint32) {
}
