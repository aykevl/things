// +build itsybitsy_m4

package main

import (
	"machine"

	"github.com/aykevl/things/hub75"
)

var display = hub75.New(hub75.Config{
	Data:         machine.NoPin,
	Clock:        machine.NoPin,
	Latch:        machine.D5,
	OutputEnable: machine.D7,
	A:            machine.D9,
	B:            machine.D10,
	C:            machine.D11,
	D:            machine.D12,
})
