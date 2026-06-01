package main

import (
	"image/color"

	"github.com/aykevl/ledsgo"
)

const numLEDs = 36

const initialMode = 0 // first mode

const (
	modeTrace = iota
	modeNoise
	modeFire
	modeFireSound
	modeStatic
	modeFlag
	modeSparkle
	modeVUMeter
	modeLast

	modeTest
	modePowerOn
)

var variantsPerMode = [...]uint8{
	modeTrace:     3,
	modeNoise:     uint8(len(noisePatterns)),
	modeFire:      6,
	modeFireSound: 6,
	modeStatic:    uint8(len(staticColors)),
	modeVUMeter:   uint8(len(staticColors)),
	modeFlag:      uint8(len(allFlags)),
	modeSparkle:   uint8(len(sparkleVariants)),
}

// Cycle to the next variant within a mode.
func animationNextVariant(mode, variant int) int {
	variant++
	if mode >= len(variantsPerMode) {
		variant = 0 // out of range
	} else if variant >= int(variantsPerMode[mode]) {
		variant = 0
	}
	return variant
}

func animate(mode, variant, led, frame int) Color {
	switch mode {
	case modeTrace:
		return trace(led, frame, variant)
	case modeNoise:
		return noise(led, frame, variant)
	case modeFire:
		return fire(led, frame, variant)
	case modeFireSound:
		return fireSound(led, frame, variant)
	case modeStatic:
		return staticColor(led, frame, variant)
	case modeFlag:
		return showFlag(led, frame, variant)
	case modeVUMeter:
		return vuMeter(led, frame, variant)
	case modeSparkle:
		return makeSparkle(led, frame, variant)
	case modeTest:
		return testPulse(led, frame)
	case modePowerOn:
		return powerOn(led, frame)
	default:
		// bug
		return errorPattern(led, frame)
	}
}

func animationNeedsMic(mode int) bool {
	switch mode {
	case modeVUMeter, modeFireSound, modeSparkle:
		return true
	default:
		return false
	}
}

// Called on every new frame, to do per-frame processing.
// It should be very fast, otherwise the earrings will dip in brightness at 30Hz
// or so which looks annoying.
func newFrame(mode, variant int) {
	switch mode {
	case modeVUMeter:
		addPower(uint16(processSamples()))
	case modeFireSound:
		addPower(uint16(processSamples()))
		vol := uint8(min(currentVolume()*3/512, 0xff))
		volumeHistoryIndex++
		if volumeHistoryIndex == 36 {
			volumeHistoryIndex = 0
		}
		volumeHistory[volumeHistoryIndex] = vol
	case modeSparkle:
		sparkleNextFrame(variant)
	}
}

var noisePatterns = [...]struct {
	speed   uint8            // higher means faster
	spread  uint8            // higher means colors more close together
	palette ledsgo.Palette16 // FastLED palette
}{
	{10, 11, ledsgo.PartyColors},
	{10, 11, ledsgo.RainbowColors},
	{9, 12, ledsgo.ForestColors},
	{9, 12, ledsgo.OceanColors},
	{9, 12, ledsgo.CloudColors},
}

// Show some Simplex noise on the earring, with various predefined patterns.
func noise(led, frame, variant int) Color {
	// Determine palette to show.
	if variant >= len(noisePatterns) {
		return NewColor(0, 0, 0) // shouldn't happen
	}
	pattern := &noisePatterns[variant]

	x := uint32(frame) << uint32(pattern.speed)
	y := uint32(led) << uint32(pattern.spread)
	c := pattern.palette.ColorAt(ledsgo.Perlin2(x, y) * 2)
	return NewColor(c.R, c.G, c.B)
}

var traceIndex uint8
var traceIndexFrame int

func trace(led, frame, variant int) Color {
	switch variant {
	case 0, 1:
		return twoTracers(led, frame, variant)
	default:
		return purpleCircles(led, frame)
	}
}

