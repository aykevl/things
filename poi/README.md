# LED poi

This is some code for a [poi](https://en.wikipedia.org/wiki/Poi_(performance_art)). You can see a video in [this tweet](https://twitter.com/aykevl/status/1280135326913232896). [This](https://www.youtube.com/watch?v=LGy0neDXxAE0) is a professionally made video of a commercially available variant.

As of this writing, I've built three types of poi, all with Nordic chips inside (loosely based on [this build](http://orchardelica.com/wp/?page_id=597)):

  * The first one with a GT82C_02 board inside with a nRF51822 chip.
  * The second one with a GT832E_01 board inside with a nRF52832 chip.
  * The third one with a GT832C_01 board inside. It is similar to the GT832E_01 except that it also contains a BMI160 accelerometer/gyroscope and is slightly smaller.

I quite like these boards. They're easy to solder with 2.54mm or 2.00mm distance between outside connections, are very small, and only require an external voltage regulator to use them on a lithium ion battery.

The electronics of the third type (which I might stick with, it works well) are the following:

  * [GT832C_01](https://www.aliexpress.com/item/4000123520442.html) module (nRF52832 with BMI160).
  * Microchip MPC1702 3.3V regulator, see note 1 below.
  * Fairchild FQP30N06L MOSFET to turn the LED strips on and off. Bought from AliExpress so no idea whether it was genuine.
  * 1m 144 LED SK9822 LED strip, see note 2 below. Each side uses a quarter (25cm or 36 LEDs) of the strip. Both sides are connected in parallel.
  * Protected 14500 Li-ion battery: [this one](https://eu.nkon.nl/keeppower-16588.html). Protected so that an accidental short won't blow up the battery.
  * AA battery holder. Be warned that not all AA battery holders can fit protected Li-ion batteries because they are slightly longer than regular AA batteries. Battery holders that are 57mm long on the outside seem to be long enough usually.

Note 1: you can in principle use any 3.3V voltage regulator, but to be usable with a lithium battery (li-ion or li-poly), it has to be one with a low drop out rate. I tried an L78L33 from STMicroelectronics but the dropout voltage was too big to be usable: the LED strips started misbehaving when the battery was still at 3.85V (nearly full). Be [especially careful](https://hackaday.com/2018/11/10/what-good-are-counterfeit-parts-believe-it-or-not-maybe-a-refund/) with 3.3V voltage regulators from AliExpress. They might work well to drop from 5V to 3.3V, but may have a too large dropout rate for a clean signal between the MCU and the LED strip.

Note 2: use [SK9822 LED strips](https://cpldcpu.wordpress.com/2016/12/13/sk9822-a-clone-of-the-apa102/). Not APA102 because it [does not support flicker-free brightness control](https://cpldcpu.wordpress.com/2014/08/27/apa102/) and not WS2812 because it is a different protocol and it is not suitable for persistence of vision applications (it will cause tons of flickering). SK9822 is a near-clone of APA102 but the brightness control improvement is actually very important for this electronics project.

Note 3: the GT832C_01 board seems to be memory protected from the factory. Perhaps it has some firmware installed on it by default, I haven't checked whether it does anything. See [this post](https://devzone.nordicsemi.com/f/nordic-q-a/17015/how-do-i-disable-control-access-port-protection-on-nrf52-using-openocd) for how to remove the protection to make it programmable.

You can control the poi from [here](https://aykevl.nl/apps/poi/), for example. It's just Bluetooth Low Energy, and it's also possible to control it using nRF Connect ([Android version](https://play.google.com/store/apps/details?id=no.nordicsemi.android.mcp)).
