#include <stdint.h>

static void bitbang_update_bitplane_1(uint32_t led0, uint32_t led1, uint32_t led2, uint32_t *bitplane) {
    uint32_t tmp;
    uint32_t N8 = 8;
    uint32_t xor = 0b1000011111111111;

    // bitplane[0] lo bits: 2 cycle bitplane
    // bitplane[0] hi bits: 1 cycle bitplane
    // bitplane[1] lo bits: 4 cycle bitplane
    // bitplane[1] hi bits: 8 cycle bitplane

    __asm__ __volatile__(
        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 1 cycle bitplane
        "lsrs %[led0], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "lsls %[tmp], #3\n\t"             // gap between PA9 and PA5
        "lsrs %[led2], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A7/PA5)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A8/PA4)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A9/PA3)
        "lsrs %[led1], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // 2 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "lsls %[tmp], #3\n\t"             // gap between PA9 and PA5
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A7/PA5)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A8/PA4)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A9/PA3)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // Store 1-cycle and 2-cycle bitplanes.
        "str %[tmp], [%[bitplane]]\n\t"

        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 4 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "lsls %[tmp], #3\n\t"             // gap between PA9 and PA5
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A7/PA5)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A8/PA4)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A9/PA3)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // 8 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "lsls %[tmp], #3\n\t"             // gap between PA9 and PA5
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A7/PA5)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A8/PA4)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A9/PA3)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // Swap 4-cycle and 8-cycle bitplane (by rotating by 16 bits).
        "rors %[tmp], %[N8]\n\t"
        "rors %[tmp], %[N8]\n\t"

        // Store 4-cycle and 8-cycle bitplanes.
        "str %[tmp], [%[bitplane], #4]\n\t"

        : [tmp]"=&r"(tmp),
          [led0]"+&r"(led0),
          [led1]"+&r"(led1),
          [led2]"+&r"(led2)
        : [bitplane]"r"(bitplane),
          [N8]"r"(N8),
          [xor]"r"(xor)
    );
}

static void bitbang_update_bitplane_2(uint32_t led0, uint32_t led1, uint32_t led2, uint32_t *bitplane) {
    uint32_t tmp;
    uint32_t N8 = 8;
    uint32_t xor = 0b1000011111111111;

    // bitplane[0] lo bits: 2 cycle bitplane
    // bitplane[0] hi bits: 1 cycle bitplane
    // bitplane[1] lo bits: 4 cycle bitplane
    // bitplane[1] hi bits: 8 cycle bitplane

    __asm__ __volatile__(
        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 1 cycle bitplane
        "lsrs %[led0], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "lsrs %[led2], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "lsls %[tmp], #3\n\t"             // gap between PA6 and PA2
        "lsrs %[led1], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // 2 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "lsls %[tmp], #3\n\t"             // gap between PA6 and PA2
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // Store 1-cycle and 2-cycle bitplanes.
        "str %[tmp], [%[bitplane]]\n\t"

        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 4 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "lsls %[tmp], #3\n\t"             // gap between PA6 and PA2
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // 8 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "lsls %[tmp], #3\n\t"             // gap between PA6 and PA2
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A10/PA2)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A11/PA1)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // Swap 4-cycle and 8-cycle bitplane (by rotating by 16 bits).
        "rors %[tmp], %[N8]\n\t"
        "rors %[tmp], %[N8]\n\t"

        // Store 4-cycle and 8-cycle bitplanes.
        "str %[tmp], [%[bitplane], #4]\n\t"

        : [tmp]"=&r"(tmp),
          [led0]"+&r"(led0),
          [led1]"+&r"(led1),
          [led2]"+&r"(led2)
        : [bitplane]"r"(bitplane),
          [N8]"r"(N8),
          [xor]"r"(xor)
    );
}

static void bitbang_update_bitplane_3(uint32_t led0, uint32_t led1, uint32_t led2, uint32_t *bitplane) {
    uint32_t tmp;
    uint32_t N8 = 8;
    uint32_t xor = 0b1000011111111111;

    // bitplane[0] lo bits: 2 cycle bitplane
    // bitplane[0] hi bits: 1 cycle bitplane
    // bitplane[1] lo bits: 4 cycle bitplane
    // bitplane[1] hi bits: 8 cycle bitplane

    __asm__ __volatile__(
        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 1 cycle bitplane
        "lsrs %[led0], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "lsrs %[led2], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "lsrs %[led1], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "lsls %[tmp], #3\n\t"             // gap at the start (PA3-0)
        "eors %[tmp], %[xor]\n\t"

        // 2 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "lsls %[tmp], #3\n\t"             // gap at the start (PA3-0)
        "eors %[tmp], %[xor]\n\t"

        // Store 1-cycle and 2-cycle bitplanes.
        "str %[tmp], [%[bitplane]]\n\t"

        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 4 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "lsls %[tmp], #3\n\t"             // gap at the start (PA3-0)
        "eors %[tmp], %[xor]\n\t"

        // 8 cycle bitplane
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A1/PA15)
        "lsls %[tmp], #4\n\t"             // gap between PA15 and PA10
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A2/PA10)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A3/PA9)
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "lsls %[tmp], #3\n\t"             // gap at the start (PA3-0)
        "eors %[tmp], %[xor]\n\t"

        // Swap 4-cycle and 8-cycle bitplane (by rotating by 16 bits).
        "rors %[tmp], %[N8]\n\t"
        "rors %[tmp], %[N8]\n\t"

        // Store 4-cycle and 8-cycle bitplanes.
        "str %[tmp], [%[bitplane], #4]\n\t"

        : [tmp]"=&r"(tmp),
          [led0]"+&r"(led0),
          [led1]"+&r"(led1),
          [led2]"+&r"(led2)
        : [bitplane]"r"(bitplane),
          [N8]"r"(N8),
          [xor]"r"(xor)
    );
}

