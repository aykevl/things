package main

// Bluetooth support for the watch.
// This should be kept reasonably portable, so that at least testing on Linux
// will continue to work.

import (
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func InitBluetooth() error {
	err := adapter.Enable()
	if err != nil {
		return err
	}

	// TODO: use a shorter advertisement interval after start and after losing
	// connection. For example, a 20ms interval for 30 seconds as stated in the
	// Apple guidelines.
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

	return nil
}
