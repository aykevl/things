package main

import (
	"encoding/hex"
	"machine"
	"time"

	"codeberg.org/maaike328p/bthome"

	"tinygo.org/x/bluetooth"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/scd4x"
)

var adapter = bluetooth.DefaultAdapter

const bindkey = "d7fe7a79883ec7ef0e43c39ce7ade70f" // yes this is public, I know, don't worry

func main() {
	// Configure CO2 sensor.
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		SCL:       machine.P0_31,
		SDA:       machine.P0_29,
		Frequency: 400_000,
	})
	checkError(err)
	sensor := scd4x.New(i2c)
	err = sensor.Configure()
	checkError(err)
	err = sensor.StartLowPowerPeriodicMeasurement()
	checkError(err)

	// Configure Bluetooth.
	err = adapter.Enable()
	checkError(err)
	adv := adapter.DefaultAdvertisement()
	opts := bluetooth.AdvertisementOptions{
		// Use the longest duration recommended by Apple.
		// This is still definitely long enough to be detected by Home Assistant
		// (the value only updates every 30s anyway).
		// https://developer.apple.com/accessories/Accessory-Design-Guidelines.pdf
		Interval:  bluetooth.NewDuration(1285 * time.Millisecond),
		LocalName: "CO2", // max 5 characters
		ServiceData: []bluetooth.ServiceDataElement{
			{
				UUID: bluetooth.New16BitUUID(bthome.ServiceUUID),
				Data: nil, // filled in later
			},
		},
	}

	// Configure BTHome.
	address, err := adapter.Address()
	checkError(err)
	var key [16]uint8
	hex.Decode(key[:], []byte(bindkey))
	payload := bthome.NewPayload(bthome.Config{
		Bindkey: key,
		MAC:     address.MAC.Address(),
	})
	co2 := payload.AddCO2()
	temp := payload.AddTemperature2()
	hum := payload.AddHumidity0()

	for {
		ready, err := sensor.DataReady()
		checkError(err)
		if ready {
			err := sensor.Update(drivers.Concentration | drivers.Temperature | drivers.Humidity)
			checkError(err)

			println("co2:        ", sensor.CO2())
			println("temperature:", sensor.Temperature())
			println("humidity:   ", sensor.Humidity())

			// Update Bluetooth data.
			adv.Stop()
			co2.Set(uint16(sensor.CO2()))
			temp.Set(int16(sensor.Temperature() / 10))
			hum.Set(uint8(sensor.Humidity() / 100))
			opts.ServiceData[0].Data = payload.EncryptedData()
			err = adv.Configure(opts)
			checkError(err)
			err = adv.Start()
			checkError(err)

			// Data is ready every ~30 seconds. Try again very close to the next
			// expected readout.
			time.Sleep(time.Second * 29)
		} else {
			// Not yet ready, try in one second.
			time.Sleep(time.Second)
		}
	}
}

func checkError(err error) {
	if err != nil {
		for {
			println("got error:", err.Error())
			time.Sleep(time.Second)
		}
	}
}
