package game

import (
	"math/rand"
	"time"

	"github.com/fburtin/total-medieval-destruction/internal/entities"
	"github.com/fburtin/total-medieval-destruction/internal/world"
)

const (
	riverMinWidth = 3
	riverMaxWidth = 4
)

func (g *Game) initializeLevel() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	g.generateRandomRiver(rng)
	g.spawnEnemyShip()
}

func (g *Game) generateRandomRiver(rng *rand.Rand) {
	riverCenter := g.grid.Columns() - 5
	riverWidth := riverMinWidth

	for row := 0; row < g.grid.Rows(); row++ {
		if row > 0 {
			riverCenter += rng.Intn(3) - 1
		}

		if rng.Intn(4) == 0 {
			riverWidth = riverMinWidth + rng.Intn(
				riverMaxWidth-riverMinWidth+1,
			)
		}

		halfWidth := riverWidth / 2

		minCenter := halfWidth + 1
		maxCenter := g.grid.Columns() - halfWidth - 2

		if riverCenter < minCenter {
			riverCenter = minCenter
		}

		if riverCenter > maxCenter {
			riverCenter = maxCenter
		}

		startColumn := riverCenter - halfWidth
		endColumn := startColumn + riverWidth

		for column := startColumn; column < endColumn; column++ {
			g.grid.SetTile(column, row, world.TileWater)
		}
	}
}

func (g *Game) spawnEnemyShip() {
	for row := 0; row < g.grid.Rows(); row++ {
		for column := 0; column < g.grid.Columns(); column++ {
			if g.grid.TileAt(column, row) != world.TileWater {
				continue
			}

			x := float64(gridOffsetX + column*tileSize)
			y := float64(gridOffsetY + row*tileSize)

			g.enemyShip = entities.NewShip(x, y)

			return
		}
	}
}
