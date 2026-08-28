package main

// Not finished. Start of some work on a Home Assistant integration using
// ESPHome.

import (
	"hash/fnv"
	"log/slog"
	"os"
	"time"

	"github.com/hsmade/esphome-go/pkg/server"
	"github.com/hsmade/esphome-go/pkg/server/conf"
)

func main() {
	// set up logging (optional)
	var programLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	// See:
	// https://www.home-assistant.io/integrations/sensor/

	sensorCO2 := conf.Sensor{
		Definition: conf.GenericSensorDefinition{
			BaseSensorDefinition: conf.BaseSensorDefinition{
				ObjectId: "co2",
				Key:      fnvKey("co2"),
				Name:     "CO2",
				UniqueId: "co2",
			},
			DisabledByDefault: false,
			UnitOfMeasurement: "ppm",
			AccuracyDecimals:  0,
			DeviceClass:       "carbon_dioxide", // https://www.home-assistant.io/integrations/sensor#device-class
		},
	}

	sensorTemperature := conf.Sensor{
		Definition: conf.GenericSensorDefinition{
			BaseSensorDefinition: conf.BaseSensorDefinition{
				ObjectId: "temperature",
				Key:      fnvKey("temperature"),
				Name:     "temperature",
				UniqueId: "temperature",
			},
			DisabledByDefault: false,
			UnitOfMeasurement: "°C",
			AccuracyDecimals:  1,
			DeviceClass:       "temperature",
		},
	}

	sensorHumidity := conf.Sensor{
		Definition: conf.GenericSensorDefinition{
			BaseSensorDefinition: conf.BaseSensorDefinition{
				ObjectId: "humidity",
				Key:      fnvKey("humidity"),
				Name:     "humidity",
				UniqueId: "humidity",
			},
			DisabledByDefault: false,
			UnitOfMeasurement: "%",
			AccuracyDecimals:  0,
			DeviceClass:       "humidity",
		},
	}

	sensorUpdates := make(chan conf.SensorUpdate, 1)
	S := server.Server{
		Addr: "localhost:6053",
		Config: conf.Config{
			DeviceInfo: conf.DeviceInfo{
				Name:         "Cloud",
				UsesPassword: false,
			},
			Updates: sensorUpdates,
			Sensors: []conf.Sensor{
				sensorCO2,
				sensorTemperature,
				sensorHumidity,
			},
		},
	}

	// and finally start the server
	go func() {
		err := S.Listen()
		if err != nil {
			slog.Error("main: server exited: %w", err)
		}
	}()

	for {
		// TODO: connect to the local serial port
		time.Sleep(time.Second)
	}
}

func fnvKey(name string) uint32 {
	key := fnv.New32()
	_, _ = key.Write([]byte(name))
	return key.Sum32()
}
