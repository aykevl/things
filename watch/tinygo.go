//go:build tinygo

package main

import (
	"runtime"
	"time"
)

func adjustTime(offset time.Duration) {
	runtime.AdjustTimeOffset(int64(offset))
}

func watchTime() time.Time {
	return time.Now()
}
