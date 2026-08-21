//go:build pico
// +build pico

package main

import "machine"

const (
	LED_PIN       = machine.GPIO29
	LED_POWER_PIN = machine.GP14
)
