# LED poi

This is some code for a [poi](https://en.wikipedia.org/wiki/Poi_(performance_art)). You can see a video in [this tweet](https://twitter.com/aykevl/status/1280135326913232896). [This](https://www.youtube.com/watch?v=LGy0neDXxAE0) is a professionally made video of a similar device.

As of this writing, I have built two of these: one with a GT82C_02 board inside and one with a GT832E_01 board (both Nordic chips). These boards are easy to solder, are very small and only require a single external component (a voltage regulator) to run them. It is very loosely based on [this build](http://orchardelica.com/wp/?page_id=597).

You can control the poi from [here](https://aykevl.nl/apps/poi/), for example. It's just Bluetooth Low Energy, and it's also possible to control it using nRF Connect ([Android version](https://play.google.com/store/apps/details?id=no.nordicsemi.android.mcp)).
