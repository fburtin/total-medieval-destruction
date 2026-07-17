package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 42,
		G: 72,
		B: 46,
		A: 255,
	})

	g.drawGrid(screen)

	for _, ship := range g.enemyShips {
		if ship != nil {
			ship.Draw(screen)
		}
	}

	g.drawHover(screen)
	g.drawRestartButton(screen)
}

func (g *Game) Layout(
	outsideWidth int,
	outsideHeight int,
) (int, int) {
	return screenWidth, screenHeight
}
