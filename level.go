package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Sprite string
}

type Level interface {
	GetTiles() []Tile
	GetSprites() map[string]*ebiten.Image
	UpdateLevel()
	SetTiles([]Tile)
}

type Level1 struct {
	tiles   []Tile
	sprites map[string]*ebiten.Image
}

func DrawLevel(screen *ebiten.Image, level Level) {
	tiles := level.GetTiles()
	sprites := level.GetSprites()

	for _, tile := range tiles {

		tileImage := sprites[tile.Sprite]
		if tileImage == nil {
			fmt.Println("Tile image not found for sprite:", tile.Sprite)
			continue
		}

		gridSize := 100.0
		drawImage(screen, tileImage, Transform{
			x:        tile.X * float64(gridSize),
			y:        tile.Y * float64(gridSize),
			width:    tile.Width * float64(gridSize),
			height:   tile.Height * float64(gridSize),
			rotation: 0,
		})
	}

}

func (level *Level1) GetTiles() []Tile {
	return level.tiles
}

func (level *Level1) GetSprites() map[string]*ebiten.Image {
	return level.sprites
}

func (level *Level1) UpdateLevel() {}

func (level *Level1) SetTiles(tiles []Tile) {
	level.tiles = tiles
}
