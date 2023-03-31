package main

// Small brick breaker game for the Adafruit PyBadge.

import (
	"image/color"
	"strconv"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/pixel"
	"github.com/aykevl/tinygl/style"
	"tinygo.org/x/tinyfont/freesans"
)

func main() {
	board.Buttons.Configure()
	display := board.Display.Configure()
	runUI(display)
}

func runUI[T pixel.Color](display board.Displayer[T]) {
	width, height := display.Size()

	// Base style (100% scale, blue background, white foreground).
	font := &freesans.Regular9pt7b
	foreground := pixel.NewColor[T](0xff, 0xff, 0xff)
	background := pixel.NewColor[T](64, 64, 64)
	base := style.New(100, foreground, background, font)

	buf := make([]T, int(width)*int(height)/10)
	screen := tinygl.NewScreen(display, base, buf)

	title := tinygl.NewText(base.WithBackground(color.RGBA{R: 255, A: 255}), "")
	canvas := gfx.NewCanvas(base.WithBackground(color.RGBA{A: 255}), 96, 96)
	canvas.SetGrowable(0, 1)
	all := tinygl.NewVBox[T](base, title, canvas)
	screen.SetChild(all)

	// run brick breaker game
	screen.Layout()
	for {
		runBrickBreaker(screen, canvas, title)
	}
}

func runBrickBreaker[T pixel.Color](screen *tinygl.Screen[T], canvas *gfx.Canvas[T], title *tinygl.Text[T]) {
	// Collect some game constants.
	_, _, cw, ch := canvas.Bounds()
	const paddleWidth = 24
	const paddleHeight = 6
	const ballSize = 3 * 256
	const initialBallRotation = 256 - 40 // 64 equals 90°
	const (
		stateStart = iota
		statePlay
		stateFinished
	)

	// Initialize the game.
	ballX := -10 * 256 // off-screen, will be moved in first frame
	ballY := -10 * 256
	ballRotation := uint8(initialBallRotation)
	points := 0
	state := stateStart
	canvas.Clear()
	title.SetText("Brick Breaker")

	brickSize := cw / 13
	brickPadding := brickSize / 10
	brickInnerSize := brickSize - brickPadding*2
	brickYStart := ch/6 + brickPadding
	var bricks []*gfx.Rect[T]
	for i := 0; i < 10; i++ {
		x := cw/2 + (i-5)*brickSize + brickPadding
		y := brickYStart
		bricks = append(bricks, canvas.CreateRect(x, y, brickInnerSize, brickInnerSize, color.RGBA{R: 200, G: 200, B: 255}))
	}
	for i := 0; i < 11; i++ {
		x := cw/2 + (i-5)*brickSize - brickSize/2 + brickPadding
		y := brickYStart + brickSize
		bricks = append(bricks, canvas.CreateRect(x, y, brickInnerSize, brickInnerSize, color.RGBA{R: 200, G: 200, B: 255}))
	}
	paddle := canvas.CreateRect(cw/2-paddleWidth/2, ch-paddleHeight, paddleWidth, paddleHeight, color.RGBA{R: 255})
	ball := canvas.CreateRect(ballX/256, ballY/256, ballSize/256, ballSize/256, color.RGBA{R: 255, G: 255})

	var leftDown, rightDown bool
	for {
		frameStart := time.Now()

		// Read input (buttons etc).
		board.Buttons.ReadInput()
		for {
			event := board.Buttons.NextEvent()
			if event == board.NoKeyEvent {
				break
			}
			switch event.Key() {
			case board.KeyLeft:
				leftDown = event.Pressed()
			case board.KeyRight:
				rightDown = event.Pressed()
			case board.KeyStart, board.KeyA, board.KeySpace:
				if event.Pressed() {
					switch state {
					case stateStart:
						// Start the game.
						state = statePlay
						ballRotation = initialBallRotation
						title.SetText("0")
					case stateFinished:
						// Exit the game.
						return
					}
				}
			}
		}
		if state == stateStart || state == statePlay {
			if leftDown {
				x, y, _, _ := paddle.Bounds()
				x -= 3
				if x < 0 {
					x = 0
				}
				paddle.Move(x, y)
			}
			if rightDown {
				x, y, _, _ := paddle.Bounds()
				x += 3
				if x > cw-paddleWidth {
					x = cw - paddleWidth
				}
				paddle.Move(x, y)
			}
		}

		// Update game state (movement, etc).
		switch state {
		case stateStart:
			// Beginning of game. Move ball with paddle.
			x, y, w, _ := paddle.Bounds()
			ballX = x*256 + w*128 - ballSize/2
			ballY = y*256 - ballSize
		case statePlay, stateFinished:
			// Playing the game.
			ballMoveX, ballMoveY := rotationToVector(ballRotation)
			ballX += ballMoveX
			ballY += ballMoveY
			if ballX+ballSize > cw*256 {
				ballRotation = rotationFlipX(ballRotation)
				ballX -= (ballX + ballSize) - cw*256
			}
			if ballX < 0 {
				ballRotation = rotationFlipX(ballRotation)
				ballX -= 0
			}
			if ballY < 0 {
				ballRotation = rotationFlipY(ballRotation)
				ballY -= ballY
			}
			switch state {
			case statePlay:
				// Check whether the ball hit the ground or the paddle.
				paddleX, paddleY, paddleW, _ := paddle.Bounds()
				paddleX *= 256
				paddleY *= 256
				paddleW *= 256
				if _, ballMoveY := rotationToVector(ballRotation); ballMoveY > 0 { // moving down
					if ballY+ballSize >= paddleY && ballX >= paddleX && ballX <= paddleX+paddleW {
						// Ball is on the paddle, so bounce it.
						ballRotation = rotationFlipY(ballRotation)
						ballRotation += uint8(frameStart.UnixNano()>>16)/16 - 8 // add a bit of randomness
						ballRotation += uint8((ballX - (paddleX + paddleW/2)) / 128)
					}
					if ballY+ballSize >= ch*256 {
						// Ball fell on the ground, game over.
						state = stateFinished
						title.SetText("game over: " + strconv.Itoa(points))
					}
				}
				// Check whether the ball hit any of the bricks.
				for _, brick := range bricks {
					if brick.Hidden() {
						continue
					}
					brickX, brickY, brickW, brickH := brick.Bounds()
					// Use some Pythagoras to calculate the distance between the
					// ball and the brick! Normally you'd take the square root
					// of the result, but it's easier to just square the
					// expected distance too to avoid sqrt.
					dx := (brickX + brickW/2) - (ballX+ballSize/2)/256
					dy := (brickY + brickH/2) - (ballY+ballSize/2)/256
					distance := dx*dx + dy*dy
					if distance < brickInnerSize*brickInnerSize {
						brick.SetHidden(true)
						points++
						if points == len(bricks) {
							title.SetText("finished!")
							state = stateFinished
						} else {
							title.SetText(strconv.Itoa(points))
							// TODO: treat brick as round while bouncing (to
							// make it look more random).
							if dx > dy {
								ballRotation = rotationFlipY(ballRotation)
							} else {
								ballRotation = rotationFlipX(ballRotation)
							}
						}
					}
				}
			case stateFinished:
				// Don't continue to (visibly) move the ball once it leaves the screen.
				if ballY > ch {
					ball.SetHidden(true)
				}
			}
		}
		ball.Move(ballX/256, ballY/256)

		screen.Update()
		duration := time.Since(frameStart)
		if duration == 0 {
			println("fps: ∞")
		} else {
			println("fps:", time.Second/duration)
		}
		time.Sleep(time.Second/60 - duration) // try to hit 60fps
	}
}

