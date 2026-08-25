package main

func animation(cycles uint32) {
	// Show animation
	switch (cycles >> 10) % 4 {
	case 0:
		setColor(0xff0000)
	case 1:
		setColor(0x00ff00)
	case 2:
		setColor(0x0000ff)
	case 3:
		setColor(0x555555)
	}
}
