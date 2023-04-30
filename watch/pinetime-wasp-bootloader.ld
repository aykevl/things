/* Slightly modified version of targets/nrf52-s132v6.ld in the TinyGo source
 * code. It changes the RAM start address to 0x200039e0 to avoid PNVRAM area.
 * For details, see:
 * https://wasp-os.readthedocs.io/en/latest/wasp.html#pnvram-protocol
 */

/* TODO: exclude bootloader area from FLASH_TEXT */

MEMORY
{
    FLASH_TEXT (rw) : ORIGIN = 0x00000000 + 0x00026000 , LENGTH = 512K - 0x00026000 /* .text */
    RAM (xrw)       : ORIGIN = 0x20000000 + 0x000039e0,  LENGTH = 64K  - 0x000039e0
}

_stack_size = 4K;

/* This value is needed by the Nordic SoftDevice. */
__app_ram_base = ORIGIN(RAM);

INCLUDE "targets/arm.ld"
