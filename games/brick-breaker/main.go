package main

// Small brick breaker game for the Adafruit PyBadge.

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/gfx"
	"github.com/aykevl/tinygl/pixel"
	"github.com/aykevl/tinygl/style"
	"tinygo.org/x/drivers/shifter"
	"tinygo.org/x/drivers/st7735"
	"tinygo.org/x/tinyfont/freesans"
)

var buttons shifter.Device

func main() {
	time.Sleep(time.Second)
	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       machine.SPI1_SCK_PIN,
		SDO:       machine.SPI1_SDO_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		Frequency: 15_000_000, // datasheet says 66ns (~15.15MHz) is the max speed
	})
	display := st7735.New(machine.SPI1, machine.TFT_RST, machine.TFT_DC, machine.TFT_CS, machine.TFT_LITE)
	display.Configure(st7735.Config{
		Rotation: st7735.ROTATION_90,
	})
	width, height := display.Size()

	// Base style (100% scale, blue background, white foreground).
	font := &freesans.Regular9pt7b
	foreground := pixel.NewRGB565BE(0xff, 0xff, 0xff)
	background := pixel.NewRGB565BE(64, 64, 64)
	base := style.New(100, foreground, background, font)

	buf := make([]pixel.RGB565BE, width*height/10)
	screen := tinygl.NewScreen(&display, base, buf)

	title := tinygl.NewText(base.WithBackground(color.RGBA{R: 255, A: 255}), "")
	canvas := gfx.NewCanvas(base.WithBackground(color.RGBA{A: 255}), 96, 96)
	all := tinygl.NewVBox[pixel.RGB565BE](base, title, canvas)
	screen.SetChild(all)

	buttons = shifter.NewButtons()
	buttons.Configure()

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
	const initialBallSpeedX = 200
	const initialBallSpeedY = -300
	const (
		stateStart = iota
		statePlay
		stateFinished
	)

	// Initialize the game.
	ballX := -10 * 256 // off-screen, will be moved in first frame
	ballY := -10 * 256
	ballMoveX := 0
	ballMoveY := 0
	points := 0
	state := 0
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

	startPressed := true // usually false, but when restarting the game, 'start' is still pressed
	for {
		frameStart := time.Now()

		// Read input (buttons etc).
		buttons.ReadInput()
		if state == stateStart || state == statePlay {
			if buttons.Pins[shifter.BUTTON_LEFT].Get() {
				x, y, _, _ := paddle.Bounds()
				x -= 3
				if x < 0 {
					x = 0
				}
				paddle.Move(x, y)
			}
			if buttons.Pins[shifter.BUTTON_RIGHT].Get() {
				x, y, _, _ := paddle.Bounds()
				x += 3
				if x > cw-paddleWidth {
					x = cw - paddleWidth
				}
				paddle.Move(x, y)
			}
		}
		if buttons.Pins[shifter.BUTTON_START].Get() {
			if !startPressed {
				switch state {
				case stateStart:
					// Start the game.
					state = statePlay
					ballMoveX = initialBallSpeedX
					ballMoveY = initialBallSpeedY
					title.SetText("0")
				case stateFinished:
					// Exit the game.
					return
				}
			}
			startPressed = true
		} else {
			startPressed = false
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
			ballX += ballMoveX
			ballY += ballMoveY
			if ballX+ballSize > cw*256 {
				ballMoveX = -ballMoveX
				ballX -= (ballX + ballSize) - cw*256
			}
			if ballX < 0 {
				ballMoveX = -ballMoveX
				ballX -= 0
			}
			if ballY < 0 {
				ballMoveY = -ballMoveY
				ballY -= ballY
			}
			switch state {
			case statePlay:
				// Check whether the ball hit the ground or the paddle.
				paddleX, paddleY, paddleW, _ := paddle.Bounds()
				paddleX *= 256
				paddleY *= 256
				paddleW *= 256
				if ballMoveY > 0 { // moving down
					if ballY+ballSize >= paddleY && ballX >= paddleX && ballX <= paddleX+paddleW {
						// Ball is on the paddle, so bounce it.
						ballMoveY = -ballMoveY
						ballY -= (ballY + ballSize) - paddleY
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
						}
						title.SetText(strconv.Itoa(points))
						// TODO: treat brick as round while bouncing
						if dx > dy {
							ballMoveY = -ballMoveY
						} else {
							ballMoveX = -ballMoveX
						}
					}
				}
			case stateFinished:
				// Make sure the ball doesn't continue moving after it left the screen.
				if ballY > ch {
					ballMoveX = 0
					ballMoveY = 0
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