func twoTracers(led, frame, variant int) Color {
	// This essentially calculates (frame%72) but faster.
	if frame != traceIndexFrame {
		traceIndexFrame = frame
		traceIndex++
	}
	if traceIndex >= 72 {
		traceIndex = 0
	}

	// Calculate where in the circle we are for this LED.
	index := uint8(led*2) + 71 - traceIndex
	if index >= 72 {
		index -= 72
	}
	if index >= 72 {
		index -= 72
	}

	// Split out in the two traces.
	trace := 0
	if index >= 36 {
		index -= 36
		trace = 1
	}

	// Determine the color for the LED.
	var c Color
	switch variant {
	case 0:
		// Rainbow trace.
		var col color.RGBA
		switch trace {
		case 0:
			// First tracer.
			col = ledsgo.Color{H: uint16(frame * 128), S: 255, V: 255}.Rainbow()
		default:
			// Second tracer, offset 180° on the color wheel and on the actual LED ring.
			col = ledsgo.Color{H: uint16(frame*128) + 0x8000, S: 255, V: 255}.Rainbow()
		}
		c = NewColor(
			uint8(uint32(col.R)*(uint32(index)*7)/256),
			uint8(uint32(col.G)*(uint32(index)*7)/256),
			uint8(uint32(col.B)*(uint32(index)*7)/256))
	case 1:
		// Red (fire) and blue (ice) swirling around in circles.
		switch trace {
		case 0:
			c = NewColor(index*7, index*3/8, 0)
		default:
			c = NewColor(index, 0, index*7)
		}
	}

	// Dim at the start (fade in)
	if index == 35 {
		c = NewColor(c.R()/3, c.G()/3, c.B()/3)
	}

	return c
}

// Three purple tracers running in circles.
func purpleCircles(led, frame int) Color {
	if frame != traceIndexFrame {
		traceIndexFrame = frame
		traceIndex++
	}
	if traceIndex >= 24 {
		// This animation has three tracers.
		traceIndex = 0
	}

	index := uint8(led*2) + 23 - traceIndex
	if index >= 24 {
		index -= 24
	}
	if index >= 24 {
		index -= 24
	}
	if index >= 24 {
		index -= 24
	}

	if index == 23 {
		// fade in at the beginning
		return NewColor(0x20, 0, 0x10)
	}
	return NewColor(uint8(index*10), 0, uint8(index*5))
}

// Determine fire color for both fire animations.
func fireColor(variant int) color.RGBA {
	switch variant {
	case 0: // red
		return color.RGBA{R: 255}
	case 1: // orange
		return color.RGBA{R: 220, G: 40}
	case 2: // green
		return color.RGBA{G: 255}
	case 3: // teal-ish
		return color.RGBA{G: 127, B: 127}
	case 4: // blue
		return color.RGBA{B: 255}
	default: // purple
		return color.RGBA{R: 127, B: 127}
	}
}

// Fire animation in various colors.
func fire(led, frame, variant int) Color {
	intensityIndex := indexFromBottom(led)
	noiseIndex := uint32(frame) - uint32(intensityIndex)
	if led > numLEDs/2 {
		// Use a different flame on the otherside of the earring.
		// Without this, it would simply mirror the flame on both sides.
		// The constant is just an arbitrary number to give it enough distance
		// from the left side.
		noiseIndex += 0x1234
	}

	// Calculate the amount of heat on this particular pixel.
	heat := uint8(ledsgo.Noise1(noiseIndex<<10) >> 8)
	cooling := uint8(intensityIndex * 11)
	if heat < uint8(cooling) {
		heat = 0
	} else {
		heat -= uint8(cooling)
	}

	// Turn it into a flame, based on the given palette.
	// Perhaps we could use an actual 0-255 (or 0-64) palette instead? That
	// might be faster.
	c := coloredFlame(heat, fireColor(variant))
	return NewColor(c.R, c.G, c.B)
}

var volumeHistory [36]uint8
var volumeHistoryIndex uint8

