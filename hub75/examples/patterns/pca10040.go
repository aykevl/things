// +build pca10040

package main

import (
	"github.com/aykevl/things/hub75"
)

var display = hub75.New(hub75.Config{
	Data:         22,
	Clock:        23,
	Latch:        24,
	OutputEnable: 25,
	A:            20,
	B:            19,
	C:            18,
	D:            17,
})
