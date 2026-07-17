package game

import (
	"github.com/fburtin/total-medieval-destruction/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) updateMousePosition() {
	mouseX, mouseY := ebiten.CursorPosition()

	column := (mouseX - gridOffsetX) / tileSize
	row := (mouseY - gridOffsetY) / tileSize

	g.hoverColumn = column
	g.hoverRow = row
	g.hoverValid = g.grid.IsInside(column, row)
}

func (g *Game) processInput() {
	if !g.hoverValid {
		return
	}

	if !g.canBuild() {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.grid.ToggleWall(
			g.hoverColumn,
			g.hoverRow,
		)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.grid.ClearTile(
			g.hoverColumn,
			g.hoverRow,
		)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if g.grid.TileAt(
			g.hoverColumn,
			g.hoverRow,
		) == world.TileEmpty {
			g.grid.SetCastle(
				g.hoverColumn,
				g.hoverRow,
			)
		}
	}
}