func fireSound(led, frame, variant int) Color {
	// LED indices calculated, so they start from the bottom and alternate
	// between the two sides:
	// LED 14: index 6
	// LED 15: index 4
	// LED 16: index 2
	// LED 17: index 0
	// LED 18: index 1
	// LED 18: index 3
	// LED 18: index 5
	var historyIndex int
	var intensityIndex uint8
	if led < 18 {
		intensityIndex = uint8(18 - led)
		historyIndex = (17 - led) * 2
	} else {
		historyIndex = (led-18)*2 + 1
		intensityIndex = uint8(led - 18)
	}

	// Calculate index into volumeHistory, since that is a circular buffer (with
	// index 0 being the most recent and higher numbers being further in the
	// past).
	historyIndex = int(volumeHistoryIndex) - historyIndex
	if historyIndex < 0 {
		historyIndex += 36
	}

	// Cool down the flame the higher it gets.
	heat := volumeHistory[historyIndex]
	cooling := uint8(intensityIndex * 8)
	if heat < uint8(cooling) {
		heat = 0
	} else {
		heat -= uint8(cooling)
	}

	c := coloredFlame(heat, fireColor(variant))
	return NewColor(c.R, c.G, c.B)
}

// Colored flame. Like a heat map, but the lowest temperatures are not fixed red
// but instead use the configured color.
//
//go:inline
func coloredFlame(index uint8, fireColor color.RGBA) color.RGBA {
	if index < 128 {
		// <color>
		c := ledsgo.ApplyAlpha(fireColor, index*2)
		c.A = 255
		return c
	}
	if index < 224 {
		// <color>-yellow
		c1 := ledsgo.ApplyAlpha(fireColor, 255-uint8(uint32(index-128)*8/3))
		c2 := ledsgo.ApplyAlpha(color.RGBA{255, 255, 0, 255}, uint8(uint32(index-128)*8/3))
		return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 255}
	}
	// yellow-white
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}

func indexFromBottom(index int) int {
	// Start at the top with 18, move along the right size down to 0 at the
	// bottom, and then resume counting upwards again:
	// 18, 17, ..., 1, 0, 1, 2, ..., 17
	newIndex := 18 - index
	if index >= 18 {
		newIndex = index - 18
	}
	return newIndex
}

var (
	staticColors = [...]Color{
		NewColor(0x99, 0x00, 0x66), // purple
		NewColor(0xee, 0x00, 0x11), // hot pink
		NewColor(0xff, 0x00, 0x00), // red
		NewColor(0xcc, 0x22, 0x00), // orange
		NewColor(0x88, 0x77, 0x00), // yellow (orange-ish)
		NewColor(0x33, 0xcc, 0x00), // yellow (green-ish)
		NewColor(0x00, 0xff, 0x00), // green
		NewColor(0x00, 0xdd, 0x22), // green-turquoise
		NewColor(0x00, 0x99, 0x66), // turquoise
		NewColor(0x00, 0x33, 0xcc), // greenish blue
		NewColor(0x00, 0x00, 0xff), // blue
		NewColor(0x11, 0x00, 0xee), // purple-ish blue
		NewColor(0x55, 0x00, 0xaa), // blueish purple
	}
)

func staticColor(led, frame, variant int) Color {
	if variant >= len(staticColors) {
		return NewColor(0, 0, 0) // should be unreachable
	}
	return staticColors[variant]
}

type Palette [numLEDs]Color

