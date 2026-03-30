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
