#include <stdint.h>

static void bitbang_show_leds(uint8_t c1, uint8_t c2, uint8_t c3, uint8_t c4, uint8_t c5, uint8_t c6, volatile uint8_t *port) {
    uint8_t savedOut = *port;
    uint8_t tmp1 = savedOut;
    uint8_t tmp2 = savedOut;
    uint8_t tmp3 = savedOut;
    __asm__ __volatile__(
        // precalculate 32-cycle bitplane
        "bst %[c1], 7\n\t"
        "bld %[tmp1], 6\n\t"
        "bst %[c2], 7\n\t"
        "bld %[tmp1], 5\n\t"
        "bst %[c3], 7\n\t"
        "bld %[tmp1], 4\n\t"
        "bst %[c4], 7\n\t"
        "bld %[tmp1], 3\n\t"
        "bst %[c5], 7\n\t"
        "bld %[tmp1], 2\n\t"
        "bst %[c6], 7\n\t"
        "bld %[tmp1], 1\n\t"
        // precalculate 16-cycle bitplane
        "bst %[c1], 6\n\t"
        "bld %[tmp2], 6\n\t"
        "bst %[c2], 6\n\t"
        "out %[port], %[tmp1]\n\t" // start 32 cycles
        "bld %[tmp2], 5\n\t"
        "bst %[c3], 6\n\t"
        "bld %[tmp2], 4\n\t"
        "bst %[c4], 6\n\t"
        "bld %[tmp2], 3\n\t"
        "bst %[c5], 6\n\t"
        "bld %[tmp2], 2\n\t"
        "bst %[c6], 6\n\t"
        "bld %[tmp2], 1\n\t"
        // precalculate 8-cycle bitplane
        "bst %[c1], 5\n\t"
        "bld %[tmp1], 6\n\t"
        "bst %[c2], 5\n\t"
        "bld %[tmp1], 5\n\t"
        "bst %[c3], 5\n\t"
        "bld %[tmp1], 4\n\t"
        "bst %[c4], 5\n\t"
        "bld %[tmp1], 3\n\t"
        "bst %[c5], 5\n\t"
        "bld %[tmp1], 2\n\t"
        "bst %[c6], 5\n\t"
        "bld %[tmp1], 1\n\t"
        // precalculate 4-cycle bitplane
        "bst %[c1], 4\n\t"
        "bld %[tmp3], 6\n\t"
        "bst %[c2], 4\n\t"
        "bld %[tmp3], 5\n\t"
        "bst %[c3], 4\n\t"
        "bld %[tmp3], 4\n\t"
        "bst %[c4], 4\n\t"
        "bld %[tmp3], 3\n\t"
        "bst %[c5], 4\n\t"
        "bld %[tmp3], 2\n\t"
        "out %[port], %[tmp2]\n\t" // end 32 cycles, start 16 cycles
        "bst %[c6], 4\n\t"
        "bld %[tmp3], 1\n\t"
        // precalculate 2-cycle bitplane
        "bst %[c1], 3\n\t"
        "bld %[tmp2], 6\n\t"
        "bst %[c2], 3\n\t"
        "bld %[tmp2], 5\n\t"
        "bst %[c3], 3\n\t"
        "bld %[tmp2], 4\n\t"
        "bst %[c4], 3\n\t"
        "bld %[tmp2], 3\n\t"
        "bst %[c5], 3\n\t"
        "bld %[tmp2], 2\n\t"
        "bst %[c6], 3\n\t"
        "bld %[tmp2], 1\n\t"
        // precalculate 1-cycle bitplane
        "bst %[c1], 2\n\t"
        "out %[port], %[tmp1]\n\t" // end 16 cycles, start 8 cycles
        "bld %[tmp1], 6\n\t"
        "bst %[c2], 2\n\t"
        "bld %[tmp1], 5\n\t"
        "bst %[c3], 2\n\t"
        "bld %[tmp1], 4\n\t"
        "bst %[c4], 2\n\t"
        "bld %[tmp1], 3\n\t"
        "out %[port], %[tmp3]\n\t" // end 8 cycles, start 4 cycles
        "bst %[c5], 2\n\t"
        "bld %[tmp1], 2\n\t"
        "bst %[c6], 2\n\t"
        "out %[port], %[tmp2]\n\t" // end 4 cycles, start 2 cycles
        "bld %[tmp1], 1\n\t"
        "out %[port], %[tmp1]\n\t" // end 2 cycles, start 1 cycle
        "out %[port], %[savedOut]\n\t" // end 1 cycle, restore old port state
        : [tmp1]"+r"(tmp1),
          [tmp2]"+r"(tmp2),
          [tmp3]"+r"(tmp3)
        : [c1]"r"(c1),
          [c2]"r"(c2),
          [c3]"r"(c3),
          [c4]"r"(c4),
          [c5]"r"(c5),
          [c6]"r"(c6),
          [port]"I"(port),
          [savedOut]"r"(savedOut)
    );
    *port = savedOut;
}
