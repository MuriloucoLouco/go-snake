package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type snake struct {
	texture       *sdl.Texture
	positions     [][4]int
	direction     string
	nextDirection string
	growing       bool
	score         int
}

func createSnake(state gameState) (s snake) {
	s.texture = loadTextureFromBMP("sprites/snake.bmp", state.renderer)
	s.positions = [][4]int{
		{centerX - 3, centerY, 90, 0},
		{centerX - 2, centerY, 90, 0},
		{centerX - 1, centerY, 90, 0},
	}
	s.direction = "right"
	s.nextDirection = "right"
	s.growing = false
	s.score = 0

	return s
}

func (s *snake) move() {
	newDirection := s.nextDirection
	keys := sdl.GetKeyboardState()

	switch uint8(1) {
	case keys[sdl.SCANCODE_RIGHT]:
		if s.direction != "left" {
			newDirection = "right"
		}
	case keys[sdl.SCANCODE_LEFT]:
		if s.direction != "right" {
			newDirection = "left"
		}
	case keys[sdl.SCANCODE_UP]:
		if s.direction != "down" {
			newDirection = "up"
		}
	case keys[sdl.SCANCODE_DOWN]:
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
		state.apple.posX, state.apple.posY = randomPlace(s.positions)
		s.score++
	}
}

func (s *snake) die(state gameState) {
	*state.snake = createSnake(state)
	*state.apple = createApple(state)
}

func (s *snake) update(state gameState) {
	s.eat(state)

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

	if newPosition[0] < 0 || newPosition[0] >= gridWidth || newPosition[1] < 0 || newPosition[1] >= gridHeight {
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
}

func (s *snake) render(state gameState) {
	padX, padY, blockSize, _, _ := getRightSize(state.window)

	for i, position := range s.positions {
		var textureCoord int32

		if position[3] == 0 {
			textureCoord = 0
		} else {
			textureCoord = 8
		}

		if i == len(s.positions)-1 {
			textureCoord = 16
		}

		state.renderer.CopyEx(s.texture,
			&sdl.Rect{
				X: textureCoord,
				Y: 0,
				W: 8,
				H: 8,
			},
			&sdl.Rect{
				X: int32(float64(position[0])*blockSize + padX),
				Y: int32(float64(position[1])*blockSize+padY) - 1,
				W: int32(blockSize),
				H: int32(blockSize) + 1,
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
