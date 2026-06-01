package main

type Color uint32

func NewColor(r, g, b uint8) Color {
	return Color(r)<<0 | Color(g)<<8 | Color(b)<<16
}

func (c Color) R() uint8 {
	return uint8(c >> 0)
}

func (c Color) G() uint8 {
	return uint8(c >> 8)
}

func (c Color) B() uint8 {
	return uint8(c >> 16)
}

// Linearly scale the given color by the given intensity, which must be 0..256.
func (c Color) Scale(intensity int) Color {
	return NewColor(
		uint8((int(c.R())*intensity)/256),
		uint8((int(c.G())*intensity)/256),
		uint8((int(c.B())*intensity)/256))
}
