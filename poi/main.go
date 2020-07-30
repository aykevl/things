package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"github.com/spaolacci/murmur3"
	"tinygo.org/x/drivers/apa102"
	"tinygo.org/x/drivers/bmi160"
)

var leds = make([]color.RGBA, numLeds)

var imu *bmi160.DeviceSPI

const debug = false

// Parameters that are controlled with Bluetooth.
var (
	animationIndex uint8 = 0
	speed          uint8 = 10
)

// List of animations that can be selected over BLE.
var animations = []func(time.Time, movement){
	solid,
	new(noiseState).noise,
	fire,
	iris,
	gear,
	halfcircles,
	arrows,
	glitter,
	black,
}

func main() {
	if serialTxPin != machine.NoPin {
		// Reconfigure UART on a different pin.
		machine.UART0.Configure(machine.UARTConfig{
			TX: serialTxPin,
			RX: machine.NoPin,
		})
	}
	println("starting")
	initHardware()

	if hasBMI160 {
		machine.SPI1.Configure(machine.SPIConfig{
			SCK:       6, // connected to SCX of BMI160 (also SCK)
			SDO:       5, // connected to SDX of BMI160 (also SDI)
			SDI:       4, // connected to SDO of BMI160
			Mode:      0, // both mode 0 and mode 3 are supported
			Frequency: 8e6,
		})
		csb := machine.Pin(7)
		imu = bmi160.NewSPI(csb, machine.SPI1)
	}

	if mosfetPin != machine.NoPin {
		mosfetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	enable()

	a := apa102.New(machine.SPI0)

	for i := uint(0); ; i++ {
		if int(animationIndex) == len(animations)-1 {
			// Disable all LEDs (for compatibility with poi that don't have a
			// MOSFET to control LED power).
			for y := range leds {
				leds[y] = color.RGBA{}
			}
			a.WriteColors(leds)
			time.Sleep(time.Millisecond) // wait until data is fully sent (?)

			// Turn off many peripherals.
			disable()

			for int(animationIndex) == len(animations)-1 {
				// Wait until we're not sleeping anymore.
				if debug {
					println("sleeping")
				}
				time.Sleep(time.Second)
			}

			// Enable all peripherals again.
			enable()
		}

		now := time.Now()

		// Default values, for when the MPU isn't available.
		m := movement{
			rotationSpeed: 54000, // 1½ rotation per second
		}

		if imu != nil {
			// Ignore gyroX because it just shows how much the poi rotates around
			// its own axis and is thus noise. The other two (gyroY and gyroZ) show
			// how fast the poi is actually spinning.
			_, gyroY, gyroZ, _ := imu.ReadRotation()

			// The gyro number here is in units of .01°/s.
			m.rotationSpeed = ledsgo.Sqrt(int(((gyroY / 10000) * (gyroY / 10000)) + ((gyroZ / 10000) * (gyroZ / 10000))))
		}

		animation := animations[animationIndex]
		animation(now, m)

		if height != len(leds) {
			// This is a poi with two sides that wraps around.
			// Make sure the other side is also colored properly (animations
			// only color one side).
			for i := 0; i < len(leds)/2; i++ {
				leds[len(leds)-i-1] = leds[i]
			}
		}
		a.WriteColors(leds)

		// print speed
		if debug && i%500 == 0 {
			//duration := time.Since(now)
			//println("duration:", duration.String())
		}
	}
}

// enable enables all peripherals that might be disabled when in standby mode.
func enable() {
	// Enable LEDs.
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: spiFrequency,
		Mode:      0,
		SCK:       spiClockPin,
		SDO:       spiDataPin,
		SDI:       machine.NoPin,
	})
	if mosfetPin != machine.NoPin {
		mosfetPin.Set(true)
	}

	// Enable IMU.
	if imu != nil {
		imu.Configure()
	}
}

// disable turns off peripherals that can be disabled to save power.
// Perhaps most importantly, it turns off all LEDs to massively reduce current
// consumption.
func disable() {
	// Disable LEDs.
	machine.SPI0.Bus.ENABLE.Set(0)
	if mosfetPin != machine.NoPin {
		mosfetPin.Set(false)
	}

	// Disable IMU.
	if imu != nil {
		imu.Reset()
	}
}

// Type movement contains information about the current movement that might be
// relevant to animations. Some animations may change depending on how the poi
// moves.
type movement struct {
	// Rotation speed in 0.01°/s.
	// A poi will usually spin at 1-2 rotations per second, which means the
	// rotation speed will usually be in the range of 36000 to 72000 when
	// rotating.
	// This number is either zero or positive (although exactly zero is
	// unlikely).
	rotationSpeed int
}

// State for the noise function.
type noiseState struct {
	lastTime      time.Time
	noisePosition int64
}

