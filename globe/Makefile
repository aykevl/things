
all: flash

flash: globe.hex
	openocd -f interface/stlink-v2.cfg -f target/nrf51.cfg -c 'program globe.hex reset exit'

.PHONY: globe.hex # always rebuild
globe.hex:
	tinygo build -o globe.hex -target microbit -size=short .
