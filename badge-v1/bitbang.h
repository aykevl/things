#include <stdint.h>

static void bitbang_update_bitplane_1(uint32_t led0, uint32_t led1, uint32_t led2, uint32_t *bitplane, uint32_t mask) {
    uint32_t tmp;
    uint32_t N7 = 7;

    __asm__ __volatile__(
        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 1 and 2 cycle bitplane
        "lsls %[tmp], #14\n\t"            // gap at the beginning
        "lsrs %[led0], #3\n\t"            // LED0 red   (A4/PB8)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 green (A5/BP7)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 blue  (A6/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #3\n\t"            // LED1 red   (A7/PB5)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 green (A8/BP4)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 blue  (A9/PB3)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #3\n\t"            // LED2 red   (A10/PB2)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 green (A11/BP1)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 blue  (A12/PB0)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"

        // Store 1-cycle and 2-cycle bitplanes.
        "eors %[tmp], %[tmp], %[mask]\n\t"
        "str %[tmp], [%[bitplane], #0]\n\t"

        // 4 and 8 cycle bitplane
        "lsls %[tmp], #14\n\t"            // gap at the beginning
        "rors %[led0], %[N7]\n\t"
        "lsrs %[led0], #10\n\t"           // LED0 red   (A4/PB8)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 green (A5/BP7)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 blue  (A6/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"
        "lsrs %[led1], #10\n\t"           // LED1 red   (A7/PB5)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 green (A8/BP4)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 blue  (A9/PB3)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"
        "lsrs %[led2], #10\n\t"           // LED2 red   (A10/PB2)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 green (A11/BP1)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 blue  (A12/PB0)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"

        // Store 4 and 8 cycle bitplanes.
        "eors %[tmp], %[tmp], %[mask]\n\t"
        "str %[tmp], [%[bitplane], #4]\n\t"

        // 16 and 32 cycle bitplane
        "lsls %[tmp], #14\n\t"            // gap at the beginning
        "rors %[led0], %[N7]\n\t"
        "lsrs %[led0], #10\n\t"           // LED0 red   (A4/PB8)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 green (A5/BP7)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 blue  (A6/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"
        "lsrs %[led1], #10\n\t"           // LED1 red   (A7/PB5)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 green (A8/BP4)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 blue  (A9/PB3)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"
        "lsrs %[led2], #10\n\t"           // LED2 red   (A10/PB2)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 green (A11/BP1)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 blue  (A12/PB0)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"

        // Store 16 and 32 cycle bitplanes.
        "eors %[tmp], %[tmp], %[mask]\n\t"
        "str %[tmp], [%[bitplane], #8]\n\t"

        : [tmp]"=&r"(tmp),
          [led0]"+&r"(led0),
          [led1]"+&r"(led1),
          [led2]"+&r"(led2)
        : [bitplane]"r"(bitplane),
          [mask]"r"(mask),
          [N7]"r"(N7)
    );
}

