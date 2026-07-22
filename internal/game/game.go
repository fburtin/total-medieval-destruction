package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/fburtin/total-medieval-destruction/assets"
	"github.com/fburtin/total-medieval-destruction/internal/entities"
	"github.com/fburtin/total-medieval-destruction/internal/ui"
	"github.com/fburtin/total-medieval-destruction/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1920
	screenHeight = 1080

	tileSize = 40

	gridCols = 40
	gridRows = 24

	gridOffsetX = 80
	gridOffsetY = 30

	restartButtonWidth  = 220
	restartButtonHeight = 60
	restartButtonMargin = 25
)

type Game struct {
	grid *world.Grid

	sprites *assets.Sprites

	hoverColumn int
	hoverRow    int
	hoverValid  bool

	phase         Phase
	phaseTimeLeft time.Duration
	roundNumber   int
	lastUpdate    time.Time

	enemyShips []*entities.Ship

	restartButtonHovered bool
}

func (g *Game) Update() error {
	now := time.Now()
	deltaTime := now.Sub(g.lastUpdate)
	g.lastUpdate = now

	g.updateRestartButton()

	g.updateMousePosition()

	if !g.restartButtonHovered {
		g.processInput()
	}

	g.updateRound(deltaTime)

	if g.phase == PhaseBattle {
		for _, ship := range g.enemyShips {
			g.updateEnemyShip(ship, deltaTime)
		}
	}

	return nil
}

func New() (*Game, error) {
	sprites, err := assets.Load()
	if err != nil {
		return nil, fmt.Errorf("load game assets: %w", err)
	}

	game := &Game{
		grid:        world.NewGrid(gridCols, gridRows),
		sprites:     sprites,
		hoverColumn: -1,
		hoverRow:    -1,
		roundNumber: 1,
		lastUpdate:  time.Now(),
	}

	game.initializeLevel()
	game.startPhase(PhaseBuild)

	return game, nil
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for row := 0; row < g.grid.Rows(); row++ {
		for column := 0; column < g.grid.Columns(); column++ {
			x := gridOffsetX + column*tileSize
			y := gridOffsetY + row*tileSize

			tileColor := g.tileColor(
				g.grid.TileAt(column, row),
			)

			ebitenutil.DrawRect(
				screen,
				float64(x),
				float64(y),
				float64(tileSize-1),
				float64(tileSize-1),
				tileColor,
			)
		}
	}

}

func (g *Game) drawHover(screen *ebiten.Image) {
	if !g.hoverValid {
		return
	}

	x := gridOffsetX + g.hoverColumn*tileSize
	y := gridOffsetY + g.hoverRow*tileSize

	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		float64(tileSize-1),
		float64(tileSize-1),
		color.RGBA{
			R: 255,
			G: 255,
			B: 180,
			A: 110,
		},
	)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	castleColumn, castleRow, castleExists := g.grid.FindCastle()

	data := ui.HUDData{
		HoverColumn:  g.hoverColumn,
		HoverRow:     g.hoverRow,
		HoverValid:   g.hoverValid,
		CastleExists: castleExists,
		CastleColumn: castleColumn,
		CastleRow:    castleRow,

		FortressEnclosed: castleExists &&
			g.grid.IsCastleEnclosed(),

		Phase:       g.phase.String(),
		TimeLeft:    g.phaseTimeLeft.Seconds(),
		RoundNumber: g.roundNumber,
		CanBuild:    g.canBuild(),
	}

	ui.DrawHUD(screen, data)
}

func (g *Game) tileColor(tile world.TileType) color.Color {
	switch tile {
	case world.TileWall:
		return color.RGBA{
			R: 125,
			G: 105,
			B: 85,
			A: 255,
		}

	case world.TileCastle:
		return color.RGBA{
			R: 190,
			G: 160,
			B: 55,
			A: 255,
		}

	case world.TileWater:
		return color.RGBA{
			R: 45,
			G: 105,
			B: 145,
			A: 255,
		}

	case world.TileEmpty:
		return color.RGBA{
			R: 92,
			G: 125,
			B: 75,
			A: 255,
		}

	default:
		return color.RGBA{
			R: 92,
			G: 125,
			B: 75,
			A: 255,
		}
	}
}

func (g *Game) updateEnemyShip(
	ship *entities.Ship,
	deltaTime time.Duration,
) {
	if ship == nil || !ship.Alive {
		return
	}

	currentRow := int(
		(ship.Y - float64(gridOffsetY)) /
			float64(tileSize),
	)

	nextRow := currentRow + 1

	if nextRow >= g.grid.Rows() {
		nextRow = 0
	}

	waterColumn, found := g.findWaterColumn(nextRow)
	if !found {
		return
	}

	targetX := float64(
		gridOffsetX + waterColumn*tileSize,
	)

	targetY := float64(
		gridOffsetY + nextRow*tileSize,
	)

	moveSpeed := 60.0 * deltaTime.Seconds()

	ship.X = moveTowards(
		ship.X,
		targetX,
		moveSpeed,
	)

	ship.Y = moveTowards(
		ship.Y,
		targetY,
		moveSpeed,
	)
}

func (g *Game) findWaterColumn(row int) (int, bool) {
	waterColumns := make([]int, 0)

	for column := 0; column < g.grid.Columns(); column++ {
		if g.grid.TileAt(column, row) == world.TileWater {
			waterColumns = append(waterColumns, column)
		}
	}

	if len(waterColumns) == 0 {
		return 0, false
	}

	return waterColumns[len(waterColumns)/2], true
}

func moveTowards(current, target, maximumDelta float64) float64 {
	if current < target {
		current += maximumDelta

		if current > target {
			return target
		}

		return current
	}

	if current > target {
		current -= maximumDelta

		if current < target {
			return target
		}
	}

	return current
}

func (g *Game) drawGrassTile(
	screen *ebiten.Image,
	column int,
	row int,
) {
	variation := (column + row*gridCols) % len(g.sprites.Grass)
	sprite := g.sprites.Grass[variation]

	op := &ebiten.DrawImageOptions{}

	spriteWidth := sprite.Bounds().Dx()
	spriteHeight := sprite.Bounds().Dy()

	op.GeoM.Scale(
		float64(tileSize)/float64(spriteWidth),
		float64(tileSize)/float64(spriteHeight),
	)

	x := gridOffsetX + column*tileSize
	y := gridOffsetY + row*tileSize

	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(sprite, op)
}
