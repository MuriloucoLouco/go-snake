package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 400
	screenHeight = 400
	gridWidth    = 20
	gridHeight   = 20
	blockWidth   = screenWidth / gridWidth
	blockHeight  = screenHeight / gridHeight
	centerX      = int(gridWidth / 2)
	centerY      = int(gridHeight / 2)
)

func contains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}

func loadTextureFromBMP(spriteName string, renderer *sdl.Renderer) (texture *sdl.Texture) {

	img, err := sdl.LoadBMP(spriteName)
	if err != nil {
		panic(fmt.Errorf("loading sprite: %v", err))
	}
	defer img.Free()

	texture, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture: %v", err))
	}
	return texture
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"snake",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL,
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

	s := createSnake(renderer)
	a := createApple(renderer, s)

	lastTime := time.Now().UnixNano() / int64(time.Millisecond)
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		s.eat(&a)
		s.move()
		tempNow := time.Now().UnixNano() / int64(time.Millisecond)
		if tempNow > lastTime+125 {
			lastTime = tempNow
			s.update()

		}

		a.render(renderer)
		s.render(renderer)

		renderer.Present()
	}
}
