package entities

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	shipWidth  = 50
	shipHeight = 22
)

type Ship struct {
	X         float64
	Y         float64
	Speed     float64
	Direction float64
	Health    int
	Alive     bool
}

func NewShip(x, y float64) *Ship {
	return &Ship{
		X:         x,
		Y:         y,
		Speed:     80,
		Direction: 1,
		Health:    3,
		Alive:     true,
	}
}

func (s *Ship) Update(
	deltaTime time.Duration,
	minX float64,
	maxX float64,
) {
	if !s.Alive {
		return
	}

	s.X += s.Speed * s.Direction * deltaTime.Seconds()

	if s.X <= minX {
		s.X = minX
		s.Direction = 1
	}

	if s.X+shipWidth >= maxX {
		s.X = maxX - shipWidth
		s.Direction = -1
	}
}

func (s *Ship) Draw(screen *ebiten.Image) {
	if !s.Alive {
		return
	}

	// Hull
	ebitenutil.DrawRect(
		screen,
		s.X,
		s.Y,
		shipWidth,
		shipHeight,
		color.RGBA{R: 90, G: 55, B: 35, A: 255},
	)

	// Deck
	ebitenutil.DrawRect(
		screen,
		s.X+10,
		s.Y-8,
		30,
		8,
		color.RGBA{R: 145, G: 100, B: 55, A: 255},
	)

	// Mast
	ebitenutil.DrawRect(
		screen,
		s.X+24,
		s.Y-35,
		3,
		27,
		color.RGBA{R: 70, G: 45, B: 25, A: 255},
	)

	// Sail
	ebitenutil.DrawRect(
		screen,
		s.X+27,
		s.Y-32,
		17,
		20,
		color.RGBA{R: 220, G: 215, B: 180, A: 255},
	)
}
