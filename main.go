package main

import (
	"fmt"
	"io/ioutil"
	"time"
    "path"

	"github.com/BurntSushi/toml"
	"github.com/veandco/go-sdl2/sdl"
)

type configFile struct {
	GridWidth     int
	GridHeight    int
	Speed         int
	ScreenWidth   int32
	ScreenHeight  int32
	SnakeTextures string
	AppleTextures string
    SnakeFile     string
    AppleFile     string
	FontTexture   string
}

type textureState struct {
	snake *sdl.Texture
	apple *sdl.Texture
	font  *sdl.Texture
}

type gameState struct {
	window    *sdl.Window
	renderer  *sdl.Renderer
	snake     *snake
	apple     *apple
	menu      *menu
	textures  textureState
	config    configFile
	paused    bool
	exited    bool
}

func main() {
	//load configs
	var state gameState
	state.paused = true
	state.exited = false

	state.config = configFile{
		20,
		20,
		125,
		400,
		400,
		"sprites/snakes/",
		"sprites/fruits/",
        "snake.bmp",
        "apple.bmp",
		"sprites/font.bmp",
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
    
    snakeFiles, err := ioutil.ReadDir(state.config.SnakeTextures)
    if err != nil || len(snakeFiles) == 0 {
        panic(err)
    }
    state.config.SnakeFile = snakeFiles[0].Name()
    
    appleFiles, err := ioutil.ReadDir(state.config.AppleTextures)
    if err != nil || len(appleFiles) == 0 {
        panic(err)
    }
    state.config.AppleFile = appleFiles[0].Name()

	state.textures.snake = loadTextureFromBMP(
        path.Join(state.config.SnakeTextures, snakeFiles[0].Name()),
        state.renderer,
    )
	state.textures.apple = loadTextureFromBMP(
        path.Join(state.config.AppleTextures, appleFiles[0].Name()),
        state.renderer,
    )
	state.textures.font = loadTextureFromBMP(state.config.FontTexture, state.renderer)

	s := createSnake(state)
	state.snake = &s
	a := createApple(state)
	state.apple = &a
	m := createMenu(state)
	state.menu = &m

	//main loop
	go render(&state)

	for {
		sdl.WaitEvent()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				state.exited = true
			}
		}
		if state.exited {
			return
		}
		state.menu.control(&state)
		if !state.paused {
			state.snake.control()
		}
	}
}

func render(state *gameState) {
	for {
		padX, padY, _, screenWidth, screenHeight := getRightSize(*state)
		state.renderer.SetDrawColor(0, 0, 0, 255)
		state.renderer.Clear()
		state.renderer.SetDrawColor(255, 255, 255, 255)
		state.renderer.FillRect(&sdl.Rect{
            X: int32(padX),
            Y: int32(padY),
            W: int32(screenWidth),
            H: int32(screenHeight),
        })
		state.renderer.SetDrawColor(0, 0, 0, 255)
		state.renderer.FillRect(&sdl.Rect{
            X: int32(padX) + 1,
            Y: int32(padY) + 1,
            W: int32(screenWidth) - 2,
            H: int32(screenHeight) - 2,
        })

		if !state.paused {
			state.snake.update(*state)
		}
		state.snake.render(*state)
		state.apple.render(*state)
		if state.paused {
			state.menu.render(*state)
		}

		state.renderer.Present()

		state.window.SetTitle("Snake - score: " + fmt.Sprintf("%d", state.snake.score))
		time.Sleep(time.Duration(state.config.Speed) * time.Millisecond)
	}
}
