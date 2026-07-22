package assets

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const grassSpriteCount = 8

type Sprites struct {
	Grass []*ebiten.Image
}

func Load() (*Sprites, error) {
	grassSprites, err := loadNumberedSprites(
		"assets/sprites/grass/grass_%d.png",
		grassSpriteCount,
	)
	if err != nil {
		return nil, fmt.Errorf("load grass sprites: %w", err)
	}

	return &Sprites{
		Grass: grassSprites,
	}, nil
}

func loadNumberedSprites(pattern string, count int) ([]*ebiten.Image, error) {
	sprites := make([]*ebiten.Image, 0, count)

	for i := 1; i <= count; i++ {
		path := fmt.Sprintf(pattern, i)

		sprite, err := loadImage(path)
		if err != nil {
			return nil, err
		}

		sprites = append(sprites, sprite)
	}

	return sprites, nil
}

func loadImage(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image %q: %w", path, err)
	}
	defer file.Close()

	decodedImage, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode image %q: %w", path, err)
	}

	return ebiten.NewImageFromImage(decodedImage), nil
}
