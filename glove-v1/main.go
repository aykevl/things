package main

const (
	modeTest = iota
	modePurple
	modeTrail
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
	case modeTrail:
		c = trail(cycles)
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

var trailColors = [...]Color{
	0x550055,
	0x550055,
	0x550055,
	0x550055,
	0x550055,
	0x550055,
	0x550055,
	0x550055,
	0,
	0,
	0,
	0,
	0x005500,
	0x005500,
	0x005500,
	0x005500,
	0,
	0,
	0,
	0,
}
var trailIndex = uint8(255)

func trail(cycles uint32) Color {
	if cycles%2 == 0 {
		trailIndex++
	}
	if int(trailIndex) >= len(trailColors[:]) {
		trailIndex = 0
	}

	return trailColors[trailIndex]
}
