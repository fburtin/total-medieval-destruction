package game

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const restartButtonFontSize = 24

var restartButtonFontSource *text.GoTextFaceSource

func init() {
	var err error

	restartButtonFontSource, err = text.NewGoTextFaceSource(
		bytes.NewReader(goregular.TTF),
	)
	if err != nil {
		panic(fmt.Sprintf(
			"failed to load restart button font: %v",
			err,
		))
	}
}

func (g *Game) updateRestartButton() {
	mouseX, mouseY := ebiten.CursorPosition()

	x, y, width, height := g.restartButtonBounds()

	g.restartButtonHovered =
		mouseX >= x &&
			mouseX < x+width &&
			mouseY >= y &&
			mouseY < y+height

	if g.restartButtonHovered &&
		inpututil.IsMouseButtonJustPressed(
			ebiten.MouseButtonLeft,
		) {
		g.restart()
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.restart()
	}
}

func (g *Game) drawRestartButton(screen *ebiten.Image) {
	x, y, width, height := g.restartButtonBounds()

	backgroundColor := color.RGBA{
		R: 80,
		G: 80,
		B: 80,
		A: 255,
	}

	if g.restartButtonHovered {
		backgroundColor = color.RGBA{
			R: 115,
			G: 115,
			B: 115,
			A: 255,
		}
	}

	buttonImage := ebiten.NewImage(width, height)
	buttonImage.Fill(backgroundColor)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(buttonImage, options)

	drawButtonBorder(
		screen,
		x,
		y,
		width,
		height,
	)

	face := &text.GoTextFace{
		Source: restartButtonFontSource,
		Size:   restartButtonFontSize,
	}

	label := "RESTART GAME"

	textWidth, textHeight := text.Measure(
		label,
		face,
		0,
	)

	textOptions := &text.DrawOptions{}

	textOptions.GeoM.Translate(
		float64(x)+(float64(width)-textWidth)/2,
		float64(y)+(float64(height)-textHeight)/2,
	)

	textOptions.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		label,
		face,
		textOptions,
	)
}

func (g *Game) restartButtonBounds() (
	x int,
	y int,
	width int,
	height int,
) {
	width = restartButtonWidth
	height = restartButtonHeight

	x = screenWidth -
		restartButtonWidth -
		restartButtonMargin

	y = restartButtonMargin

	return
}

func drawButtonBorder(
	screen *ebiten.Image,
	x int,
	y int,
	width int,
	height int,
) {
	borderColor := color.RGBA{
		R: 210,
		G: 210,
		B: 210,
		A: 255,
	}

	border := ebiten.NewImage(width, 2)
	border.Fill(borderColor)

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(border, options)

	options.GeoM.Reset()
	options.GeoM.Translate(
		float64(x),
		float64(y+height-2),
	)
	screen.DrawImage(border, options)

	verticalBorder := ebiten.NewImage(2, height)
	verticalBorder.Fill(borderColor)

	options.GeoM.Reset()
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(verticalBorder, options)

	options.GeoM.Reset()
	options.GeoM.Translate(
		float64(x+width-2),
		float64(y),
	)
	screen.DrawImage(verticalBorder, options)
}