var (
	flagRed    = NewColor(0xff, 0x00, 0x00)
	flagOrange = NewColor(0xff, 0x55, 0x00)
	flagYellow = NewColor(0x88, 0xff, 0x00)
	flagGreen  = NewColor(0x00, 0xff, 0x00)
	flagBlue   = NewColor(0x00, 0x00, 0xff)
	flagPurple = NewColor(0x80, 0x00, 0x80)
	flagLGBT   = Palette{
		flagRed, flagRed, flagRed, flagRed, flagRed, flagRed,
		flagOrange, flagOrange, flagOrange, flagOrange, flagOrange, flagOrange,
		flagYellow, flagYellow, flagYellow, flagYellow, flagYellow, flagYellow,
		flagGreen, flagGreen, flagGreen, flagGreen, flagGreen, flagGreen,
		flagBlue, flagBlue, flagBlue, flagBlue, flagBlue, flagBlue,
		flagPurple, flagPurple, flagPurple, flagPurple, flagPurple, flagPurple,
	}

	flagPastelBlue  = NewColor(0x11, 0x33, 0x88)
	flagPastelPink  = NewColor(0x80, 0x22, 0x22)
	flagPastelWhite = NewColor(0x88, 0xaa, 0xaa)
	flagTrans       = Palette{
		flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelWhite, flagPastelWhite, flagPastelWhite, flagPastelWhite,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue, flagPastelBlue,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelWhite, flagPastelWhite, flagPastelWhite, flagPastelWhite,
		flagPastelPink, flagPastelPink, flagPastelPink, flagPastelPink,
		flagPastelBlue, flagPastelBlue,
	}

	// Five-stripe lesbian flag (five colors are easier to distinguish than
	// seven).
	flagLesbianRed    = NewColor(0xC0, 0x00, 0x00)
	flagLesbianOrange = NewColor(0xC0, 0x40, 0x00)
	flagLesbianWhite  = NewColor(0x60, 0x50, 0x90)
	flagLesbianPink   = NewColor(0x40, 0x08, 0x30)
	flagLesbianPurple = NewColor(0x20, 0x00, 0x10)
	flagLesbian       = Palette{
		flagLesbianRed, flagLesbianRed, flagLesbianRed, flagLesbianRed,
		flagLesbianOrange, flagLesbianOrange, flagLesbianOrange, flagLesbianOrange,
		flagLesbianWhite, flagLesbianWhite, flagLesbianWhite, flagLesbianWhite,
		flagLesbianPink, flagLesbianPink, flagLesbianPink, flagLesbianPink,
		flagLesbianPurple, flagLesbianPurple, flagLesbianPurple, flagLesbianPurple, flagLesbianPurple, flagLesbianPurple,
		flagLesbianPink, flagLesbianPink, flagLesbianPink, flagLesbianPink,
		flagLesbianWhite, flagLesbianWhite, flagLesbianWhite, flagLesbianWhite,
		flagLesbianOrange, flagLesbianOrange, flagLesbianOrange, flagLesbianOrange,
		flagLesbianRed, flagLesbianRed,
	}

	// Flag for gay men. This is the less common 5 stripe version, since that's
	// easier to make on the earrings. (Not sure how many men will wear these,
	// but it's there for those who want it).
	flagGayDarkGreen  = NewColor(0x00, 0x18, 0x08)
	flagGayLightGreen = NewColor(0x10, 0x40, 0x20)
	flagGayWhite      = NewColor(0x60, 0x50, 0x90)
	flagGayLightBlue  = NewColor(0x10, 0x10, 0x90)
	flagGayDarkBlue   = NewColor(0x08, 0x00, 0x20)
	flagGay           = Palette{
		flagGayDarkGreen, flagGayDarkGreen, flagGayDarkGreen, flagGayDarkGreen,
		flagGayLightGreen, flagGayLightGreen, flagGayLightGreen, flagGayLightGreen,
		flagGayWhite, flagGayWhite, flagGayWhite, flagGayWhite,
		flagGayLightBlue, flagGayLightBlue, flagGayLightBlue, flagGayLightBlue,
		flagGayDarkBlue, flagGayDarkBlue, flagGayDarkBlue, flagGayDarkBlue, flagGayDarkBlue, flagGayDarkBlue,
		flagGayLightBlue, flagGayLightBlue, flagGayLightBlue, flagGayLightBlue,
		flagGayWhite, flagGayWhite, flagGayWhite, flagGayWhite,
		flagGayLightGreen, flagGayLightGreen, flagGayLightGreen, flagGayLightGreen,
		flagGayDarkGreen, flagGayDarkGreen,
	}

	// This one really pops! While the number of LEDs for each color is
	// balanced, it might look better with a bit more for yellow?
	flagNonBinaryYellow = NewColor(0xaa, 0xaa, 0x00)
	flagNonBinaryWhite  = NewColor(0x60, 0x50, 0x90)
	flagNonBinaryPurple = NewColor(0x50, 0x00, 0x40)
	flagNonBinaryBlack  = NewColor(0x00, 0x00, 0x00)
	flagNonBinary       = Palette{
		flagNonBinaryYellow, flagNonBinaryYellow, flagNonBinaryYellow, flagNonBinaryYellow,
		flagNonBinaryWhite, flagNonBinaryWhite, flagNonBinaryWhite, flagNonBinaryWhite, flagNonBinaryWhite,
		flagNonBinaryPurple, flagNonBinaryPurple, flagNonBinaryPurple, flagNonBinaryPurple, flagNonBinaryPurple,
		flagNonBinaryBlack, flagNonBinaryBlack, flagNonBinaryBlack, flagNonBinaryBlack, flagNonBinaryBlack, flagNonBinaryBlack, flagNonBinaryBlack, flagNonBinaryBlack,
		flagNonBinaryPurple, flagNonBinaryPurple, flagNonBinaryPurple, flagNonBinaryPurple, flagNonBinaryPurple,
		flagNonBinaryWhite, flagNonBinaryWhite, flagNonBinaryWhite, flagNonBinaryWhite, flagNonBinaryWhite,
		flagNonBinaryYellow, flagNonBinaryYellow, flagNonBinaryYellow, flagNonBinaryYellow,
	}

	flagBiPink   = NewColor(0xD0, 0x00, 0x08)
	flagBiPurple = NewColor(0x40, 0x00, 0x30)
	flagBiBlue   = NewColor(0x00, 0x00, 0x80)
	flagBi       = Palette{
		flagBiPink, flagBiPink, flagBiPink, flagBiPink, flagBiPink, flagBiPink, flagBiPink, flagBiPink,
		flagBiPurple, flagBiPurple, flagBiPurple, flagBiPurple,
		flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue, flagBiBlue,
		flagBiPurple, flagBiPurple, flagBiPurple, flagBiPurple,
		flagBiPink, flagBiPink, flagBiPink, flagBiPink, flagBiPink, flagBiPink,
	}

	flagPanPink   = NewColor(0xD0, 0x00, 0x10)
	flagPanYellow = NewColor(0xaa, 0xaa, 0x00)
	flagPanBlue   = NewColor(0x08, 0x08, 0xFF)
	flagPan       = Palette{
		flagPanPink, flagPanPink, flagPanPink, flagPanPink, flagPanPink, flagPanPink,
		flagPanYellow, flagPanYellow, flagPanYellow, flagPanYellow, flagPanYellow, flagPanYellow,
		flagPanBlue, flagPanBlue, flagPanBlue, flagPanBlue, flagPanBlue, flagPanBlue,
		flagPanBlue, flagPanBlue, flagPanBlue, flagPanBlue, flagPanBlue, flagPanBlue,
		flagPanYellow, flagPanYellow, flagPanYellow, flagPanYellow, flagPanYellow, flagPanYellow,
		flagPanPink, flagPanPink, flagPanPink, flagPanPink, flagPanPink, flagPanPink,
	}

	// Hard to represent because of the black, but I did my best. Since it's a
	// circle, the black does stand out.
	flagAceBlack  = NewColor(0x00, 0x00, 0x00)
	flagAceGray   = NewColor(0x10, 0x10, 0x18)
	flagAceWhite  = NewColor(0x60, 0x50, 0x90)
	flagAcePurple = NewColor(0x50, 0x00, 0x40)
	flagAce       = Palette{
		flagAceBlack, flagAceBlack, flagAceBlack, flagAceBlack,
		flagAceGray, flagAceGray, flagAceGray, flagAceGray, flagAceGray,
		flagAceWhite, flagAceWhite, flagAceWhite, flagAceWhite, flagAceWhite,
		flagAcePurple, flagAcePurple, flagAcePurple, flagAcePurple, flagAcePurple, flagAcePurple, flagAcePurple, flagAcePurple,
		flagAceWhite, flagAceWhite, flagAceWhite, flagAceWhite, flagAceWhite,
		flagAceGray, flagAceGray, flagAceGray, flagAceGray, flagAceGray,
		flagAceBlack, flagAceBlack, flagAceBlack, flagAceBlack,
	}

	// Aromantic flag.
	flagAroGreen1 = NewColor(0x00, 0x80, 0x08)
	flagAroGreen2 = NewColor(0x10, 0x60, 0x28)
	flagAroWhite  = NewColor(0x60, 0x50, 0x90)
	flagAroGray   = NewColor(0x10, 0x10, 0x18)
	flagAroBlack  = NewColor(0x00, 0x00, 0x00)
	flagAro       = Palette{
		flagAroGreen1, flagAroGreen1, flagAroGreen1, flagAroGreen1,
		flagAroGreen2, flagAroGreen2, flagAroGreen2, flagAroGreen2,
		flagAroWhite, flagAroWhite, flagAroWhite, flagAroWhite,
		flagAroGray, flagAroGray, flagAroGray, flagAroGray,
		flagAroBlack, flagAroBlack, flagAroBlack, flagAroBlack, flagAroBlack, flagAroBlack,
		flagAroGray, flagAroGray, flagAroGray, flagAroGray,
		flagAroWhite, flagAroWhite, flagAroWhite, flagAroWhite,
		flagAroGreen2, flagAroGreen2, flagAroGreen2, flagAroGreen2,
		flagAroGreen1, flagAroGreen1,
	}

	// The new polyamory flag, with the yellow heart:
	// https://nl.wikipedia.org/wiki/Bestand:Tricolor_Polyamory_Pride_Flag.svg
	// It's hard to represent the dark purple at the bottom, this seems as dark
	// and purplish as possible.
	flagPolyBlue   = NewColor(0x08, 0x08, 0xFF)
	flagPolyRed    = NewColor(0xC0, 0x00, 0x10)
	flagPolyPurple = NewColor(0x08, 0x00, 0x15)
	flagPolyYellow = NewColor(0xaa, 0xaa, 0x00)
	flagPolyWhite  = NewColor(0x60, 0x50, 0x90)
	flagPoly       = Palette{
		flagPolyBlue, flagPolyBlue, flagPolyBlue, flagPolyBlue, flagPolyBlue, flagPolyBlue,
		flagPolyRed, flagPolyRed, flagPolyRed, flagPolyRed, flagPolyRed, flagPolyRed,
		flagPolyPurple, flagPolyPurple, flagPolyPurple, flagPolyPurple, flagPolyPurple, flagPolyPurple,
		flagPolyPurple, flagPolyPurple, flagPolyPurple, flagPolyPurple, flagPolyPurple, flagPolyPurple,
		flagPolyRed, flagPolyRed, flagPolyRed, flagPolyRed,
		flagPolyWhite, flagPolyYellow, flagPolyYellow, flagPolyYellow, flagPolyWhite, // the heart shape
		flagPolyBlue, flagPolyBlue, flagPolyBlue,
	}

	allFlags = [...]Palette{
		flagLGBT,
		flagTrans,
		flagLesbian,
		flagGay,
		flagNonBinary,
		flagBi,
		flagPan,
		flagAce,
		flagAro,
		flagPoly,
	}
)

