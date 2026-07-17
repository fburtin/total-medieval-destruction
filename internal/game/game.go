package game

import (
	"image/color"
	"time"

	"github.com/fburtin/total-medieval-destruction/internal/ui"
	"github.com/fburtin/total-medieval-destruction/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 960
	screenHeight = 540

	tileSize = 40
	gridCols = 20
	gridRows = 12

	gridOffsetX = 80
	gridOffsetY = 30
)

type Game struct {
	grid *world.Grid

	hoverColumn int
	hoverRow    int
	hoverValid  bool

	phase         Phase
	phaseTimeLeft time.Duration
	roundNumber   int
	lastUpdate    time.Time
}

func New() *Game {
	game := &Game{
		grid:        world.NewGrid(gridCols, gridRows),
		hoverColumn: -1,
		hoverRow:    -1,
		roundNumber: 1,
		lastUpdate:  time.Now(),
	}

	game.startPhase(PhaseBuild)

	return game
}

func (g *Game) Update() error {
	now := time.Now()
	deltaTime := now.Sub(g.lastUpdate)
	g.lastUpdate = now

	g.updateMousePosition()
	g.processInput()
	g.updateRound(deltaTime)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 42,
		G: 72,
		B: 46,
		A: 255,
	})

	g.drawGrid(screen)
	g.drawHover(screen)
	g.drawHUD(screen)
}

func (g *Game) Layout(
	outsideWidth int,
	outsideHeight int,
) (screenWidth int, screenHeight int) {
	return 960, 540
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

	case world.TileEmpty:
		fallthrough

	default:
		return color.RGBA{
			R: 92,
			G: 125,
			B: 75,
			A: 255,
		}
	}
}
