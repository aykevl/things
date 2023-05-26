//go:build pybadge || arduino_nano33

package main

import "machine"

const (
	LED_PIN = machine.D13
)