// Python oneliner:
//
//	list([int(round(math.sin(n / 32 * math.pi / 2)*255)) for n in range(32)])
var sinuses = [32]uint8{0, 13, 25, 37, 50, 62, 74, 86, 98, 109, 120, 131, 142, 152, 162, 171, 180, 189, 197, 205, 212, 219, 225, 231, 236, 240, 244, 247, 250, 252, 254, 255}

// Change a rotation (256 equals 360°, starting at the right side of the unit
// circle) to a vector (where positive coordinates go down and to the right).
func rotationToVector(rotation uint8) (x, y int) {
	rotation /= 2
	switch rotation / 32 {
	case 0: // bottom right
		return int(sinuses[31-rotation]), int(sinuses[rotation])
	case 1: // bottom left
		return -int(sinuses[rotation-32]), int(sinuses[63-rotation])
	case 2: // top left
		return -int(sinuses[95-rotation]), -int(sinuses[rotation-64])
	case 3: // top right
		return int(sinuses[rotation-96]), -int(sinuses[127-rotation])
	}
	return 0, 0 // unreachable
}

// Flip the rotation across the X axis.
func rotationFlipX(rotation uint8) uint8 {
	return 128 - rotation
}

// Flip the rotation across the Y axis.
func rotationFlipY(rotation uint8) uint8 {
	return rotationFlipX(rotation+64) - 64
}