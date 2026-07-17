package ui

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	hudFontSize   = 28
	hudLineHeight = 38
)

var hudFontSource *text.GoTextFaceSource

func init() {
	var err error

	hudFontSource, err = text.NewGoTextFaceSource(
		bytes.NewReader(goregular.TTF),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load HUD font: %v", err))
	}
}

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
	lines := []string{
		fmt.Sprintf("ROUND: %d", data.RoundNumber),
		fmt.Sprintf("PHASE: %s", data.Phase),
		fmt.Sprintf("TIME: %.1f", data.TimeLeft),
		fmt.Sprintf("CAN BUILD: %t", data.CanBuild),
		fmt.Sprintf("FORTRESS ENCLOSED: %t", data.FortressEnclosed),
	}

	if data.HoverValid {
		lines = append(
			lines,
			fmt.Sprintf(
				"TILE: column %d, row %d",
				data.HoverColumn,
				data.HoverRow,
			),
		)
	}

	if data.CastleExists {
		lines = append(
			lines,
			fmt.Sprintf(
				"CASTLE: column %d, row %d",
				data.CastleColumn,
				data.CastleRow,
			),
		)
	}

	face := &text.GoTextFace{
		Source: hudFontSource,
		Size:   hudFontSize,
	}

	x := 30.0
	y := 20.0

	for index, line := range lines {
		options := &text.DrawOptions{}

		options.GeoM.Translate(
			x,
			y+float64(index*hudLineHeight),
		)

		options.ColorScale.ScaleWithColor(
			color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			},
		)

		text.Draw(
			screen,
			line,
			face,
			options,
		)
	}
}
