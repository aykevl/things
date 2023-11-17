package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/scd4x"
	"tinygo.org/x/drivers/sgp30"
	"tinygo.org/x/drivers/ws2812"
)

const NUM_LEDS = 50

var leds = make([]color.RGBA, NUM_LEDS)

const animationSpeed = 2 // higher means faster
var brightness uint8 = 127
var palette = ledsgo.PartyColors

const measurementInterval = time.Second * 5

var (
	co2Sensor *scd4x.Device
	vocSensor *sgp30.Device
)

type ledPosition struct {
	X uint8
	Y uint8
}

var brightnessMap = [10]uint8{2, 7, 17, 32, 52, 77, 108, 146, 189, 238}

var positions = [...]ledPosition{
	{43, 16}, {44, 16}, {44, 19}, {37, 19}, {39, 9}, {44, 19}, {50, 13}, {41, 12}, {39, 9}, // ball 1
	{57, 20}, {66, 21}, {75, 19}, {64, 15}, {63, 24}, // ball 2
	{65, 28}, {75, 30}, {70, 33}, {80, 40}, {78, 36}, {75, 39}, // ball 3
	{70, 41}, {70, 39}, {62, 37}, {52, 43}, {62, 52}, {66, 40}, {61, 36}, {52, 40}, // ball 4
	{43, 39}, {38, 37}, {45, 36}, {42, 34}, {34, 40}, {42, 44}, {46, 39}, {40, 36}, // ball 5
	{29, 36}, {26, 43}, {29, 36}, // ball 6
	{21, 35}, {15, 37}, {20, 39}, // ball 7
	{29, 29}, {33, 20}, {28, 12}, {23, 20}, {20, 19}, {22, 21}, {28, 27}, {33, 25}, // ball 8
}

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	LED_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(LED_PIN)

	time.Sleep(time.Second * 1)
	println("start")
	go runSensors()

	var command byte
	animation := noise
	for {
		if command == 0 {
			// Read command.
			if machine.Serial.Buffered() != 0 {
				command, _ = machine.Serial.ReadByte()
			}
		}
		if command != 0 {
			switch command {
			case 'D': // disable
				animation = poweroff
				command = 0
			case 'W': // white
				animation = white
				command = 0
			case 'P':
				animation = noise
				palette = ledsgo.PartyColors
				command = 0
			case 'F':
				animation = noise
				palette = ledsgo.ForestColors
				command = 0
			case 'O':
				animation = noise
				palette = ledsgo.OceanColors
				command = 0
			case 'L': // lightning
				animation = lightning
				command = 0
			case 'b': // brightness
				if machine.Serial.Buffered() != 0 {
					b, _ := machine.Serial.ReadByte()
					if b >= '0' && b <= '9' {
						brightness = brightnessMap[b-'0']
					}
					command = 0
				}
			}
		}

		// Update colors.
		var t uint64
		if animationSpeed != 0 {
			t = uint64(time.Now().UnixNano() >> (26 - animationSpeed))
		}
		animation(t, leds)

		// Send new colors to LEDs.
		for _, c := range leds {
			strip.WriteByte(c.G) // G
			strip.WriteByte(c.R) // R
			strip.WriteByte(c.B) // B
			strip.WriteByte(c.A) // W (alpha channel, used as white channel)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func noise(t uint64, leds []color.RGBA) {
	for i := range leds {
		pos := positions[i]
		const spread = 24 // higher means more detail
		val := ledsgo.Noise3(uint32(t), uint32(pos.X)*spread, uint32(pos.Y)*spread)
		c := palette.ColorAt(val)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}

func lightning(t uint64, leds []color.RGBA) {
	const interval = 1 << 8
	elapsed := interval - t%interval
	for i := range leds[:10] {
		leds[i] = color.RGBA{0, 0, 0, uint8(elapsed / (interval / 100))}
	}
}

func white(t uint64, leds []color.RGBA) {
	for i := range leds {
		leds[i] = color.RGBA{A: brightness}
	}
}

func poweroff(t uint64, leds []color.RGBA) {
	for i := range leds {
		leds[i] = color.RGBA{}
	}
}

func xorshift64(x uint64) uint64 {
	// https://en.wikipedia.org/wiki/Xorshift
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	return x
}

func runSensors() {
	configureSensors()

	for n := uint64(0); ; n++ {
		// Sample CO₂ every 10 seconds.
		if n%10 == 0 && co2Sensor != nil {
			sampleCO2Sensor(co2Sensor)
		}
		// Sample VOC every second.
		if vocSensor != nil {
			sampleVOCSensor(vocSensor)
		}
		time.Sleep(time.Second)
	}
}

func configureSensors() {
	bus := machine.I2C1
	err := bus.Configure(machine.I2CConfig{
		SDA:       machine.GP26,
		SCL:       machine.GP27,
		Frequency: 400 * machine.KHz,
	})
	if err != nil {
		println("could not configure I2C:", bus)
		return
	}
	{
		// Configure SCD40 CO₂ sensor.
		sensor := scd4x.New(bus)
		if !sensor.Connected() {
			println("CO2 sensor is not connected!")
			return
		}
		if err := sensor.Configure(); err != nil {
			println("could not configure CO2 sensor:", err.Error())
			return
		}
		println("configured!")

		if err := sensor.StartPeriodicMeasurement(); err != nil {
			println("could not start peridic measurement:", err)
			return
		}
		co2Sensor = sensor
	}
	{
		sensor := sgp30.New(bus)
		if !sensor.Connected() {
			println("VOC sensor not connected")
			return
		}
		err := sensor.Configure(sgp30.Config{})
		if err != nil {
			println("VOC sensor could not be configured:", err.Error())
			return
		}
		vocSensor = sensor
	}
}

func sampleCO2Sensor(sensor *scd4x.Device) {
	co2, err := sensor.ReadCO2()
	if err != nil {
		println("failed to read CO2:", err.Error())
		return
	}

	temperature, err := sensor.ReadTemperature()
	if err != nil {
		println("failed to read temperature:", err.Error())
		return
	}

	humidity, err := sensor.ReadHumidity()
	if err != nil {
		println("failed to read humidity:", err.Error())
		return
	}
	println("CO2:        ", co2)
	println("temperature:", temperature)
	println("humidity:   ", humidity)
}

func sampleVOCSensor(sensor *sgp30.Device) {
	err := sensor.Update(0)
	if err != nil {
		println("could not read VOC sensor:", err.Error())
		return
	}
	println("CO2eq:      ", sensor.CO2())
	println("TVOC        ", sensor.TVOC())
}
