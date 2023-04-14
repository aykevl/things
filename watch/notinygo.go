//go:build !tinygo

package main

import "time"

var timeOffset time.Duration

func adjustTime(offset time.Duration) {
	timeOffset += offset
}

func watchTime() time.Time {
	return time.Now().Add(timeOffset)
}
