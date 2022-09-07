# If flashing fails, put a 330-660Ω resistor between MOSI and VCC.


all: main.hex
clean:
	rm -f main.elf main.hex

.PHONY: flash
flash: main.hex
	avrdude -p attiny10 -c usbasp -U flash:w:main.hex:i

main.hex: main.elf
	avr-objcopy -O ihex main.elf main.hex

main.elf: main.c
	avr-gcc -o main.elf main.c -mmcu=attiny10 -Os -Wl,--gc-sections -flto
	avr-size main.elf