func showFlag(led, frame, variant int) Color {
	if variant >= len(allFlags) {
		// This shouldn't actually happen.
		return NewColor(0, 0, 0)
	}
	return allFlags[variant][led]
}

// Basic VU meter like animation, with different colors.
func vuMeter(led, frame, variant int) Color {
	// LED indices calculated, so they start from the bottom and alternate
	// between the two sides:
	// LED 14: index 6
	// LED 15: index 4
	// LED 16: index 2
	// LED 17: index 0
	// LED 18: index 1
	// LED 18: index 3
	// LED 18: index 5
	var ledIndex int
	if led < 18 {
		ledIndex = (17 - led) * 2
	} else {
		ledIndex = (led-18)*2 + 1
	}

	// Determine color for this VU meter.
	c := NewColor(0, 0, 0)
	if variant < len(staticColors) {
		c = staticColors[variant]
	}

	frameIntensity := int(currentVolume()) / 32
	intensity := frameIntensity - ledIndex*64
	if intensity > 255 {
		return c
	} else if intensity < 0 {
		return NewColor(0, 0, 0)
	} else {
		return c.Scale(intensity)
	}
}

// Sparkle animation!
// Credits to this person for suggesting it:
// https://hachyderm.io/@wmd@chaos.social/115816155535125138
// It's been vaguely inspired by the HALO-90 animation but this one looks very
// different.

