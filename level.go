package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	X      int
	Y      int
	Width  int
	Height int
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

var levelEditorActivated bool = false

func DrawLevel(screen *ebiten.Image, level Level) {
	tiles := level.GetTiles()
	sprites := level.GetSprites()

	for _, tile := range tiles {

		tileImage := sprites[tile.Sprite]
		if tileImage == nil {
			fmt.Println("Tile image not found for sprite:", tile.Sprite)
			continue
		}

		gridSize := 100
		drawImage(screen, tileImage, Transform{
			x:        float64(tile.X * gridSize),
			y:        float64(tile.Y * gridSize),
			width:    float64(tile.Width * gridSize),
			height:   float64(tile.Height * gridSize),
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
