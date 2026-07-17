package game

import (
	"image/color"

	"github.com/fburtin/total-medieval-destruction/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// drawTileSprites adds pixel-art details over the existing tile colours.
// The pattern is deterministic, so tiles do not flicker between frames.
func (g *Game) drawTileSprites(screen *ebiten.Image) {
	