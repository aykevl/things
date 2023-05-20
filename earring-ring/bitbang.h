#include <stdint.h>

static void bitbang_show_leds(uint8_t c1, uint8_t c2, uint8_t c3, uint8_t c4, uint8_t c5, uint8_t c6, volatile uint8_t *toggleBits) {
    uint8_t savedOut = *toggleBits;
    uint8_t i;
    __asm__ __volatile__(
        "mov __tmp_reg__, %[savedOut]\n\t"

        // 32 cycles
        "bst %[c1], 7\n\t"
        "bld __tmp_reg__, 6\n\t"
        "bst %[c2], 7\n\t"
        "bld __tmp_reg__, 5\n\t"
        "bst %[c3], 7\n\t"
        "bld __tmp_reg__, 4\n\t"
        "bst %[c4], 7\n\t"
        "bld __tmp_reg__, 3\n\t"
        "bst %[c5], 7\n\t"
        "bld __tmp_reg__, 2\n\t"
        "bst %[c6], 7\n\t"
        "bld __tmp_reg__, 1\n\t"
        "out %[port], __tmp_reg__\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        // precalculate 16-cycle bitplane
        "bst %[c1], 6\n\t"
        "bld __tmp_reg__, 6\n\t"
        "bst %[c2], 6\n\t"
        "bld __tmp_reg__, 5\n\t"
        "bst %[c3], 6\n\t"
        "bld __tmp_reg__, 4\n\t"
        "bst %[c4], 6\n\t"
        "bld __tmp_reg__, 3\n\t"
        "bst %[c5], 6\n\t"
        "bld __tmp_reg__, 2\n\t"
        "bst %[c6], 6\n\t"
        "bld __tmp_reg__, 1\n\t"

        // End 32 cycles, start 16 cycles.
        "out %[port], __tmp_reg__\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        // precalculate 8-cycle bitplane
        "bst %[c1], 5\n\t"
        "bld __tmp_reg__, 6\n\t"
        "bst %[c2], 5\n\t"
        "bld __tmp_reg__, 5\n\t"
        "bst %[c3], 5\n\t"
        "bld __tmp_reg__, 4\n\t"
        "bst %[c4], 5\n\t"
        "bld __tmp_reg__, 3\n\t"
        "bst %[c5], 5\n\t"
        "bld __tmp_reg__, 2\n\t"
        "bst %[c6], 5\n\t"
        "bld __tmp_reg__, 1\n\t"

        // End 16 cycles, start 8 cycles.
        "out %[port], __tmp_reg__\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "out %[port], %[savedOut]\n\t"

        // 4 cycles
        "bst %[c1], 4\n\t"
        "bld __tmp_reg__, 6\n\t"
        "bst %[c2], 4\n\t"
        "bld __tmp_reg__, 5\n\t"
        "bst %[c3], 4\n\t"
        "bld __tmp_reg__, 4\n\t"
        "bst %[c4], 4\n\t"
        "bld __tmp_reg__, 3\n\t"
        "bst %[c5], 4\n\t"
        "bld __tmp_reg__, 2\n\t"
        "bst %[c6], 4\n\t"
        "bld __tmp_reg__, 1\n\t"
        "out %[port], __tmp_reg__\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "out %[port], %[savedOut]\n\t"

        // 2 cycles
        "bst %[c1], 3\n\t"
        "bld __tmp_reg__, 6\n\t"
        "bst %[c2], 3\n\t"
        "bld __tmp_reg__, 5\n\t"
        "bst %[c3], 3\n\t"
        "bld __tmp_reg__, 4\n\t"
        "bst %[c4], 3\n\t"
        "bld __tmp_reg__, 3\n\t"
        "bst %[c5], 3\n\t"
        "bld __tmp_reg__, 2\n\t"
        "bst %[c6], 3\n\t"
        "bld __tmp_reg__, 1\n\t"
        "out %[port], __tmp_reg__\n\t"
        "nop\n\t"
        "out %[port], %[savedOut]\n\t"

        // 1 cycle
        "bst %[c1], 2\n\t"
        "bld __tmp_reg__, 6\n\t"
        "bst %[c2], 2\n\t"
        "bld __tmp_reg__, 5\n\t"
        "bst %[c3], 2\n\t"
        "bld __tmp_reg__, 4\n\t"
        "bst %[c4], 2\n\t"
        "bld __tmp_reg__, 3\n\t"
        "bst %[c5], 2\n\t"
        "bld __tmp_reg__, 2\n\t"
        "bst %[c6], 2\n\t"
        "bld __tmp_reg__, 1\n\t"
        "out %[port], __tmp_reg__\n\t"
        "out %[port], %[savedOut]\n\t"
        : [c1]"+r"(c1),
          [c2]"+r"(c2),
          [c3]"+r"(c3),
          [c4]"+r"(c4),
          [c5]"+r"(c5),
          [c6]"+r"(c6)
        : [port]"I"(toggleBits),
          [savedOut]"r"(savedOut)
    );
    *toggleBits = savedOut;
}
