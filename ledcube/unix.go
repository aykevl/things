// +build !baremetal

package main

import (
	"log"

	"github.com/aykevl/tilegraphics/sdlscreen"
)

var display *sdlscreen.Screen

func init() {
	var err error
	const scale = 5
	display, err = sdlscreen.NewScreen("LED cube", 32*6*scale, 32*scale)
	if err != nil {
		log.Fatalln("could not instantiate screen:", err)
	}
	display.Scale = scale
}

func getFullRefreshes() uint {
	return 0
}
