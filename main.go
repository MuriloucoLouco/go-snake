package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	gridWidth  = 20
	gridHeight = 20
	centerX    = int(gridWidth / 2)
	centerY    = int(gridHeight / 2)
	speed      = 125
)

type textureState struct {
    snake *sdl.Texture
    apple *sdl.Texture
}

type gameState struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	snake    *snake
	apple    *apple
    textures textureState
}

func render(state gameState) {
	for {
		state.snake.update(state)

		padX, padY, _, screenWidth, screenHeight := getRightSize(state.window)
		state.renderer.SetDrawColor(0, 0, 0, 255)
		state.renderer.Clear()
		state.renderer.SetDrawColor(255, 255, 255, 255)
		state.renderer.FillRect(&sdl.Rect{X: int32(padX), Y: int32(padY), W: int32(screenWidth), H: int32(screenHeight)})
		state.renderer.SetDrawColor(0, 0, 0, 255)
		state.renderer.FillRect(&sdl.Rect{X: int32(padX) + 1, Y: int32(padY) + 1, W: int32(screenWidth) - 2, H: int32(screenHeight) - 2})

		state.snake.render(state)
		state.apple.render(state)

		state.renderer.Present()

		state.window.SetTitle("Snake - score: " + fmt.Sprintf("%d", state.snake.score))
		time.Sleep(speed * time.Millisecond)
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	screenWidth := int32(400)
	screenHeight := int32(400)

	window, err := sdl.CreateWindow(
		"snake",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	state := gameState{
		window:   window,
		renderer: renderer,
	}
    
    state.textures.snake = loadTextureFromBMP("sprites/snake.bmp", state.renderer)
    state.textures.apple = loadTextureFromBMP("sprites/apple.bmp", state.renderer)

	s := createSnake(state)
	state.snake = &s
	a := createApple(state)
	state.apple = &a

	go render(state)

	for {
		sdl.WaitEvent()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		s.move()
	}
}
