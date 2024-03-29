package main

import (
	"device/nrf"
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
	nil,
	solid,
	noise,
	fire,
	glitter,
	iris,
	arrows,
	colorarrows,
	halfcircles,
}

func main() {
	if debug {
		if serialTxPin != machine.NoPin {
			// Reconfigure UART on a different pin.
			machine.UART0.Configure(machine.UARTConfig{
				TX: serialTxPin,
				RX: machine.NoPin,
			})
		}
	} else {
		// The UART consumes a lot of current in sleep mode when left enabled.
		nrf.UART0.ENABLE.Set(0)
	}

	if debug {
		println("starting")
	}
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

	// Default values, for when the MPU isn't available.
	m := movement{
		rotationSpeed: 360 * 30 / 2, // 1½ rotation per second
	}

	now := time.Now()
	for i := uint(0); ; i++ {
		if animationIndex == 0 {
			// Disable all LEDs (for compatibility with poi that don't have a
			// MOSFET to control LED power).
			for y := range leds {
				leds[y] = color.RGBA{}
			}
			a.WriteColors(leds)
			time.Sleep(time.Millisecond) // wait until data is fully sent (?)

			// Turn off many peripherals.
			disable()

			for animationIndex == 0 {
				// Wait until we're not sleeping anymore.
				if debug {
					println("sleeping")
				}
				time.Sleep(time.Second)
			}

			// Enable all peripherals again.
			enable()
		}

		lastTime := now
		now = time.Now()

		if imu != nil {
			// Ignore gyroX because it just shows how much the poi rotates around
			// its own axis and is thus noise. The other two (gyroY and gyroZ) show
			// how fast the poi is actually spinning.
			_, gyroY, gyroZ, _ := imu.ReadRotation()

			// The gyro number here is in units of .1°/s.
			m.rotationSpeed = ledsgo.Sqrt(int(((gyroY / 100000) * (gyroY / 100000)) + ((gyroZ / 100000) * (gyroZ / 100000))))
		}

		const minAnimationSpeed = 50

		// Update position.
		timeElapsed := now.Sub(lastTime)
		animationSpeed := m.rotationSpeed
		if animationSpeed < minAnimationSpeed {
			animationSpeed = minAnimationSpeed
		}
		m.animationPosition += (int64(timeElapsed) * int64(animationSpeed)) >> (27 - speed)

		animation := animations[animationIndex]
		if animation == nil {
			// Avoid a race condition. If animationIndex is set to zero between
			// the sleep loop above and now, the animation will be the first one
			// and it'll result in a nil panic.
			// Still technically a race condition but hopefully not an actually
			// buggy one.
			continue
		}
		animation(now, m)

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

// applyBrightness applies the lower 3 bits of brightness control to the color.
// Normally, the SK9822 LEDs only support 5 bits of brightness, therefore the
// lower 3 bits are wasted. However, lowering the brightness further is still
// very useful in low-light situations. Therefore, this function applies the
// brightness into the color itself, to regain more brightness control at a loss
// of a bit of fidelity.
func applyBrightness(c color.RGBA) color.RGBA {
	if c.A < 8 {
		c.R = uint8(uint32(c.R) * uint32(c.A) / 8)
		c.G = uint8(uint32(c.G) * uint32(c.A) / 8)
		c.B = uint8(uint32(c.B) * uint32(c.A) / 8)
		c.A = 8
	}
	return c
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

	// Custom value that indicates the position in an animation, in a somewhat
	// undefined way. It increments quickly as the poi moves fast but keeps
	// incrementing when the poi moves slowly. This is a replacement for basing
	// an animation entirely on the current time and is usually better.
	animationPosition int64
}

// Colorful noise.
func noise(now time.Time, m movement) {
	const spread = 7

	// Color each pixel.
	for y := int16(0); y < height; y++ {
		hue := ledsgo.Noise2(uint32(m.animationPosition>>10), uint32(y<<spread)) * 2
		c := ledsgo.Color{hue, 0xff, 0xff}.Spectrum()
		c.A = baseColor.A
		setLED(y, c)
	}
}

// Looks a bit like spikes from inside to the outside.
func iris(now time.Time, m movement) {
	expansion := int16(ledsgo.Noise1(uint32(m.animationPosition>>7)) / 256) - 50
	for y := int16(0); y < height; y++ {
		intensity := expansion - y*224/height
		if intensity < 0 {
			intensity = 0
		}
		c := ledsgo.ApplyAlpha(baseColor, uint8(intensity))
		c.A = baseColor.A
		setLED(y, c)
	}
}

var rainbowColors = []color.RGBA{
	color.RGBA{255, 0, 0, 255},   // red
	color.RGBA{233, 22, 0, 255},  // orange
	color.RGBA{180, 75, 0, 255},  // yellow
	color.RGBA{0, 255, 0, 255},   // green
	color.RGBA{0, 0, 255, 255},   // blue
	color.RGBA{128, 0, 127, 255}, // purple
}

// Make a series of arrows in rainbow colors.
func colorarrows(now time.Time, m movement) {
	for y := int16(0); y < height/2; y++ {
		position := m.animationPosition>>17 - int64(y)
		c := rainbowColors[(position/(height/2))%int64(len(rainbowColors))]
		c.A = baseColor.A
		if position%(height/2) == 0 {
			// Make the position between the arrow colors black.
			c = color.RGBA{}
		}
		setLED(y, c)
		setLED(height-1-y, c)
	}
}

// Somewhat improperly named. When two poi are spinning in opposite direction,
// it has a somewhat palm tree like appearance.
func halfcircles(now time.Time, m movement) {
	chosenOne := int16((m.animationPosition >> 16) % height)
	for y := int16(0); y < height; y++ {
		c := color.RGBA{}
		if y >= chosenOne && y < chosenOne+(height/7) {
			c = baseColor
		}
		setLED(y, c)
	}
}

// Simple > arrows pointing in the direction the poi is moving.
func arrows(now time.Time, m movement) {
	// First make them all black.
	for y := int16(0); y < height; y++ {
		setLED(y, color.RGBA{})
	}

	// Turn the two LEDs on that are part of the arrow.
	index := int16((m.animationPosition >> 17) % (height / 2))
	setLED(index, baseColor)
	setLED(height-1-index, baseColor)
}

// Avoid heap allocation by allocating this globally.
var glitterBuf [4]byte

// Random colored specles. Looks great in the dark because the poi itself is
// (nearly) invisible showing only these speckles.
func glitter(now time.Time, m movement) {
	// Make all LEDs black.
	for y := int16(0); y < height; y++ {
		setLED(y, color.RGBA{})
	}

	// Get a random number based on the time.
	t := uint32(m.animationPosition >> 17)
	glitterBuf[0] = byte(t)
	glitterBuf[1] = byte(t >> 8)
	glitterBuf[2] = byte(t >> 16)
	glitterBuf[3] = byte(t >> 24)
	hash := murmur3.Sum32(glitterBuf[:])

	// Use this number to get an index.
	index := int16(hash % height)

	// Get a random color on the color wheel.
	c := ledsgo.Color{uint16((hash >> 7)), 0xff, 0xff}.Spectrum()
	c.A = baseColor.A

	setLED(index, c)
}

// Solid color. Useful to reduce distraction, for testing and as a not too
// distracting startup color.
func solid(now time.Time, m movement) {
	for y := int16(0); y < height; y++ {
		setLED(y, baseColor)
	}
}

// Fire animation. The flame is the configured color, so you can not only have a
// red flame, but also a green or blue flame. Red looks the best IMHO, though.
func fire(now time.Time, m movement) {
	var cooling = (14 * 16) / height // higher means faster cooling
	const detail = 400               // higher means more detailed flames
	for y := int16(0); y < height; y++ {
		heat := int16(ledsgo.Noise2(uint32(m.animationPosition>>8), uint32(y)*detail)/256)
		heat -= y * int16(cooling)
		if heat < 0 {
			heat = 0
		}
		c := coloredFlame(uint8(heat))
		c.A = baseColor.A
		setLED(y, c)
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
