package game

import (
	"time"

	"github.com/fburtin/total-medieval-destruction/internal/world"
)

func (g *Game) restart() {
	g.grid = world.NewGrid(gridCols, gridRows)

	g.hoverColumn = -1
	g.hoverRow = -1
	g.hoverValid = false

	g.roundNumber = 1
	g.phaseTimeLeft = 0
	g.lastUpdate = time.Now()

	g.enemyShips = nil

	g.initializeLevel()
	g.startPhase(PhaseBuild)
}