// Colorful noise.
func (n *noiseState) noise(now time.Time, m movement) {
	const spread = 7
	const minAnimationSpeed = 1000

	// Update position.
	timeElapsed := now.Sub(n.lastTime)
	animationSpeed := m.rotationSpeed
	if animationSpeed < minAnimationSpeed {
		animationSpeed = minAnimationSpeed
	}
	n.noisePosition += (int64(timeElapsed) * int64(animationSpeed)) >> (32 - speed)
	n.lastTime = now

	// Color each pixel.
	for y := int16(0); y < height; y++ {
		hue := uint16(ledsgo.Noise2(int32(n.noisePosition>>10), int32(y<<spread))) * 2
		c := ledsgo.Color{hue, 0xff, 0xff}.Spectrum()
		c.A = baseColor.A
		leds[y] = c
	}
}

// Looks a bit like spikes from inside to the outside.
func iris(now time.Time, m movement) {
	expansion := (ledsgo.Noise1(int32(now.UnixNano()>>(21-speed))) / 256) + 128 - 50
	for y := int16(0); y < height; y++ {
		intensity := expansion - y*224/height
		if intensity < 0 {
			intensity = 0
		}
		c := ledsgo.ApplyAlpha(baseColor, uint8(intensity))
		c.A = baseColor.A
		leds[y] = c
	}
}

// Looks like a typical blocky gear, with square gear teeth.
func gear(now time.Time, m movement) {
	long := int16((now.UnixNano()>>(32-speed))%8) == 0
	for y := int16(0); y < height; y++ {
		c := color.RGBA{}
		if long || y < height/4 {
			c = baseColor
		}
		leds[y] = c
	}
}

// Somewhat improperly named. When two poi are spinning in opposite direction,
// it has a somewhat palm tree like appearance.
func halfcircles(now time.Time, m movement) {
	chosenOne := int16((now.UnixNano() >> (32 - speed)) % height)
	for y := int16(0); y < height; y++ {
		c := color.RGBA{}
		if y >= chosenOne && y < chosenOne+(height/7) {
			c = baseColor
		}
		leds[y] = c
	}
}

// Simple > arrows pointing in the direction the poi is moving.
func arrows(now time.Time, m movement) {
	// First make them all black.
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}

	// Turn the two LEDs on that are part of the arrow.
	index := int16((now.UnixNano() >> (32 - speed)) % (height / 2))
	leds[index] = baseColor
	leds[height-1-index] = baseColor
}

// Random colored specles. Looks great in the dark because the poi itself is
// (nearly) invisible showing only these speckles.
func glitter(now time.Time, m movement) {
	// Make all LEDs black.
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}

	// Get a random number based on the time.
	t := uint32(now.UnixNano() >> (32 - speed))
	hash := murmur3.Sum32([]byte{byte(t), byte(t >> 8), byte(t >> 16), byte(t >> 24)})

	// Use this number to get an index.
	index := hash % (height * 2)
	if index >= height {
		return // don't sparkle all the time
	}

	// Get a random color on the color wheel.
	c := ledsgo.Color{uint16((hash >> 7)), 0xff, 0xff}.Spectrum()
	c.A = baseColor.A

	leds[index] = c
}

// Solid color. Useful to reduce distraction, for testing and as a not too
// distracting startup color.
func solid(now time.Time, m movement) {
	for y := int16(0); y < height; y++ {
		leds[y] = baseColor
	}
}

// Disable all LEDs. LEDs will still consume power as they use around 1mA even
// when they're dark.
func black(now time.Time, m movement) {
	for y := int16(0); y < height; y++ {
		leds[y] = color.RGBA{}
	}
}

// Fire animation. The flame is the configured color, so you can not only have a
// red flame, but also a green or blue flame. Red looks the best IMHO, though.
func fire(now time.Time, m movement) {
	var cooling = (14 * 16) / height // higher means faster cooling
	const detail = 400               // higher means more detailed flames
	for y := int16(0); y < height; y++ {
		heat := ledsgo.Noise2(int32((now.UnixNano()>>(23-speed))), int32(y)*detail)/256 + 128
		heat -= y * int16(cooling)
		if heat < 0 {
			heat = 0
		}
		c := coloredFlame(uint8(heat))
		c.A = baseColor.A
		leds[y] = c
	}
}

// Colored flame. Like a heat map, but the lowest temperatures are not fixed red
// but instead use the configured color.
func coloredFlame(index uint8) color.RGBA {
	if index < 128 {
		// <color>
		c := ledsgo.ApplyAlpha(baseColor, index*2)
		c.A = 255
		return c
	}
	if index < 224 {
		// <color>-yellow
		c1 := ledsgo.ApplyAlpha(baseColor, 255-uint8(uint32(index-128)*8/3))
		c2 := ledsgo.ApplyAlpha(color.RGBA{255, 255, 0, 255}, uint8(uint32(index-128)*8/3))
		return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 255}
		//return color.RGBA{255, uint8(uint32(index-128) * 8 / 3), 0, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}

// Heat map, where lower numbers are colder. Previously used in the fire
// animation.
func heatMap(index uint8) color.RGBA {
	if index < 128 {
		return color.RGBA{index * 2, 0, 0, 255}
	}
	if index < 224 {
		// red-yellow
		return color.RGBA{255, uint8(uint32(index-128) * 8 / 3), 0, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}
