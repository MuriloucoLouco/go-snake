package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/veandco/go-sdl2/sdl"
)

type configFile struct {
	GridWidth    int
	GridHeight   int
	Speed        int
	ScreenWidth  int32
	ScreenHeight int32
	SnakeTexture string
	AppleTexture string
}

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
	config   configFile
}

func main() {
	//load configs
	var state gameState
	state.config = configFile{
        20,
        20,
        125,
        400,
        400,
        "sprites/snakes/snake.bmp",
        "sprites/fruits/apple.bmp",
    }
	cfgBinary, err := ioutil.ReadFile("./config.toml")
	if err != nil {
		fmt.Println("file config.toml couldn't be loaded")
	}

	configToml := string(cfgBinary)
	_, err = toml.Decode(configToml, &state.config)
	if err != nil {
		fmt.Println("file config.toml couldn't be parsed")
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	//init sdl
	window, err := sdl.CreateWindow(
		"snake",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		state.config.ScreenWidth, state.config.ScreenHeight,
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

	//define state
	state.window = window
	state.renderer = renderer

	state.textures.snake = loadTextureFromBMP(state.config.SnakeTexture, state.renderer)
	state.textures.apple = loadTextureFromBMP(state.config.AppleTexture, state.renderer)

	s := createSnake(state)
	state.snake = &s
	a := createApple(state)
	state.apple = &a

	//main loop
	go render(state)

	for {
		sdl.WaitEvent()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		state.snake.move()
	}
}

func render(state gameState) {
	for {
		state.snake.update(state)

		padX, padY, _, screenWidth, screenHeight := getRightSize(state)
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
		time.Sleep(time.Duration(state.config.Speed) * time.Millisecond)
	}
}