static void bitbang_update_bitplane_4(uint32_t led0, uint32_t led1, uint32_t led2, uint32_t *bitplane, uint32_t mask) {
    uint32_t tmp;
    uint32_t N7 = 7;

    __asm__ __volatile__(
        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 1 and 2 cycle bitplane
        "lsls %[tmp], #8\n\t"             // gap at the beginning
        "lsrs %[led0], #3\n\t"            // LED0 red   (A1/PB11)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 green (A2/BP10)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 blue  (A3/PB9)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #3\n\t"            // LED1 red   (A4/PB8)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 green (A5/BP7)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 blue  (A6/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #3\n\t"            // LED2 red   (A7/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 green (A8/BP4)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 blue  (A9/PB3)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsls %[tmp], #6\n\t"             // gap at the end

        // Store 1-cycle and 2-cycle bitplanes.
        "eors %[tmp], %[tmp], %[mask]\n\t"
        "str %[tmp], [%[bitplane], #0]\n\t"

        // 4 and 8 cycle bitplane
        "lsls %[tmp], #8\n\t"             // gap at the beginning
        "rors %[led0], %[N7]\n\t"
        "lsrs %[led0], #10\n\t"           // LED0 red   (A1/PB11)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 green (A2/BP10)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 blue  (A3/PB9)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"
        "lsrs %[led1], #10\n\t"           // LED1 red   (A4/PB8)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 green (A5/BP7)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 blue  (A6/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"
        "lsrs %[led2], #10\n\t"           // LED2 red   (A7/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 green (A8/BP4)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 blue  (A9/PB3)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsls %[tmp], #6\n\t"             // gap at the end

        // Store 4 and 8 cycle bitplanes.
        "eors %[tmp], %[tmp], %[mask]\n\t"
        "str %[tmp], [%[bitplane], #4]\n\t"

        // 16 and 32 cycle bitplane
        "lsls %[tmp], #8\n\t"             // gap at the beginning
        "rors %[led0], %[N7]\n\t"
        "lsrs %[led0], #10\n\t"           // LED0 red   (A1/PB11)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 green (A2/BP10)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led0], %[N7]\n\t"         // LED0 blue  (A3/PB9)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led0], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"
        "lsrs %[led1], #10\n\t"           // LED1 red   (A4/PB8)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 green (A5/BP7)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led1], %[N7]\n\t"         // LED1 blue  (A6/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led1], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"
        "lsrs %[led2], #10\n\t"           // LED2 red   (A7/PB6)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 green (A8/BP4)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "rors %[led2], %[N7]\n\t"         // LED2 blue  (A9/PB3)
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsrs %[led2], #1\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t"
        "lsls %[tmp], #6\n\t"             // gap at the end

        // Store 16 and 32 cycle bitplanes.
        "eors %[tmp], %[tmp], %[mask]\n\t"
        "str %[tmp], [%[bitplane], #8]\n\t"

        : [tmp]"=&r"(tmp),
          [led0]"+&r"(led0),
          [led1]"+&r"(led1),
          [led2]"+&r"(led2)
        : [bitplane]"r"(bitplane),
          [mask]"r"(mask),
          [N7]"r"(N7)
    );
}

