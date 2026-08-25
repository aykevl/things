package main

const (
	modeTest = iota
	modePurple
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
