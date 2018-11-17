
all: flash

.PHONY: globe.hex globe.elf # always rebuild

GOFLAGS=-target microbit -size=short

flash: globe.hex
	openocd -f interface/stlink-v2.cfg -f target/nrf51.cfg -c 'program globe.hex reset exit'

globe.hex:
	tinygo build -o globe.hex $(GOFLAGS) .

globe.elf:
	tinygo build -o globe.hex $(GOFLAGS) .
