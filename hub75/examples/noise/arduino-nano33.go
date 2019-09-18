// +build arduino_nano33

package main

import (
	"machine"

	"github.com/aykevl/things/hub75"
)

// data, clock, latch, oe, abcd
var display = hub75.New(machine.D2, machine.D3, machine.D4, machine.D5, machine.D8, machine.D9, machine.D10, machine.D11)
