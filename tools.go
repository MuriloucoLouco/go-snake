package main

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func loadTextureFromBMP(spriteName string, renderer *sdl.Renderer) (texture *sdl.Texture) {

	img, err := sdl.LoadBMP(spriteName)
	if err != nil {
		panic(fmt.Errorf("loading sprite: %v", err))
	}
	defer img.Free()

	key := sdl.MapRGB(img.Format, 0, 0, 0)

	if err := img.SetColorKey(true, key); err != nil {

	}

	texture, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture: %v", err))
	}
	return texture
}

func getRightSize(state gameState) (padX, padY, blockSize, screenX, screenY float64) {
	gridWidth := float64(state.config.GridWidth)
	gridHeight := float64(state.config.GridHeight)
	screenWidth, screenHeight := state.window.GetSize()
	proportion := float64(gridWidth) / float64(gridHeight)
	screenProportion := float64(screenWidth) / float64(screenHeight)

	if screenProportion > proportion {
		screenX = float64(screenHeight) * proportion
		screenY = float64(screenHeight)
		blockSize = screenY / gridHeight
		padX = (float64(screenWidth) - screenX) / 2
	} else {
		screenX = float64(screenWidth)
		screenY = float64(screenWidth) / proportion
		blockSize = screenX / gridWidth
		padY = (float64(screenHeight) - screenY) / 2
	}

	return padX, padY, blockSize, screenX, screenY
}

func writeText(text string, x, y, size int32, state gameState) {
	for i, letter := range strings.ToLower(text) {
		offsetX := int32(0)
		offsetY := int32(0)
		fontSize := int32(8)
		if int(letter) >= 97 && int(letter) <= 122 {
			offsetX = fontSize * (int32(letter) - 96)
		}
		if int(letter) >= 48 && int(letter) <= 57 {
			offsetX = fontSize * (int32(letter) - 48)
			offsetY = fontSize
		}
		state.renderer.Copy(state.textures.font,
			&sdl.Rect{X: offsetX, Y: offsetY, W: 8, H: 8},
			&sdl.Rect{X: x + size*int32(i), Y: y, W: size, H: size},
		)
	}
}