// This is one giant assembly function that will update all LEDs.
static void bitbang_show_leds(volatile uint32_t *bitplanes, volatile void *gpio) {
    uint32_t v1, v2, v3;

    // Helpers to construct MODER values.
    // The default is to set all pins to analog mode, except for LED_VCC_ON
    // which must be kept as an output at all times.
    //                        aaaaaaaabbbbbbbbccccccccdddddddd
    uint32_t mask         = 0b11111101010101010101010101010101; // ORed with bitplanes
    uint32_t default_mode = 0b11111101111111111111111111111111; // used to clear bitplanes

    // Using the GPIO registers directly with memory offsets to avoid using one
    // more register:
    // MODER:  [%[gpio], #0x0]
    // OTYPER: [%[gpio], #0x4]
    // ODR:    [%[gpio], #0x14]

    __asm__ __volatile__(
        // Prepare for first anode.
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A1/PB11 high
        "lsls %[v2], %[v2], #11\n\t"

        // Update LEDs with anode A1/PB11.
        "strh %[v2], [%[gpio], #0x14]\n\t"   // [1] update ODR to set only the current anode high
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 32 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #0]\n\t" // [2] load 2 and 1 cycle bitplanes
        "lsls %[v2], %[v1], #1\n\t"          // [2] prepare 2 cycle bitplane
        "orrs %[v2], %[v2], %[mask]\n\t"
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 1 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #4]\n\t" // [2] load 8 (and 4) cycle bitplanes
        "lsls %[v3], %[v3], #1\n\t"          // [2] prepare 8 cycle bitplane
        "orrs %[v3], %[v3], %[mask]\n\t"
        "nop\n\t"                            // [22] nop
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
        "nop\n\t"
        "nop\n\t"
        "str  %[v2], [%[gpio], #0]\n\t"      // [1] ** start 2 cycle bitplane
        "nop\n\t"                            // [1] nop
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 1 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 8 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #4]\n\t" // [2] load (8 and) 4 cycle bitplanes
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "nop\n\t"                            // [4] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 4 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #8]\n\t" // [2] load (32 and) 16 cycle bitplanes
        "orrs %[v3], %[v3], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 16 cycle bitplane
        "adds %[bitplanes], #12\n\t"         // [1] increment bitplanes for next 3 words
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes for next anode
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A2/PB10 high
        "lsls %[v2], %[v2], #10\n\t"
        "nop\n\t"                            // [8] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[defmode], [%[gpio], #0]\n\t" // [1] ** turn off all LEDs

        // Update LEDs with anode A2/PB10.
        "strh %[v2], [%[gpio], #0x14]\n\t"   // [1] update ODR to set only the current anode high
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 32 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #0]\n\t" // [2] load 2 and 1 cycle bitplanes
        "lsls %[v2], %[v1], #1\n\t"          // [2] prepare 2 cycle bitplane
        "orrs %[v2], %[v2], %[mask]\n\t"
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 1 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #4]\n\t" // [2] load 8 (and 4) cycle bitplanes
        "lsls %[v3], %[v3], #1\n\t"          // [2] prepare 8 cycle bitplane
        "orrs %[v3], %[v3], %[mask]\n\t"
        "nop\n\t"                            // [22] nop
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
        "nop\n\t"
        "nop\n\t"
        "str  %[v2], [%[gpio], #0]\n\t"      // [1] ** start 2 cycle bitplane
        "nop\n\t"                            // [1] nop
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 1 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 8 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #4]\n\t" // [2] load (8 and) 4 cycle bitplanes
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "nop\n\t"                            // [4] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 4 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #8]\n\t" // [2] load (32 and) 16 cycle bitplanes
        "orrs %[v3], %[v3], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 16 cycle bitplane
        "adds %[bitplanes], #12\n\t"         // [1] increment bitplanes for next 3 words
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes for next anode
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A3/PB9 high
        "lsls %[v2], %[v2], #9\n\t"
        "nop\n\t"                            // [8] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[defmode], [%[gpio], #0]\n\t" // [1] ** turn off all LEDs

        // Update LEDs with anode A3/PB9.
        "strh %[v2], [%[gpio], #0x14]\n\t"   // [1] update ODR to set only the current anode high
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 32 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #0]\n\t" // [2] load 2 and 1 cycle bitplanes
        "lsls %[v2], %[v1], #1\n\t"          // [2] prepare 2 cycle bitplane
        "orrs %[v2], %[v2], %[mask]\n\t"
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 1 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #4]\n\t" // [2] load 8 (and 4) cycle bitplanes
        "lsls %[v3], %[v3], #1\n\t"          // [2] prepare 8 cycle bitplane
        "orrs %[v3], %[v3], %[mask]\n\t"
        "nop\n\t"                            // [22] nop
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
        "nop\n\t"
        "nop\n\t"
        "str  %[v2], [%[gpio], #0]\n\t"      // [1] ** start 2 cycle bitplane
        "nop\n\t"                            // [1] nop
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 1 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 8 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #4]\n\t" // [2] load (8 and) 4 cycle bitplanes
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "nop\n\t"                            // [4] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 4 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #8]\n\t" // [2] load (32 and) 16 cycle bitplanes
        "orrs %[v3], %[v3], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 16 cycle bitplane
        "adds %[bitplanes], #12\n\t"         // [1] increment bitplanes for next 3 words
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes for next anode
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A4/PB8 high
        "lsls %[v2], %[v2], #8\n\t"
        "nop\n\t"                            // [8] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[defmode], [%[gpio], #0]\n\t" // [1] ** turn off all LEDs

        // XXX increment bitplane to skip a bunch of LEDs
        "adds %[bitplanes], #72\n\t"

        // XXX prepare for next anode
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes for next anode
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A10/PB2 high
        "lsls %[v2], %[v2], #2\n\t"

        // Update LEDs with anode A10/PB2.
        "strh %[v2], [%[gpio], #0x14]\n\t"   // [1] update ODR to set only the current anode high
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 32 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #0]\n\t" // [2] load 2 and 1 cycle bitplanes
        "lsls %[v2], %[v1], #1\n\t"          // [2] prepare 2 cycle bitplane
        "orrs %[v2], %[v2], %[mask]\n\t"
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 1 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #4]\n\t" // [2] load 8 (and 4) cycle bitplanes
        "lsls %[v3], %[v3], #1\n\t"          // [2] prepare 8 cycle bitplane
        "orrs %[v3], %[v3], %[mask]\n\t"
        "nop\n\t"                            // [22] nop
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
        "nop\n\t"
        "nop\n\t"
        "str  %[v2], [%[gpio], #0]\n\t"      // [1] ** start 2 cycle bitplane
        "nop\n\t"                            // [1] nop
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 1 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 8 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #4]\n\t" // [2] load (8 and) 4 cycle bitplanes
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "nop\n\t"                            // [4] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 4 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #8]\n\t" // [2] load (32 and) 16 cycle bitplanes
        "orrs %[v3], %[v3], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 16 cycle bitplane
        "adds %[bitplanes], #12\n\t"         // [1] increment bitplanes for next 3 words
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes for next anode
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A11/PB1 high
        "lsls %[v2], %[v2], #1\n\t"
        "nop\n\t"                            // [8] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[defmode], [%[gpio], #0]\n\t" // [1] ** turn off all LEDs

        // Update LEDs with anode A11/PB1.
        "strh %[v2], [%[gpio], #0x14]\n\t"   // [1] update ODR to set only the current anode high
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 32 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #0]\n\t" // [2] load 2 and 1 cycle bitplanes
        "lsls %[v2], %[v1], #1\n\t"          // [2] prepare 2 cycle bitplane
        "orrs %[v2], %[v2], %[mask]\n\t"
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 1 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #4]\n\t" // [2] load 8 (and 4) cycle bitplanes
        "lsls %[v3], %[v3], #1\n\t"          // [2] prepare 8 cycle bitplane
        "orrs %[v3], %[v3], %[mask]\n\t"
        "nop\n\t"                            // [22] nop
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
        "nop\n\t"
        "nop\n\t"
        "str  %[v2], [%[gpio], #0]\n\t"      // [1] ** start 2 cycle bitplane
        "nop\n\t"                            // [1] nop
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 1 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 8 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #4]\n\t" // [2] load (8 and) 4 cycle bitplanes
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "nop\n\t"                            // [4] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 4 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #8]\n\t" // [2] load (32 and) 16 cycle bitplanes
        "orrs %[v3], %[v3], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 16 cycle bitplane
        "adds %[bitplanes], #12\n\t"         // [1] increment bitplanes for next 3 words
        "ldr  %[v1], [%[bitplanes], #8]\n\t" // [2] load 32 (and 16) cycle bitplanes for next anode
        "lsls %[v1], %[v1], #1\n\t"          // [2] prepare 32 cycle bitplane
        "orrs %[v1], %[v1], %[mask]\n\t"
        "movs %[v2], #1\n\t"                 // [2] load bits to set only A12/PB0 high
        "lsls %[v2], %[v2], #0\n\t"
        "nop\n\t"                            // [8] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[defmode], [%[gpio], #0]\n\t" // [1] ** turn off all LEDs

        // Update LEDs with anode A12/PB0.
        "strh %[v2], [%[gpio], #0x14]\n\t"   // [1] update ODR to set only the current anode high
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 32 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #0]\n\t" // [2] load 2 and 1 cycle bitplanes
        "lsls %[v2], %[v1], #1\n\t"          // [2] prepare 2 cycle bitplane
        "orrs %[v2], %[v2], %[mask]\n\t"
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 1 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #4]\n\t" // [2] load 8 (and 4) cycle bitplanes
        "lsls %[v3], %[v3], #1\n\t"          // [2] prepare 8 cycle bitplane
        "orrs %[v3], %[v3], %[mask]\n\t"
        "nop\n\t"                            // [22] nop
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
        "nop\n\t"
        "nop\n\t"
        "str  %[v2], [%[gpio], #0]\n\t"      // [1] ** start 2 cycle bitplane
        "nop\n\t"                            // [1] nop
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 1 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 8 cycle bitplane
        "ldr  %[v1], [%[bitplanes], #4]\n\t" // [2] load (8 and) 4 cycle bitplanes
        "orrs %[v1], %[v1], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "nop\n\t"                            // [4] nop
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "str  %[v1], [%[gpio], #0]\n\t"      // [1] ** start 4 cycle bitplane
        "ldr  %[v3], [%[bitplanes], #8]\n\t" // [2] load (32 and) 16 cycle bitplanes
        "orrs %[v3], %[v3], %[mask]\n\t"     // [1] prepare 4 cycle bitplane
        "str  %[v3], [%[gpio], #0]\n\t"      // [1] ** start 16 cycle bitplane
        "nop\n\t"                            // [15] nop
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
        "str  %[defmode], [%[gpio], #0]\n\t" // [1] ** turn off all LEDs

        : [v1]"=&r"(v1),
          [v2]"=&r"(v2),
          [v3]"=&r"(v3),
          [bitplanes]"+&r"(bitplanes)
        : [gpio]"r"(gpio),
          [defmode]"r"(default_mode),
          [mask]"r"(mask)
    );
}
