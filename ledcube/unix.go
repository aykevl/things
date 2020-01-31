// +build !baremetal

package main

import (
	"log"

	"github.com/aykevl/tilegraphics/sdlscreen"
)

var display *sdlscreen.Screen

func init() {
	var err error
	const scale = 192 / size
	display, err = sdlscreen.NewScreen("LED cube", size*6*scale, size*scale)
	if err != nil {
		log.Fatalln("could not instantiate screen:", err)
	}
	display.Scale = scale
}

func getFullRefreshes() uint {
	return 0
}
