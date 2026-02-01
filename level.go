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

type EnemySpawnPoint struct {
	X float64
	Y float64
	EnemyType int
}

type LevelData struct {
	Tiles []Tile
	Colliders []Tile
	Enemies []EnemySpawnPoint
	Sprites map[string]*ebiten.Image
	CameraFollow Vector2
	CameraZoom float64
}

type Level interface {
	GetData() *LevelData
	UpdateLevel()
	StartLevel()
}

type Level1 struct {
	data LevelData
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

	PlaySoundWithLoop("music/level1", false)

	go func() {
		pausableSleep(time.Second * time.Duration(5 + rand.IntN(5)))
		train.Drive(4000, 0.2)
	}()

	go func() {
		pausableSleep(time.Second * time.Duration(20 + rand.IntN(5)))
		train2.Drive(4000, 0.2)
	}()

	for _, collider := range level.dynamicColliders {
		addGameObject(collider)
	}
}

func DrawLevel(screen *ebiten.Image, level Level) {
	tiles := level.GetData().Tiles
	sprites := level.GetData().Sprites
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
		for _, collider := range level.GetData().Colliders {

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

func (level *Level1) GetData() *LevelData {
	return &level.data
}


func (level *Level1) UpdateLevel() {
}

var currentLevel Level

func loadLevel1() {
	fmt.Println("LOADING LEVEL 1")
	loadedLevel := loadJson("level-tilesheets/level1.json", &LoadedLevel{})
	data := LevelData{
		Tiles: loadedLevel.Tiles,
		Sprites: map[string]*ebiten.Image{
			"rail":                           loadImage("assets/tiles/rail.png"),
			"rail-border-left":               loadImage("assets/tiles/rail-border-left.png"),
			"rail-border-right":              loadImage("assets/tiles/rail-border-right.png"),
			"station-floor-corner":           loadImage("assets/tiles/station-floor-corner.png"),
			"station-floor":                  loadImage("assets/tiles/station-floor.png"),
			"station-floor-protective":       loadImage("assets/tiles/station-floor-protective.png"),
			"station-floor-protective-right": loadImage("assets/tiles/station-floor-protective-right.png"),
	
			"bench":    loadImage("assets/tiles/bench.png"),
			"elevator": loadImage("assets/tiles/elevator.png"),
	
			"stairs": loadImage("assets/tiles/stairs.png"),
	
			"shadow": loadImage("assets/tiles/shadow.png"),
			"shadow-corner": loadImage("assets/tiles/shadow-corner.png"),
		},
		Colliders: loadedLevel.Colliders,
		CameraFollow: Vector2{x: 750},
		CameraZoom: 1.2,
	}
	currentLevel = &Level1{
		data: data,
	}
	currentLevel.StartLevel()
	player = CreatePlayer()
	addGameObject(player)
}

type Level2 struct {
	data LevelData
}

func (level *Level2) StartLevel() {
	// Placeholder for Level 2 start logic
}
func (level *Level2) UpdateLevel() {
	// Placeholder for Level 2 update logic
}

func (level *Level2) GetData() *LevelData {
	return &level.data
}

func loadLevel2() {
	fmt.Println("LOADING LEVEL 2")
	loadedLevel := loadJson("level-tilesheets/level2.json", &LoadedLevel{})
	data := LevelData{
		Tiles: loadedLevel.Tiles,
		Sprites: map[string]*ebiten.Image{
			"rail":                           loadImage("assets/tiles/rail.png"),
			"rail-border-left":               loadImage("assets/tiles/rail-border-left.png"),
			"rail-border-right":              loadImage("assets/tiles/rail-border-right.png"),
			"station-floor-corner":           loadImage("assets/tiles/station-floor-corner.png"),
			"station-floor":                  loadImage("assets/tiles/station-floor.png"),
			"station-floor-protective":       loadImage("assets/tiles/station-floor-protective.png"),
			"station-floor-protective-right": loadImage("assets/tiles/station-floor-protective-right.png"),
	
			"bench":    loadImage("assets/tiles/bench.png"),
			"elevator": loadImage("assets/tiles/elevator.png"),
	
			"stairs": loadImage("assets/tiles/stairs.png"),
	
			"shadow": loadImage("assets/tiles/shadow.png"),
			"shadow-corner": loadImage("assets/tiles/shadow-corner.png"),
	
			"wall": loadImage("assets/tiles/wall.png"),
			"level-layout": loadImage("assets/tiles/level-layout-hotlinemiami.png"),
		},
		Colliders: loadedLevel.Colliders,
		CameraZoom: 1.1,
	}
	currentLevel = &Level2{
		data: data,
	}
	currentLevel.StartLevel()
	player = CreatePlayer()
	addGameObject(player)
}

type Level0 struct {
	data LevelData
}

func (level *Level0) StartLevel() {
	// Placeholder for Level 0 start logic
}
func (level *Level0) UpdateLevel() {
	// Placeholder for Level 0 update logic
}

func (level *Level0) GetData() *LevelData { 
	return &level.data
}

func loadLevel0() {
	fmt.Println("LOADING LEVEL 0")
	loadedLevel := loadJson("level-tilesheets/level0.json", &LoadedLevel{})
	data := LevelData{
		Tiles: loadedLevel.Tiles,
		Sprites: map[string]*ebiten.Image{
			"rail":                           loadImage("assets/tiles/rail.png"),
			"rail-border-left":               loadImage("assets/tiles/rail-border-left.png"),
			"rail-border-right":              loadImage("assets/tiles/rail-border-right.png"),
			"station-floor-corner":           loadImage("assets/tiles/station-floor-corner.png"),
			"station-floor":                  loadImage("assets/tiles/station-floor.png"),
			"station-floor-protective":       loadImage("assets/tiles/station-floor-protective.png"),
			"station-floor-protective-right": loadImage("assets/tiles/station-floor-protective-right.png"),
	
			"bench":    loadImage("assets/tiles/bench.png"),
			"elevator": loadImage("assets/tiles/elevator.png"),
	
			"stairs": loadImage("assets/tiles/stairs.png"),
	
			"shadow": loadImage("assets/tiles/shadow.png"),
			"shadow-corner": loadImage("assets/tiles/shadow-corner.png"),
	
			"wall": loadImage("assets/tiles/wall.png"),
		},
		Colliders: loadedLevel.Colliders,
		CameraZoom: 1.1,
	}
	currentLevel = &Level0{
		data: data,
	}
	currentLevel.StartLevel()
	player = CreatePlayer()
	addGameObject(player)
}
