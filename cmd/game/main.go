package main

import (
	"log"

	gamepkg "github.com/fburtin/total-medieval-destruction/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	const (
		windowWidth  = 960
		windowHeight = 540
	)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Total Medieval Destruction")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	currentGame, err := gamepkg.New()

	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(currentGame); err != nil {
		log.Fatal(err)
	}
}
