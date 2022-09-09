#include <avr/io.h>
#include <avr/sleep.h>

// Control the number of light level bits to use for the LEDs.
// A higher number means a better resolution (noticeable at low light levels)
// but also increases the amount of flickering.
#define MAX_BRIGHTNESS 128

// For each LED, the PORTB and DDRB registers.
// DDRB is the lower 4 bits while PORTB is the higher 4 bits.
static const uint8_t states[6] = {
  0b00100110, // 5
  0b00010101, // 3
  0b00010011, // 1
  0b01000110, // 6
  0b01000101, // 4
  0b00100011, // 2
};

static const uint8_t sinewave[64] = {
  0, 0, 2, 4, 7, 11, 17, 23, 31, 39, 49, 59, 70, 82, 94, 106, 119, 132, 145, 157, 170, 182, 193, 204, 214, 223, 231, 238, 244, 249, 252, 254, 255, 254, 252, 249, 244, 238, 231, 223, 214, 204, 193, 182, 170, 157, 145, 132, 119, 106, 94, 82, 70, 59, 49, 39, 31, 23, 17, 11, 7, 4, 2, 0
};

// Source: https://www.avrfreaks.net/forum/tiny-fast-prng
__attribute__((section(".noinit")))
static uint8_t rnd_s, rnd_a;
static uint8_t rnd(void) {
  rnd_s ^= rnd_s<<3;
  rnd_s ^= rnd_s>>5;
  rnd_s ^= rnd_a++>>2;
  return rnd_s;
}

__attribute__((section(".noinit")))
static uint8_t brightness[6];

__attribute__((section(".noinit")))
static uint8_t mode;

int main(void) {
  // Reduce clock speed to 128kHz.
  //CCP = 0xD8;
  //CLKMSR = 0b01;
  //CCP = 0xD8;
  //CLKPSR = 0;

  if ((RSTFLR & EXTRF) == 0) {
    // Not an external reset, so probably a power on reset.
    // Initialize global variables.
    mode = 0;
    rnd_s = 0xaa;
    rnd_a = 0;
  } else {
    // Reset pin triggered, so go to the next mode.
    mode++;
    if (mode > 2) {
      mode = 0;
    }
  }

  if (mode == 2) {
    // Power down the chip entirely. It will only awake when the reset pin
    // is pressed.
    // This reduces power consumption to around 1µA, or as low as my
    // multimeter will measure.
    SMCR = 0b0101; // SM = 0b010 (power down), SE = 1 (sleep enabled)
    while (1) {
      sleep_cpu();
    }
  }

  uint8_t cycle = 0;
  while (1) {
    cycle++;

    // Run next step in the animation.
    if (mode == 1) {
      // Do a sine wave in a circle.
      for (int8_t i=5; i >= 0; i--) {
        uint8_t index = cycle - (uint8_t)(i * 42);
        brightness[i] = 0;
        if (index < 128) {
          brightness[i] = sinewave[index / 2];
        }
      }
    } else {
      // Turn LEDs randomly on, and let them fade out.
      for (int8_t i=5; i >= 0; i--) {
        // Reduce brightness. Quicly at first, but slower at lower
        // brightnesses.
        // This is supposed to look like a power law, but I'm not sure it
        // does.
        uint8_t b = brightness[i];
        if (b & 128) {
          b -= 4;
        }
        if (b & 64) {
          b -= 2;
        }
        if (b & 32) {
          b--;
        }
        if (b > (256 / MAX_BRIGHTNESS)) {
          if ((cycle % 4) == 0) {
            b--;
          }
        }
        brightness[i] = b;
      }
      uint8_t r = rnd();
      if ((r & 0b11110000) == 0b11110000) {
        // Pick one LED at random and turn it on.
        uint8_t index = r % 6;
        if (brightness[index] < r) {
          brightness[index] = r;
        }
      }
    }

    // Update LEDs using charlieplexing.
    for (uint8_t delay = 0; delay < (1024 / MAX_BRIGHTNESS); delay++) {
      // Turn each LED on for just the right amount of time.
      for (uint8_t i=0; i<6; i++) {
        uint8_t state = states[i];
        // TODO: use dithering to increase the perceived resolution at a higher
        // brightness while reducing the amount of flickering?
        uint8_t b = brightness[i] / (256 / MAX_BRIGHTNESS);
        for (uint8_t bit = MAX_BRIGHTNESS / 2; bit != 0; bit >>= 1) {
          PORTB = state >> 4;
          if ((b & bit) != 0) {
            // Note: this should be (state & 0x0f) but the hardware
            // ignores the upper 4 bits and avoiding the mask avoids two
            // instructions.
            DDRB = state;
          }
          for (uint8_t j=bit; j != 0; j--) {
            __asm__ volatile("");
          }
          DDRB = 0;
        }
      }
    }
  }

  // unreachable
}
