package main

import (
	"fmt"
	"time"

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
	GetColliders() []Point
	UpdateLevel()
	SetTiles([]Tile)
	SetColliders([]Point)
	StartLevel()
}

type Level1 struct {
	tiles            []Tile
	colliders        []Point
	sprites          map[string]*ebiten.Image
	dynamicColliders []*Collider

	train *Train
}

func (level *Level1) StartLevel() {
	level.dynamicColliders = []*Collider{
		{
			transform: Transform{
				x:      260,
				y:      0,
				width:  100,
				height: 1200,
			},
		},
	}

	level.train = &Train{
		transform: Transform{
			x:      240,
			y:      -1000,
			width:  100,
			height: 100,
		},
		parts: map[string]*ebiten.Image{
			"front":  loadImage("assets/sprites/u2-front.png"),
			"middle": loadImage("assets/sprites/u2-middle.png"),
		},
		length: 3,
	}

	gameObjects = append(gameObjects, level.train)

	go func() {
		time.Sleep(time.Second * 2)
		level.train.Drive(1500, 0.2)
	}()

	for _, collider := range level.dynamicColliders {
		gameObjects = append(gameObjects, collider)
	}
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
func (level *Level1) GetColliders() []Point {
	return level.colliders
}

func (level *Level1) SetColliders(colliders []Point) {
	level.colliders = colliders
}
