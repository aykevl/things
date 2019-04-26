package main

import (
	"device/arm"
	"image/color"
	"machine"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tinygo-org/drivers/lis3dh"
	"github.com/tinygo-org/drivers/ws2812"
)

var (
	i2c    = machine.I2C1
	ledPin = machine.GPIO{machine.NEOPIXELS}
	ws     = ws2812.New(ledPin)
	leds   = make([]color.RGBA, 10)
)

// Coordinates of the LEDs in (x, y) format
// (assuming they're on a circle with 12 points).
var coords = [10][2]float32{
	// first half
	{0.500, -0.866},
	{0.866, -0.500},
	{1.000, 0.000},
	{0.866, 0.500},
	{0.500, 0.866},
	{-0.500, 0.866},

	// second half
	{-0.866, 0.500},
	{-1.000, 0.000},
	{-0.866, -0.500},
	{-0.500, -0.866},
}

func main() {
	ledPin.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})

	i2c.Configure(machine.I2CConfig{})
	accel := lis3dh.New(i2c)
	accel.Address = lis3dh.Address1 // address on the Circuit Playground Express
	accel.Configure()
	accel.SetRange(lis3dh.RANGE_2_G)

	println("connected:", accel.Connected())

	for {
		time.Sleep(time.Second / 20)

		xi, yi, zi, _ := accel.ReadAcceleration()
		println("accel:", xi/1000, yi/1000, zi/1000)

		vec := mgl32.Vec3{
			float32(xi) / 1000000,
			float32(yi) / 1000000,
			float32(zi) / 1000000,
		}
		vec = vec.Normalize()
		println("  x:", int32(vec[0]*1000))
		println("  y:", int32(vec[1]*1000))

		for i := range leds {
			ledX := coords[i][0]
			ledY := coords[i][1]
			height := (vec[0] * ledX) + (vec[1] * ledY)
			red := (height * 10) + 1
			if red < 0 {
				red = 0
			}
			leds[i] = color.RGBA{uint8(red), 0, 0, 0}
		}

		mask := arm.DisableInterrupts()
		ws.WriteColors(leds)
		arm.EnableInterrupts(mask)
	}
}
