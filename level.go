package main

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	X      float64
	Y      float64
	Z      float64
	Width  float64
	Height float64
	Rotation float64
	Sprite string
}

type Level interface {
	GetTiles() []Tile
	GetSprites() map[string]*ebiten.Image
	GetColliders() []Tile
	UpdateLevel()
	SetTiles([]Tile)
	SetColliders([]Tile)
	StartLevel()
}

type Level1 struct {
	tiles            []Tile
	colliders        []Tile
	sprites          map[string]*ebiten.Image
	dynamicColliders []*Collider

	// phase int
}

func (level *Level1) StartLevel() {
	train := CreateTrain(Transform{
		x:      240,
		y:      -1000,
		width:  100,
		height: 100,
	}, 1)

	train2 := CreateTrain(Transform{
		x: 1250,
		y: 1500,
	}, -1)

	gameObjects = append(gameObjects, train)
	gameObjects = append(gameObjects, train2)

	go func() {
		pausableSleep(time.Second * time.Duration(5 + rand.IntN(5)))
		train.Drive(2000, 0.2)
	}()

	go func() {
		pausableSleep(time.Second * time.Duration(5 + rand.IntN(5)))
		train2.Drive(2500, 0.2)
	}()

	for _, collider := range level.dynamicColliders {
		gameObjects = append(gameObjects, collider)
	}
}

func DrawLevel(screen *ebiten.Image, level Level) {
	tiles := level.GetTiles()
	sprites := level.GetSprites()
	gridSize := 100.0

	// Create a map to group tiles by their Z order
	tilesByZ := make(map[float64][]Tile)
	for _, tile := range tiles {
		tilesByZ[tile.Z] = append(tilesByZ[tile.Z], tile)
	}

	// Extract Z orders and sort them
	var zOrders []float64
	for z := range tilesByZ {
		zOrders = append(zOrders, z)
	}
	sort.Float64s(zOrders)

	// Iterate through tiles in sorted Z order
	for _, z := range zOrders {
		for _, tile := range tilesByZ[z] {
			tileImage := sprites[tile.Sprite]
			if tileImage == nil {
				fmt.Println("Tile image not found for sprite:", tile.Sprite)
				continue
			}

			op := defaultImageOptions()
			if levelEditorActivated && tile.Z != float64(currentZEditor) {
				op.Alpha = 50

				if !onionSkin {
					continue
				}
			}

			drawImageWithOptions(screen, tileImage, Transform{
				x:        tile.X * float64(gridSize),
				y:        tile.Y * float64(gridSize),
				width:    tile.Width * float64(gridSize),
				height:   tile.Height * float64(gridSize),
				rotation: tile.Rotation,
			}, op)
		}

		for _, gameObject := range gameObjects {
			
			if levelEditorActivated && hideGameobjects {
				continue
			}

			if gameObject.GetTransform().z == z {
				gameObject.Draw(screen)
			}
		}
	}

	if levelEditorActivated && selectedTool == 1{
		for _, collider := range level.GetColliders() {

			if collider.Z != float64(currentZEditor){
				continue
			}

			drawRect(screen, Transform{
				x:      float64(collider.X) * gridSize,
				y:      float64(collider.Y) * gridSize,
				width:  collider.Width * gridSize,
				height: collider.Height * gridSize,
			}, color.RGBA{255, 100, 200, 50})
		}
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
func (level *Level1) GetColliders() []Tile {
	return level.colliders
}

func (level *Level1) SetColliders(colliders []Tile) {
	level.colliders = colliders
}
