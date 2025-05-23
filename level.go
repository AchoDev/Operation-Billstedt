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
		y:      -3000,
		width:  100,
		height: 100,
	}, 1)

	train2 := CreateTrain(Transform{
		x: 1250,
		y: 3000,
	}, -1)

	addGameObject(train)
	addGameObject(train2)

	go func() {
		pausableSleep(time.Second * time.Duration(5 + rand.IntN(5)))
		train.Drive(4000, 0.2)
	}()

	go func() {
		pausableSleep(time.Second * time.Duration(5 + rand.IntN(5)))
		train2.Drive(4000, 0.2)
	}()

	for _, collider := range level.dynamicColliders {
		addGameObject(collider)
	}
}

func DrawLevel(screen *ebiten.Image, level Level) {
	tiles := level.GetTiles()
	sprites := level.GetSprites()
	gridSize := 100.0



	// Create a map to group tiles by their Z order
    itemsByZ := make(map[float64][]interface{})
	for _, tile := range tiles {
		itemsByZ[tile.Z] = append(itemsByZ[tile.Z], tile)
	}
	for _, gameObject := range gameObjects {
        z := gameObject.GetTransform().z
        itemsByZ[z] = append(itemsByZ[z], gameObject)
    }

	// Extract Z orders and sort them
	var zOrders []float64
	for z := range itemsByZ {
		zOrders = append(zOrders, z)
	}
	sort.Float64s(zOrders)

	// Iterate through tiles in sorted Z order
	for _, z := range zOrders {
		for _, item := range itemsByZ[z] {

			switch v := item.(type) {
			case Tile:
				tileImage := sprites[v.Sprite]
				if tileImage == nil {
					fmt.Println("Tile image not found for sprite:", v.Sprite)
					continue
				}
	
				op := defaultImageOptions()
				if levelEditorActivated && v.Z != float64(currentZEditor) {
					op.Alpha = 50
	
					if !onionSkin {
						continue
					}
				}
	
				drawImageWithOptions(screen, tileImage, Transform{
					x:        v.X * float64(gridSize),
					y:        v.Y * float64(gridSize),
					width:    v.Width * float64(gridSize),
					height:   v.Height * float64(gridSize),
					rotation: v.Rotation,
				}, op)
			case GameObject:
				if levelEditorActivated && hideGameobjects {
					continue
				}
				v.Draw(screen)
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
