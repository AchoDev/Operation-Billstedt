package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var gameObjects []GameObject = []GameObject{}
var player Player = CreatePlayer()

var camera Vector2 = Vector2{
	x: 0,
	y: 0,
}

var currentLevel Level = &Level1{
	tiles: loadJson("level-tilesheets/level1.json", &[]Tile{}),
	sprites: map[string]*ebiten.Image{
		"rail":                     loadImage("assets/tiles/rail.png"),
		"station-floor-corner":     loadImage("assets/tiles/station-floor-corner.png"),
		"station-floor":            loadImage("assets/tiles/station-floor.png"),
		"station-floor-protective": loadImage("assets/tiles/station-floor-protective.png"),
	},
}

type Game struct{}

func (g *Game) Update() error {

	playerX := player.transform.x
	playerY := player.transform.y

	if !levelEditorActivated {
		for _, gameObject := range gameObjects {
			gameObject.Update()
		}

		if isKeyJustPressed(ebiten.Key9) {
			fmt.Println("Creating new player")
			if findPlayer() == nil {
				player = CreatePlayer()
				gameObjects = append(gameObjects, &player)
			}
		}
		currentLevel.UpdateLevel()
		checkCollisions(playerX, playerY)
	}

	UpdateLevelEditor(currentLevel)

	updateKeyState()
	updateMouseState()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	rect := ebiten.NewImage(200, 100)
	rect.Fill(color.RGBA{0, 0, 0, 255})
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("%.2f", ebiten.ActualFPS()), 0, 0)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Current gun: %s", player.currentGun.Name()), 0, 20)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Cooldown: %.2f", player.currentGun.GetCooldownTimer()), 0, 40)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Camera pos: %.2f %.2f", camera.x, camera.y), 0, 60)
	screen.DrawImage(rect, nil)

	DrawLevel(screen, currentLevel)

	for _, gameObject := range gameObjects {
		gameObject.Draw(screen)
	}

	DrawLevelEditor(screen, currentLevel)
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

func checkCollisions(playerX, playerY float64) {
	xCenterCircle := Circle{
		Center: Vector2{player.transform.x, playerY},
		Radius: player.transform.width / 2,
	}

	yCenterCircle := Circle{
		Center: Vector2{playerX, player.transform.y},
		Radius: player.transform.width / 2,
	}

	var bullets []*Bullet

	for _, gameObj := range gameObjects {
		if bullet, ok := gameObj.(*Bullet); ok {
			bullets = append(bullets, bullet)
		}
	}

	for _, gameObj := range gameObjects {
		if collider, ok := gameObj.(*Collider); ok {
			rect := Rect{
				Center: Vector2{
					collider.transform.x,
					collider.transform.y,
				},
				Width:  collider.transform.width,
				Height: collider.transform.height,
			}

			if CircleRotatedRectColliding(xCenterCircle, rect) {
				player.transform.x = playerX
			}

			if CircleRotatedRectColliding(yCenterCircle, rect) {
				player.transform.y = playerY
			}

			for _, bullet := range bullets {
				if RotatedRectsColliding(createRectFromTransform(bullet.transform), rect) {
					removeGameObject(bullet)
				}
			}

		}
	}
}

func createRectFromTransform(transform Transform) Rect {
	return Rect{
		Center: Vector2{
			transform.x, transform.y,
		},
		Width:  transform.width,
		Height: transform.height,
		Angle:  transform.rotation,
	}
}

func removeGameObject(target GameObject) {
	for i, gameObj := range gameObjects {
		if gameObj == target {
			gameObjects = append(gameObjects[:i], gameObjects[i+1:]...)
		}
	}
}

func getGameobjectsOfType[T GameObject]() []T {
	var list []T
	for _, gameObj := range gameObjects {
		if obj, ok := gameObj.(T); ok {
			list = append(list, obj)
		}
	}

	return list
}

func main() {
	gameObjects = append(gameObjects, &player)

	collider := Collider{
		transform: Transform{
			1000, 500, 100, 100, 0,
		},
	}
	gameObjects = append(gameObjects, &collider)

	ebiten.SetWindowSize(1920, 1080)
	// ebiten.SetWindowSize(2000, 1700)
	// ebiten.SetWindowSize(2000 / 2, 1700 / 2)
	ebiten.SetWindowTitle("Operation Billstedt")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(&Game{})

	if err != nil {
		log.Fatal(err)
	}
}

func findPlayer() *Player {
	for _, gameobj := range gameObjects {
		if player, ok := gameobj.(*Player); ok {
			return player
		}
	}

	return nil
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func loadJson[T any](path string, target *T) T {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loaded JSON file:", path)

	return *target
}