var sparkleColors [numLEDs]Color
var sparkleBrightness [numLEDs]uint8

var sparkleVariants = [...][2]Color{
	[2]Color{NewColor(0xff, 0x00, 0xff), NewColor(0x00, 0x99, 0x66)}, // purple, turquoise
	[2]Color{NewColor(0xff, 0x00, 0xff), NewColor(0xff, 0x00, 0x00)}, // purple, red
	[2]Color{NewColor(0xff, 0x00, 0x00), NewColor(0xff, 0xff, 0x00)}, // red, yellow
	[2]Color{NewColor(0xff, 0xff, 0x00), NewColor(0x00, 0xff, 0x00)}, // yellow, green
}

func sparkleNextFrame(variant int) {
	addPower(uint16(processSamples()))

	vol := int(currentVolume()) - 4096

	r := rand()
	if uint32(vol) > r&0x7fff {
		// Create a new sparkle!

		// Pick an index at random.
		index := fastModulo36(vol & 0xff)

		// Pick one of the two colors.
		if variant > len(sparkleVariants) {
			variant = 0 // shouldn't happen, but just to be sure
		}
		c := sparkleVariants[variant][r>>31]

		// Store this sparkle.
		sparkleColors[index] = c
		sparkleBrightness[index] = 0xff
	}
}

