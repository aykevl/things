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

// Define modes that can be cycled through.
enum {
  mode_sparkling, // also the power on mode
  mode_wave,
  mode_off,
  num_modes,
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

  // Reduce clock speed to 250kHz.
  CCP = 0xD8;      // unlock protected registers
  CLKPSR = 0b0101; // division factor 32 (8 / 32 = 0.25)

  if ((RSTFLR & EXTRF) == 0) {
    // Not an external reset, so probably a power on reset.
    // Initialize global variables.
    mode = 0;
    rnd_s = 0xaa;
    rnd_a = 0;
  } else {
    // Reset pin triggered, so go to the next mode.
    mode++;
    if (mode >= num_modes) {
      mode = 0;
    }
  }

  if (mode == mode_off) {
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
    if (mode == mode_wave) {
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
          b -= 2;
        }
        if (b & 64) {
          b -= 1;
        }
        if ((b & 32) && ((cycle % 2) == 0)) {
          b--;
        }
        if (b > (256 / MAX_BRIGHTNESS)) {
          if ((cycle % 4) == 0) {
            b--;
          }
        }
        brightness[i] = b;
      }
      if (cycle % 8 == 0) {
        uint8_t r = rnd();
        if ((r & 0b11000000) == 0b11000000) {
          // Pick one LED at random and turn it on.
          uint8_t index = r % 6;
          if (brightness[index] < r) {
            brightness[index] = r;
          }
        }
      }
    }

    // Update LEDs using charlieplexing.
    for (uint8_t delay = 0; delay < (768 / MAX_BRIGHTNESS); delay++) {
      // Turn each LED on for just the right amount of time.
      for (uint8_t i=0; i<6; i++) {
        uint8_t state = states[i];

        // Configure port to only turn a particular LED on.
        PORTB = state >> 4;

        // The number of clock cycles the output should be active for.
        uint8_t numCycles = brightness[i];

        // Use custom assembly to turn the LEDs on for exactly the given number
        // of clock cycles, while keeping the code constant-time.
        // By lowering the resolution to 1 cycle, we can reduce the chip speed a
        // lot, which is a significant reduction in power consumption.
        __asm__ volatile(
          // 1 cycle: output bit 0 (least-significant bit)
          "lsr %[num]\n\t"
          "brcc 1f\n\t"
          "out %[port], %[state]\n\t"
          "1:\n\t"
          "out %[port], __zero_reg__\n\t"

          // 2 cycles: output bit 1
          "lsr %[num]\n\t"
          "brcc 1f\n\t"
          "out %[port], %[state]\n\t"
          "1:\n\t"
          "nop\n\t"
          "out %[port], __zero_reg__\n\t"

          // 4 cycles: output bit 2
          "lsr %[num]\n\t"
          "brcc 1f\n\t"
          "out %[port], %[state]\n\t"
          "1:\n\t"
          "nop\n\t"
          "nop\n\t"
          "nop\n\t"
          "out %[port], __zero_reg__\n\t"

          // 8 cycle loop: output bit 3-7
          "cpse %[num], __zero_reg__\n\t"      // if num != 0:
          "out %[port], %[state]\n\t"          //   LEDs on
          "ldi __tmp_reg__, %[iterations]\n\t" // i = 32
          "1:\n\t"                             // loop:
          "nop\n\t"                            //   wait 2 cycles
          "nop\n\t"                            //
          "subi %[num], 1\n\t"                 //   num--
          "brpl 2f\n\t"                        //   if num == -1:
          "out %[port], __zero_reg__\n\t"      //     LEDs off
          "2:\n\t"                             //
          "subi __tmp_reg__, 1\n\t"            //   i--
          "brne 1b\n\t"                        //   if i == 0: goto loop
          : [num]"+r"(numCycles)
          : [port]"I"(_SFR_IO_ADDR(DDRB)),
            [state]"r"(state),
            [iterations]"I"(MAX_BRIGHTNESS / 8)
        );
      }
    }
  }

  // unreachable
}
