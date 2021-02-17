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

func createMenu(state gameState) (m menu) {

	m.options = []option{
		{
			"Continue",
			func(state *gameState) {
				keys := sdl.GetKeyboardState()
				if keys[sdl.SCANCODE_RETURN] == 1 {
					state.paused = false
				}
			},
		},
		{
			fmt.Sprintf("Grid Width %d", state.config.GridWidth),
			func(state *gameState) {
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
			},
		},
		{
			fmt.Sprintf("Grid Height %d", state.config.GridHeight),
			func(state *gameState) {
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
			},
		},
        {
            "Change skin",
            func(state *gameState) {
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
                state.textures.snake = loadTextureFromBMP(path.Join(state.config.SnakeTextures, files[index].Name()), state.renderer)
                
            },
        },
        {
            "Change fruit",
            func(state *gameState) {
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
                state.textures.apple = loadTextureFromBMP(path.Join(state.config.AppleTextures, files[index].Name()), state.renderer)
                
            },
        },
        {
			"Exit Game",
			func(state *gameState) {
				keys := sdl.GetKeyboardState()
				if keys[sdl.SCANCODE_RETURN] == 1 {
					state.exited = true
				}
			},
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
	padX, padY, blockSize, _, _ := getRightSize(state)

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