func makeSparkle(led, frame, variant int) Color {
	dimming := sparkleBrightness[led]
	if dimming < 32 {
		return NewColor(0, 0, 0)
	}
	sparkleBrightness[led] = dimming - 32
	return sparkleColors[led].Scale(int(dimming))
}

// Blink the first LED, roughly 0.5s on, 0.5s off.
func errorPattern(led, frame int) Color {
	if led == 0 {
		// Roughly 500ms on, 500ms off (assuming 30fps).
		return NewColor(uint8(frame%32)/16*128, 0, 0)
	}
	// Other LEDs are dark.
	return NewColor(0, 0, 0)
}

func rotateSingleColor(led, frame int) Color {
	idx := int(frame / 8 % 36)
	value := uint8(0)
	if led == idx {
		value = 128
	}
	if (led+1)%36 == idx {
		value = 128 / 2
	}
	if (led+2)%36 == idx {
		value = 128 / 4
	}
	if (led+3)%36 == idx {
		value = 128 / 8
	}
	if (led+4)%36 == idx {
		value = 128 / 16
	}
	if (led+5)%36 == idx {
		value = 128 / 32
	}
	return NewColor(0, 0, value)
}

// Pulse red LEDs around once per second, for testing.
func testPulse(led, frame int) Color {
	return NewColor(uint8((frame%32)<<3), 0, 0)
}

// Startup animation. The duration until it will power on is set in
// waitForPowerOn, but this gives a nice animation to show how far it is.
func powerOn(led, frame int) Color {
	if frame*2 > led {
		return NewColor(0, 0x3f, 0) // green
	}
	return NewColor(0, 0, 0)
}
