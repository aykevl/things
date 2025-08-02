package main

type Color uint32

func NewColor(r, g, b uint8) Color {
	return Color(r)<<0 | Color(g)<<8 | Color(b)<<16
}
