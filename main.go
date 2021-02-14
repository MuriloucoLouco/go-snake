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
)

type gameState struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	snake    *snake
	apple    *apple
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

	s := createSnake(state)
	state.snake = &s
	a := createApple(state)
	state.apple = &a

	lastTime := time.Now().UnixNano() / int64(time.Millisecond)
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		padX, padY, _, screenX, screenY := getRightSize(window)
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.FillRect(&sdl.Rect{X: padX, Y: padY, W: screenX, H: screenY})

		s.eat(state)
		window.SetTitle("Snake - score: " + fmt.Sprintf("%d", s.score))

		s.move()
		tempNow := time.Now().UnixNano() / int64(time.Millisecond)
		if tempNow > lastTime+125 {
			lastTime = tempNow
			s.update(state)
		}

		a.render(state)
		s.render(state)

		renderer.Present()
	}
}
