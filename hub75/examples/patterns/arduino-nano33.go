// +build arduino_nano33

package main

import (
	"machine"

	"github.com/aykevl/things/hub75"
)

var display = hub75.New(hub75.Config{
	Data:         machine.NoPin,
	Clock:        machine.NoPin,
	Latch:        machine.D5,
	OutputEnable: machine.D4,
	A:            machine.D8,
	B:            machine.D9,
	C:            machine.D10,
	D:            machine.D11,
})
