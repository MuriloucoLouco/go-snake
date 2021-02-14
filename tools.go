package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

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

func getRightSize(w *sdl.Window) (padX, padY, blockSize, screenX, screenY int32) {
	screenWidth, screenHeight := w.GetSize()
	proportion := float64(gridWidth) / float64(gridHeight)
	screenProportion := float64(screenWidth) / float64(screenHeight)

	if screenProportion > proportion {
		screenX = int32(float64(screenHeight) * proportion)
		screenY = screenHeight
		blockSize = screenY / gridHeight
		padX = (screenWidth - screenX) / 2
	} else {
		screenX = screenWidth
		screenY = int32(float64(screenWidth) / proportion)
		blockSize = screenX / gridWidth
		padY = (screenHeight - screenY) / 2
	}

	return padX, padY, blockSize, screenX, screenY
}