static void bitbang_update_bitplane_4(uint32_t led0, uint32_t led1, uint32_t led2, uint32_t *bitplane) {
    uint32_t tmp;
    uint32_t N8 = 8;
    uint32_t xor = 0b1000011111111111;

    // bitplane[0] lo bits: 2 cycle bitplane
    // bitplane[0] hi bits: 1 cycle bitplane
    // bitplane[1] lo bits: 4 cycle bitplane
    // bitplane[1] hi bits: 8 cycle bitplane

    __asm__ __volatile__(
        // 1 cycle bitplane
        "lsls %[tmp], #7\n\t"             // gap at the start
        "lsrs %[led2], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "lsrs %[led1], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "lsrs %[led0], #5\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A10/PA2)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A11/PA1)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // 2 cycle bitplane
        "lsls %[tmp], #7\n\t"             // gap at the start
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A10/PA2)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A11/PA1)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // Store 1-cycle and 2-cycle bitplanes.
        "str %[tmp], [%[bitplane]]\n\t"

        // Not clearing %[tmp] here since we will be shifting all 32 bits out.

        // 4 cycle bitplane
        "lsls %[tmp], #7\n\t"             // gap at the start
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A10/PA2)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A11/PA1)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // 8 cycle bitplane
        "lsls %[tmp], #7\n\t"             // gap at the start
        "rors %[led2], %[N8]\n\t"
        "lsrs %[led2], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 red   (A4/PA8)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 green (A5/PA7)
        "rors %[led2], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED2 blue  (A6/PA6)
        "rors %[led1], %[N8]\n\t"
        "lsrs %[led1], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 red   (A7/PA5)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 green (A8/PA4)
        "rors %[led1], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED1 blue  (A9/PA3)
        "rors %[led0], %[N8]\n\t"
        "lsrs %[led0], #9\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 red   (A10/PA2)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 green (A11/PA1)
        "rors %[led0], %[N8]\n\t"
        "adcs %[tmp], %[tmp], %[tmp]\n\t" // LED0 blue  (A12/PA0)
        "eors %[tmp], %[xor]\n\t"

        // Swap 4-cycle and 8-cycle bitplane (by rotating by 16 bits).
        "rors %[tmp], %[N8]\n\t"
        "rors %[tmp], %[N8]\n\t"

        // Store 4-cycle and 8-cycle bitplanes.
        "str %[tmp], [%[bitplane], #4]\n\t"

        : [tmp]"=&r"(tmp),
          [led0]"+&r"(led0),
          [led1]"+&r"(led1),
          [led2]"+&r"(led2)
        : [bitplane]"r"(bitplane),
          [N8]"r"(N8),
          [xor]"r"(xor)
    );
}

static void bitbang_show_leds(volatile uint32_t *bitplanes, volatile uint16_t *out) {
    uint32_t v1, v2;
    uint32_t mask = 0b1000011111111111;

    // TODO: make this one giant assembly function that updates all LEDs in one
    // asm statement.

    __asm__ __volatile__(
        "ldr %[v1], [%[bitplanes], #0]\n\t"
        "ldr %[v2], [%[bitplanes], #4]\n\t"
        "strh %[v1], [%[out]]\n\t" // 2 cycle bitplane
        "lsrs %[v1], #16\n\t"
        "strh %[v1], [%[out]]\n\t" // 1 cycle bitplane
        "strh %[v2], [%[out]]\n\t" // 4 cycle bitplane
        "lsrs %[v2], #16\n\t"
        "nop\n\t"
        "nop\n\t"
        "strh %[v2], [%[out]]\n\t" // 8 cycle bitplane
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "nop\n\t"
        "strh %[mask], [%[out]]\n\t"
        : [v1]"=&r"(v1),
          [v2]"=&r"(v2),
          [bitplanes]"+&r"(bitplanes)
        : [out]"r"(out),
          [mask]"r"(mask)
    );
}
