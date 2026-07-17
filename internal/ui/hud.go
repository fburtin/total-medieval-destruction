package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type HUDData struct {
	HoverColumn int
	HoverRow    int
	HoverValid  bool

	CastleExists bool
	CastleColumn int
	CastleRow    int

	FortressEnclosed bool

	Phase       string
	TimeLeft    float64
	RoundNumber int
	CanBuild    bool
}

func DrawHUD(screen *ebiten.Image, data HUDData) {
	mousePosition := "Outside grid"

	if data.HoverValid {
		mousePosition = fmt.Sprintf(
			"%d, %d",
			data.HoverColumn,
			data.HoverRow,
		)
	}

	castleStatus := "Castle: not placed"
	fortressStatus := "Fortress: unavailable"

	if data.CastleExists {
		castleStatus = fmt.Sprintf(
			"Castle: %d, %d",
			data.CastleColumn,
			data.CastleRow,
		)

		if data.FortressEnclosed {
			fortressStatus = "Fortress: ENCLOSED"
		} else {
			fortressStatus = "Fortress: OPEN"
		}
	}

	message := fmt.Sprintf(
		"TOTAL MEDIEVAL DESTRUCTION\n"+
			"Phase: %s\n"+
			"Time: %.1f\n"+
			"Round: %d\n"+
			"Building enabled: %t\n"+
			"Left click: place/remove wall\n"+
			"Right click: clear tile\n"+
			"C: place castle\n"+
			"Mouse tile: %s\n"+
			"%s\n"+
			"%s",
		data.Phase,
		data.TimeLeft,
		data.RoundNumber,
		data.CanBuild,
		mousePosition,
		castleStatus,
		fortressStatus,
	)

	ebitenutil.DebugPrint(screen, message)
}
