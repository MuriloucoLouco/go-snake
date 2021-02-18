package main

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type apple struct {
	posX, posY int
}

func canPlaceApple(state gameState) bool {
	width := state.config.GridWidth
	height := state.config.GridHeight
	positions := make([][]bool, height)
	for i := range positions {
		positions[i] = make([]bool, width)
	}

	for _, apple := range state.apples {
		positions[apple.posY][apple.posX] = true
	}
	for _, snakePos := range state.snake.positions {
		positions[snakePos[1]][snakePos[0]] = true
	}
	var numberOfBlanks int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !positions[y][x] {
				numberOfBlanks++
			}
		}
	}
	if numberOfBlanks >= 2 {
		return true
	}
	return false
}

func randomPlace(state gameState) (x, y int) {
	rand.Seed(time.Now().UnixNano())
	x = rand.Intn(state.config.GridWidth)
	y = rand.Intn(state.config.GridHeight)

	for _, position := range state.snake.positions {
		if x == position[0] && y == position[1] {
			x, y = randomPlace(state)
			break
		}
	}

	for _, apple := range state.apples {
		if x == apple.posX && y == apple.posY {
			x, y = randomPlace(state)
			break
		}
	}

	return x, y
}

func createApple(state gameState) (a apple) {
	a.posX, a.posY = randomPlace(state)

	return a
}

func (a *apple) render(state gameState) {
	padX, padY, blockSize, _, _ := getRightSize(state)

	state.renderer.Copy(state.textures.apple,
		&sdl.Rect{X: 0, Y: 0, W: 8, H: 8},
		&sdl.Rect{
			X: int32(float64(a.posX)*blockSize + padX),
			Y: int32(float64(a.posY)*blockSize + padY),
			W: int32(blockSize),
			H: int32(blockSize),
		},
	)
}
