package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type snake struct {
	positions     [][4]int
	direction     string
	nextDirection string
	growing       bool
	score         int
}

func createSnake(state gameState) (s snake) {
	gridWidth := state.config.GridWidth
	gridHeight := state.config.GridHeight
	centerX := int(gridWidth / 2)
	centerY := int(gridHeight / 2)

	s.positions = [][4]int{
		{centerX - 2, centerY, 90, 0},
		{centerX - 1, centerY, 90, 0},
		{centerX - 0, centerY, 90, 0},
	}
	s.direction = "right"
	s.nextDirection = "right"
	s.growing = false
	s.score = 0

	return s
}

func (s *snake) control() {
	newDirection := s.nextDirection
	keys := sdl.GetKeyboardState()

	switch uint8(1) {
	case keys[sdl.SCANCODE_RIGHT], keys[sdl.SCANCODE_D]:
		if s.direction != "left" {
			newDirection = "right"
		}
	case keys[sdl.SCANCODE_LEFT], keys[sdl.SCANCODE_A]:
		if s.direction != "right" {
			newDirection = "left"
		}
	case keys[sdl.SCANCODE_UP], keys[sdl.SCANCODE_W]:
		if s.direction != "down" {
			newDirection = "up"
		}
	case keys[sdl.SCANCODE_DOWN], keys[sdl.SCANCODE_S]:
		if s.direction != "up" {
			newDirection = "down"
		}
	}

	s.nextDirection = newDirection
}

func (s *snake) eat(state gameState) {
	snakeHead := s.positions[len(s.positions)-1]
	if snakeHead[0] == state.apple.posX && snakeHead[1] == state.apple.posY {
		s.growing = true
		state.apple.posX, state.apple.posY = randomPlace(
			state.snake.positions,
			state.config.GridWidth,
			state.config.GridHeight,
		)
		s.score++
	}
}

func (s *snake) die(state gameState) {
	*state.snake = createSnake(state)
	*state.apple = createApple(state)
}

func (s *snake) update(state gameState) {
	gridWidth := state.config.GridWidth
	gridHeight := state.config.GridHeight

	if !s.growing {
		s.positions = s.positions[1:]
	} else {
		s.growing = false
	}

	if s.nextDirection != "" {
		s.direction = s.nextDirection
	}

	lastPosition := s.positions[len(s.positions)-1]
	var newPosition [4]int
	switch s.direction {
	case "right":
		newPosition = [4]int{lastPosition[0] + 1, lastPosition[1], 90, 0}
	case "left":
		newPosition = [4]int{lastPosition[0] - 1, lastPosition[1], 270, 0}
	case "up":
		newPosition = [4]int{lastPosition[0], lastPosition[1] - 1, 0, 0}
	case "down":
		newPosition = [4]int{lastPosition[0], lastPosition[1] + 1, 180, 0}
	}

	if newPosition[0] < 0 || newPosition[0] >= gridWidth ||
     newPosition[1] < 0 || newPosition[1] >= gridHeight {
		s.die(state)
		return
	}

	for _, position := range s.positions {
		if position[0] == newPosition[0] && position[1] == newPosition[1] {
			s.die(state)
			return
		}
	}

	s.positions = append(s.positions, newPosition)

	//absolute distance between last and antepenultimate positions in both X and Y directions
	thirdPosition := s.positions[len(s.positions)-3]
	cornerDistanceX := float64(newPosition[0] - thirdPosition[0])
	cornerDistanceY := float64(newPosition[1] - thirdPosition[1])
	if math.Abs(cornerDistanceX) == 1 && math.Abs(cornerDistanceY) == 1 {
		secondPosition := s.positions[len(s.positions)-2]
		var cornerAngle int

		if cornerDistanceX == 1 && cornerDistanceY == 1 {
			if newPosition[1]-secondPosition[1] == 0 {
				cornerAngle = 0
			} else {
				cornerAngle = 180
			}
		}
		if cornerDistanceX == -1 && cornerDistanceY == 1 {
			if newPosition[1]-secondPosition[1] == 0 {
				cornerAngle = 270
			} else {
				cornerAngle = 90
			}
		}
		if cornerDistanceX == 1 && cornerDistanceY == -1 {
			if newPosition[1]-secondPosition[1] == 0 {
				cornerAngle = 90
			} else {
				cornerAngle = 270
			}
		}
		if cornerDistanceX == -1 && cornerDistanceY == -1 {
			if newPosition[1]-secondPosition[1] == 0 {
				cornerAngle = 180
			} else {
				cornerAngle = 0
			}
		}

		s.positions[len(s.positions)-2][2] = cornerAngle
		s.positions[len(s.positions)-2][3] = 1
	}

	s.eat(state)
}

func (s *snake) render(state gameState) {
	padX, padY, blockSize, _, screenHeight := getRightSize(state)

	for i, position := range s.positions {
		var textureCoord int32
		var adjust int32
		adjust = int32(screenHeight / 500)

		if position[3] == 0 {
			textureCoord = 0
		} else {
			textureCoord = 8
		}

		if i == len(s.positions)-1 {
			textureCoord = 16
		}

		var adjustX, adjustY, stretchX, stretchY int32
		if screenHeight > 500 && position[3] == 1 {
			switch position[2] {
			case 0:
				adjustY = -2
				if i+1 <= len(s.positions)-1 && s.positions[i+1][0] > position[0] {
					adjustX = -2
					stretchX = 1
					stretchY = 3
				}
			case 90:
				adjustX = 2
				if i+1 <= len(s.positions)-1 && s.positions[i+1][1] > position[1] {
					adjustY = -1
					stretchX = 0
					stretchY = 3
				}
			case 180:
				adjustY = 2
				if i+1 <= len(s.positions)-1 && s.positions[i+1][0] < position[0] {
					adjustX = 2
					stretchY = 3
					stretchX = 1
				}
			case 270:
				adjustX = -2
				if i+1 <= len(s.positions)-1 && s.positions[i+1][1] < position[1] {
					adjustY = 3
					stretchX = 3
					stretchY = 3
				}
			}
		}

		state.renderer.CopyEx(state.textures.snake,
			&sdl.Rect{
				X: textureCoord,
				Y: 0,
				W: 8,
				H: 8,
			},
			&sdl.Rect{
				X: int32(float64(position[0])*blockSize+padX) - int32(adjust) + adjustX,
				Y: int32(float64(position[1])*blockSize+padY) - int32(adjust) + adjustY,
				W: int32(blockSize) + int32(adjust) + stretchX,
				H: int32(blockSize) + int32(adjust) + stretchY,
			},
			float64(position[2]),
			&sdl.Point{
				X: int32(blockSize / 2),
				Y: int32(blockSize / 2),
			},
			sdl.FLIP_NONE,
		)
	}
}
