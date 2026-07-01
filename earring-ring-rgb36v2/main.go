package main

// State: current mode, current variant, previous variants, 4 hash bytes
var storedState [32]uint8

const (
	stateOffsetMode     = 0
	stateOffsetVariants = 2
)

// Saved variant for each mode (so that the variant is kept between mode
// switches).
var modeVariants = storedState[stateOffsetVariants : stateOffsetVariants+modeLast]

func main() {
	initHardware()

	// Zero all LEDs.
	for i := 0; i < 12; i++ {
		setLEDs(i, 0, 0, 0)
	}

	index := 0 // 0..11, group of 3 LEDs that will be updated together
	frame := 0
	mode := int(storedState[stateOffsetMode])
	variant := 0
	if mode < len(modeVariants) {
		variant = int(modeVariants[mode])
	}
	if animationNeedsMic(mode) {
		enableMic()
	}
	initMode(mode)
	modeFramesPressed := 0
	previousMode := 0
	variantBtnFramesPressed := 0
	for {
		led0 := animate(mode, variant, index+0, frame)
		led1 := animate(mode, variant, index+12, frame)
		led2 := animate(mode, variant, index+24, frame)

		// Read the ADC, intentionally done in between so that the LED power
		// won't interfere too much with the microphone.
		if animationNeedsMic(mode) {
			updateMic(index)
		}

		// Update 3 LEDs at a time, since that's convenient for the
		// RGB-to-bitplane conversion.
		setLEDs(index, uint32(led0), uint32(led1), uint32(led2))

		// Bitbang the LEDs.
		updateLEDs()

		index++
		if index == 12 {
			index = 0
			frame++

			// Do per animation process, such as processing audio samples
			// collected during the previous frame.
			newFrame(mode, variant, frame)

			// Read the mode button every frame update.
			modePressed := button1Pressed()
			if modeFramesPressed == 30 {
				// Always disable the microphone when sleeping.
				disableMic()

				turnOffAnimation(mode, variant, frame)

				// Sleep until the mode button is pressed.
				sleepUntilButtonPress()

				// To continue the startup animation, set the mode to "power
				// on".
				previousMode = mode
				mode = modePowerOn
				frame = 0
				modeFramesPressed = -0x8000_0000 // don't switch to the next animation on button release
			}
			if mode == modePowerOn {
				if frame == numLEDs/2 {
					mode = previousMode

					// Woke up again, so start up interfaces.
					if animationNeedsMic(mode) {
						enableMic()
					}
				}
			} else {
				if !modePressed && modeFramesPressed > 0 {
					// Move to the next mode.
					mode++
					if mode >= modeLast {
						// Last, so wrap around.
						mode = 0
					}
					variant = int(modeVariants[mode])

					// Save the current state to flash.
					storedState[stateOffsetMode] = uint8(mode)
					saveState()

					// Clear LEDs before moving on to the next mode.
					for i := 0; i < 12; i++ {
						setLEDs(i, 0, 0, 0)
					}

					// Only enable the microphone when needed.
					if animationNeedsMic(mode) {
						enableMic()
					} else {
						disableMic()
					}

					initMode(mode)
				}

				// Switch to the next animation variant when the variant button
				// is released.
				variantBtnPressed := button2Pressed()
				if !variantBtnPressed && variantBtnFramesPressed != 0 {
					newVariant := animationNextVariant(mode, variant)
					if newVariant != variant {
						variant = newVariant
						if mode < len(modeVariants) {
							modeVariants[mode] = uint8(variant)
						}
						saveState()
					}
				}
				if variantBtnPressed {
					variantBtnFramesPressed++
				} else {
					variantBtnFramesPressed = 0
				}
				if variantBtnFramesPressed == 30 {
					slot := mode - modeCustom0
					if slot >= 0 && slot < 3 {
						// Long press is to receive data, not to switch to the
						// next variant.
						variantBtnFramesPressed = 0

						// Receive the data, and write it to flash.
						dataRecv(slot)

						// Immediately start running this pattern.
						customLoadPattern(slot, false)

						// Reset variant to 0, since the newly received one
						// might have fewer variants than the one previously
						// loaded.
						variant = 0
						modeVariants[mode] = uint8(variant)
						saveState()
					}
				}
			}
			if modePressed {
				modeFramesPressed++
			} else {
				modeFramesPressed = 0
			}
		}
	}
}

// Show an animation during shutdown. It keeps the current animation, but
// freezes it in place, and turns off LEDs in sequence.
func turnOffAnimation(mode, variant, frame int) {
	for i := range 36 / 3 {
		for index := range 12 {
			// Shut down LEDs in groups of 3.
			// We have to calculate the animation for two reasons:
			//  1. The animation may change in brightness if we don't call
			//     animate() to slow it down the same way as the normal
			//     animation.
			//  2. The 3 LEDs we can update at a time are spread over the ring
			//     of LEDs, not sequentially next to each other.
			led0 := animate(mode, variant, index+0, frame)
			if index+0 < i*3 {
				led0 = 0
			}
			led1 := animate(mode, variant, index+12, frame)
			if index+12 < i*3 {
				led1 = 0
			}
			led2 := animate(mode, variant, index+24, frame)
			if index+24 < i*3 {
				led2 = 0
			}
			setLEDs(index, uint32(led0), uint32(led1), uint32(led2))

			// Bitbang the LEDs.
			updateLEDs()
		}
	}
}
