package main

const (
	modeTest = iota
	modePurple
	modeTrail0
	modeTrail1
	modeLast

	firstMode = modeTest
)

func animation(cycles uint32, mode int) {
	var c Color
	switch mode {
	case modeTest:
		c = testPattern(cycles)
	case modePurple:
		c = 0x550055
	case modeTrail0:
		c = showTrail(cycles, 0)
	case modeTrail1:
		c = showTrail(cycles, 1)
	}
	setColor(c)
}

func testPattern(cycles uint32) Color {
	// Show animation
	switch (cycles >> 10) % 4 {
	case 0:
		return 0xff0000
	case 1:
		return 0x00ff00
	case 2:
		return 0x0000ff
	default: // 3
		return 0x555555
	}
}

var trails = [...][]Color{
	{ // green-purple
		0x550055,
		0x550055,
		0x550055,
		0x550055,
		0,
		0,
		0x005500,
		0x005500,
		0,
		0,
	},
	{ // red-purple
		0x550000,
		0x660000,
		0x770000,
		0x660000,
		0x550000,
		0,
		0,
		0x330033,
		0x440044,
		0x550055,
		0x440044,
		0x330033,
		0,
		0,
	},
}

var trailIndex = uint8(255)

func showTrail(cycles uint32, which int) Color {
	trail := trails[which]
	if cycles%4 == 0 {
		trailIndex++
	}
	if int(trailIndex) >= len(trail) {
		trailIndex = 0
	}

	return trail[trailIndex]
}
