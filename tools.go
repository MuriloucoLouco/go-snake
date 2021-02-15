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

	key := sdl.MapRGB(img.Format, 0, 0, 0)

	if err := img.SetColorKey(true, key); err != nil {

	}

	texture, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture: %v", err))
	}
	return texture
}

func getRightSize(w *sdl.Window) (padX, padY, blockSize, screenX, screenY float64) {
	screenWidth, screenHeight := w.GetSize()
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
