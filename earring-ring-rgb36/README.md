# 36-LED RGB earrings

These earrings are the next version of my [RGB earrings](../earring-ring). They contain 36 RGB LEDs that are controlled with a variant on Charlieplexing.

## Improvements

Changes compared to the previous 18 LED earrings:

  * They have better LEDs: these LEDs are 5-7x brighter at the same current consumption. This means they can run at lower power for the same brightness. The LEDs also include a little bit of diffusion which makes them look a lot nicer.
  * A 32-bit ARM chip means it can run more advanced animations at the same clock speed.
  * The LED bit depth is a bit lower: changed from 6 bits to 5 bits. But the colors should be more accurate due to added resistors which divide the current better.
  * 3 RGB LEDs get updated at the same time, instead of going through each color separately. Some people might find this less annoying (colors don't split up when moving your eyes and seeing the colors through persistence-of-vision).
  * The chip runs at a _lower_ clock speed to save power even with double the LEDs: 262kHz (it was 625kHz on the previous 18 LED earring). It can do so because of the reduced bit depth and because it updates more LEDs at the same time.
  * The chip uses a lot less power: it uses around 55µA while running while the 18 LED earring used around 410µA (reduction of ~87%, or around 7-8x lower power consumption). This is because of the lower clock frequency, but also because the STM32L0 series chips have a multispeed clock that doesn't need to run at a higher frequency than needed (the AVR chip on the 18 LED earrings run off of a divided 20MHz clock).
  * Reverse power protection was added, so now it won't harm the LEDs when the battery is inserted backwards. (With the previous LEDs, it didn't seem to harm them either but it was technically out of spec).

Overall the result is that the earrings look nicer, use less power, and are a little brighter.

## Programming

For programming, you will need an SWD programmer. Both ST-Link and DAPLink programmers have been tested and work fine. You can connect the wires as indicated on the PCB itself. Specifically:

  * Connect VCC to 3.3V, or don't connect it if you have a battery inserted (make sure to _never_ connect VCC when a battery is inserted!).
  * Connect GND to the programmer GND.
  * Connect SWC to the SWCLK pin on the programmer.
  * Connect SWD to the SWDIO pin on the programmer.

You normally don't need to connect RST, but it is exposed since it's possible to put the chip in a state where it doesn't respond to programming anymore an you need to briefly connect RST to GND to reset it.

I found the easiest way to program these earrings is using a clip with pogo-pins, like [this one](https://nl.aliexpress.com/item/1005006712952020.html) (2.54mm distance, single row, 4 or 5 pins depending on needs). It makes quick iteration much easier.

## Credits

I got a lot of inspiration from the earrings made by [California STEAM](https://www.tindie.com/stores/californiasteam/). Unfortunately they're not shipping to Europe so I had to make my own 🙂

Another cool project I found was [this earring on Hackaday](https://hackaday.io/project/186402-ws2812b-neopixel-earring). It's built quite differently however, even using WS2812 LEDs!

And after I made my [previous earrings](../earring-ring), I found these [MNT x Kolibri HALO-90 electronic earrings](https://shop.mntre.com/products/mnt-x-kolibri-halo-90-electronic-earrings) which inspired me to get even higher pixel counts. 18 LEDs just doesn't cut it anymore in 2025! They've definitely been an inspiration, though of course their choices in electronics and animation patterns are quite different.
