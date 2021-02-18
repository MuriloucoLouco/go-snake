package main

import (
	"fmt"
    "path"
    "io/ioutil"

	"github.com/veandco/go-sdl2/sdl"
)

type fn func(*gameState)

type menu struct {
	options  []option
	selected int
}

type option struct {
	text   string
	method fn
}

func unpause(state *gameState) {
    keys := sdl.GetKeyboardState()
    if keys[sdl.SCANCODE_RETURN] == 1 {
        state.paused = false
    }
}

func changeGridWidth(state *gameState) {
    newValue := state.config.GridWidth
    keys := sdl.GetKeyboardState()

    if keys[sdl.SCANCODE_LEFT] == 1 {
        newValue--
    } else if keys[sdl.SCANCODE_RIGHT] == 1 {
        newValue++
    }
    if newValue >= 5 && newValue <= 80 {
        state.config.GridWidth = newValue
        state.menu.options[state.menu.selected].text = fmt.Sprintf("Grid Width %d", state.config.GridWidth)
    }
}

func changeGridHeight(state *gameState) {
    newValue := state.config.GridHeight
    keys := sdl.GetKeyboardState()

    if keys[sdl.SCANCODE_LEFT] == 1 {
        newValue--
    } else if keys[sdl.SCANCODE_RIGHT] == 1 {
        newValue++
    }
    if newValue >= 5 && newValue <= 80 {
        state.config.GridHeight = newValue
        state.menu.options[state.menu.selected].text = fmt.Sprintf("Grid Height %d", state.config.GridHeight)
    }
}

func changeSkin(state *gameState) {
    keys := sdl.GetKeyboardState()
    if !(keys[sdl.SCANCODE_LEFT] == 1 || keys[sdl.SCANCODE_RIGHT] == 1) {
        return
    }
    files, err := ioutil.ReadDir(state.config.SnakeTextures)
    if err != nil {
        fmt.Println("couldn't load "+ state.config.SnakeTextures +" folder")
        return
    }
    if len(files) == 0 {
        fmt.Println("folder "+ state.config.SnakeTextures +" has no files")
        return
    }
    
    var index int
    
    for i, file := range files {
        if file.Name() == state.config.SnakeFile {
            index = i
        }
    }
    
    if keys[sdl.SCANCODE_RIGHT] == 1 {
        if index < len(files)-1 {
            index++
        } else {
            index = 0
        }
    } else if keys[sdl.SCANCODE_LEFT] == 1 {
        if index > 0 {
            index--
        } else {
            index = len(files)-1
        }
    }
    state.config.SnakeFile = files[index].Name()
    state.textures.snake = loadTextureFromBMP(
        path.Join(state.config.SnakeTextures, files[index].Name()),
        state.renderer,
    )
}

func changeFruit(state *gameState) {
    keys := sdl.GetKeyboardState()
    if !(keys[sdl.SCANCODE_LEFT] == 1 || keys[sdl.SCANCODE_RIGHT] == 1) {
        return
    }
    files, err := ioutil.ReadDir(state.config.AppleTextures)
    if err != nil {
        fmt.Println("couldn't load "+ state.config.AppleTextures +" folder")
        return
    }
    if len(files) == 0 {
        fmt.Println("folder "+ state.config.AppleTextures +" has no files")
        return
    }
    
    var index int
    
    for i, file := range files {
        if file.Name() == state.config.AppleFile {
            index = i
        }
    }
    
    if keys[sdl.SCANCODE_RIGHT] == 1 {
        if index < len(files)-1 {
            index++
        } else {
            index = 0
        }
    } else if keys[sdl.SCANCODE_LEFT] == 1 {
        if index > 0 {
            index--
        } else {
            index = len(files)-1
        }
    }
    state.config.AppleFile = files[index].Name()
    state.textures.apple = loadTextureFromBMP(
        path.Join(state.config.AppleTextures, files[index].Name()),
        state.renderer,
    )
}

func changeAppleNumber(state *gameState) {
    newValue := state.config.AppleNumber
    keys := sdl.GetKeyboardState()

    if keys[sdl.SCANCODE_LEFT] == 1 {
        newValue--
    } else if keys[sdl.SCANCODE_RIGHT] == 1 {
        newValue++
    }
    if newValue >= 1 && newValue <= state.config.GridWidth * state.config.GridHeight {
        state.config.AppleNumber = newValue
        state.menu.options[state.menu.selected].text = fmt.Sprintf("Fruit number %d", state.config.AppleNumber)
    }
}

func exitGame(state *gameState) {
    keys := sdl.GetKeyboardState()
    if keys[sdl.SCANCODE_RETURN] == 1 {
        state.exited = true
    }
}

func createMenu(state gameState) (m menu) {

	m.options = []option{
		{
			"Continue",
			unpause,
        },
		{
			fmt.Sprintf("Grid Width %d", state.config.GridWidth),
			changeGridWidth,
		},
		{
			fmt.Sprintf("Grid Height %d", state.config.GridHeight),
			changeGridHeight,
		},
        {
            "Change skin",
            changeSkin,
        },
        {
            "Change fruit",
            changeFruit,
        },
        {
            fmt.Sprintf("Fruit number %d", state.config.AppleNumber),
            changeAppleNumber,
        },
        {
			"Exit Game",
			exitGame,
		},
	}
	m.selected = 0
	return m
}

func (m *menu) control(state *gameState) {
	keys := sdl.GetKeyboardState()

	if keys[sdl.SCANCODE_ESCAPE] == 1 {
		state.paused = !state.paused
	}

	if !state.paused {
		return
	}

	if keys[sdl.SCANCODE_DOWN] == 1 {
		if m.selected < len(m.options)-1 {
			m.selected++
		} else {
            m.selected = 0
        }
	} else if keys[sdl.SCANCODE_UP] == 1 {
		if m.selected > 0 {
			m.selected--
		} else {
            m.selected = len(m.options)-1
        }
	}

	m.options[m.selected].method(state)
}

func (m *menu) render(state gameState) {
	padX, padY, blockSize, screenWidth, screenHeight := getRightSize(state)
    var blockX, blockY float64
    blockY = float64(len(m.options)) * 1.5
    for _, option := range m.options {
        if float64(len(option.text)) > blockX {
            blockX = float64(len(option.text))
        }
    }
    
    blockSize = screenHeight / blockY
    
    
    if blockSize * blockX > screenWidth {
        blockSize = screenWidth / blockX
    }
    
	state.renderer.SetDrawColor(255, 255, 255, 255)
	state.renderer.DrawRect(&sdl.Rect{
		X: int32(padX) - 2,
		Y: int32(padY + blockSize*float64(m.selected)*1.5),
		W: int32(blockSize)*int32(len(m.options[m.selected].text)) + 4,
		H: int32(blockSize) + 4,
	})

	for i, option := range m.options {
		writeText(option.text,
			int32(padX+2),
			int32(2+padY+blockSize*float64(i)*1.5),
			int32(blockSize),
			state,
		)
	}
}
