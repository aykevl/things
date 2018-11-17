package main

import (
	"machine"
	"time"

	"github.com/aykevl/fixpoint"
	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygo-drivers/mpu6050"
	"github.com/aykevl/tinygo-drivers/ws2812"
)

const ledPin = 21

var (
	whole = fixpoint.Q24FromFloat(0.8506508083520399)
	grinv = fixpoint.Q24FromFloat(0.5257311121191336)
	zero  = fixpoint.Q24FromFloat(0)
)

const numLeds = 12

// The LED vectors from the center of the globe.
var (
	vectors = [numLeds]fixpoint.Vec3Q24{
		{zero, grinv, whole},
		{zero, grinv.Neg(), whole},
		{whole.Neg(), zero, zero},
		{grinv.Neg(), whole, zero},
		{grinv, whole, zero},
		{whole, zero, grinv},
		{grinv, whole.Neg(), zero},
		{grinv.Neg(), whole.Neg(), zero},
		{whole.Neg(), zero, grinv.Neg()},
		{zero, grinv, whole.Neg()},
		{whole, zero, grinv.Neg()},
		{zero, grinv.Neg(), whole.Neg()},
	}
)

func main() {
	// Initialize LEDs
	pin := machine.GPIO{ledPin}
	pin.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})
	strip := ws2812.New(pin)
	colors := ledsgo.Strip(make([]uint32, numLeds))
	colors.FillSolid(0) // turn off all LEDs
	strip.WriteColors(colors)

	// Initialize MPU6050
	time.Sleep(10 * time.Millisecond) // needs a bit of startup time
	bus := machine.I2C0
	bus.Configure(machine.I2CConfig{
		SCL: 4,
		SDA: 3,
	})
	mpu := mpu6050.NewI2C(bus)

	// Run animation
	//noise(colors, strip)
	//rotate(colors, strip)
	//gradient(colors, strip)
	mpuTest2(colors, strip, mpu)
}

func configureMPU(colors ledsgo.Strip, strip ws2812.Device, mpu mpu6050.Device) {
	if !mpu.Connected() {
		flashColor(colors, strip, 0xff0000)
	}

	mpu.Configure()
}

// flashColor sets the color of the globe to the given color for one second.
// Useful for debugging.
func flashColor(colors ledsgo.Strip, strip ws2812.Device, color uint32) {
	colors.FillSolid(color)
	strip.WriteColors(colors)
	time.Sleep(time.Second)
}

// Random noise.
func noise(colors ledsgo.Strip, strip ws2812.Device) {
	for {
		fillNoise(colors)
		strip.WriteColors(colors)
		time.Sleep(time.Millisecond)
	}
}

func fillNoise(colors ledsgo.Strip) {
	now := time.Now().UnixNano()
	for i := range colors {
		x := vectors[i].X.Int32Scaled(255) + int32(now>>20)
		y := vectors[i].Y.Int32Scaled(255)
		z := vectors[i].Z.Int32Scaled(255)
		hue := uint16(ledsgo.Noise3(x, y, z))
		colors[i] = ledsgo.Color{hue, 0xff, 0x22}.Spectrum()
	}
}

// Map positive acceleration to R, G, B colors.
func mpuTest1(colors ledsgo.Strip, strip ws2812.Device, mpu mpu6050.Device) {
	configureMPU(colors, strip, mpu)

	for {
		x, y, z := mpu.ReadAcceleration()
		var r, g, b uint32
		if x > 0 {
			r = uint32(x >> 8)
		}
		if y > 0 {
			g = uint32(y >> 8)
		}
		if z > 0 {
			b = uint32(z >> 8)
		}
		colors.FillSolid(r<<16 | g<<8 | b)
		strip.WriteColors(colors)
		time.Sleep(100 * time.Millisecond)
	}
}

// Gradient sets all LEDs to a red-blue gradient across the globe.
func gradient(colors ledsgo.Strip, strip ws2812.Device) {
	for {
		for i := 0; i < len(vectors); i++ {
			vec := vectors[i]
			const brightness = 32
			val := vec.X.Int32Scaled(brightness)
			r := uint32(brightness + val)
			b := uint32(brightness - val)
			colors[i] = r<<16 | b
		}
		strip.WriteColors(colors)
		time.Sleep(10 * time.Millisecond)
	}
}

// Rotate sets a rotating gradient across all LEDs.
func rotate(colors ledsgo.Strip, strip ws2812.Device) {
	// Set up rotation
	inc := fixpoint.QuatQ24{fixpoint.Q24FromInt32(0), fixpoint.Vec3Q24FromFloat(0, 0, 0.005)}
	inc.W = fixpoint.Q24FromInt32(1).Sub(fixpoint.Q24FromFloat(0.5).Mul(inc.X().Mul(inc.X()).Add(inc.Y().Mul(inc.Y())).Add(inc.Z().Mul(inc.Z()))))

	// How far we've rotated so far.
	rot := fixpoint.QuatIdent()
	for cycle := 0; ; cycle++ {
		rot = rot.Mul(inc)
		for i := 0; i < len(vectors); i++ {
			vec := rot.Rotate(vectors[i])
			const brightness = 32
			val := vec.X.Int32Scaled(brightness)
			r := uint32(brightness + val)
			b := uint32(brightness - val)
			colors[i] = r<<16 | b
		}
		strip.WriteColors(colors)
		time.Sleep(10 * time.Millisecond)
	}
}

// mpuTest2 sets a linear gradient across the globe that should be stable when
// it gets rotated.
func mpuTest2(colors ledsgo.Strip, strip ws2812.Device, mpu mpu6050.Device) {
	configureMPU(colors, strip, mpu)

	// How far we've rotated so far.
	orientation := fixpoint.QuatIdent()
	for {
		x, y, z := mpu.ReadRotation()

		// Create a rotation quaternion from the rotation we've just read.
		const mul = 15
		inc := fixpoint.QuatQ24{
			fixpoint.Q24FromInt32(0), // W
			fixpoint.Vec3Q24{ // V
				fixpoint.Q24{int32(z) * mul}, // V.X
				fixpoint.Q24{int32(y) * mul}, // V.Y
				fixpoint.Q24{int32(x) * mul}, // V.Z
			},
		}
		inc.W = fixpoint.Q24FromInt32(1).Sub(fixpoint.Q24FromFloat(0.5).Mul(inc.X().Mul(inc.X()).Add(inc.Y().Mul(inc.Y())).Add(inc.Z().Mul(inc.Z()))))

		// Update estimated orientation.
		orientation = orientation.Mul(inc)

		// Update LEDs with.
		for i := 0; i < len(vectors); i++ {
			vec := orientation.Rotate(vectors[i])
			const brightness = 32
			val := vec.X.Int32Scaled(brightness)
			r := uint32(brightness + val)
			b := uint32(brightness - val)
			colors[i] = r<<16 | b
		}
		strip.WriteColors(colors)

		time.Sleep(10 * time.Millisecond)
	}
}
