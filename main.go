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

type Camera struct {
	x      float64
	y      float64
	width  float64
	height float64
	zoom   float64
}

var camera Camera = Camera{
	x:      0,
	y:      0,
	width:  1920,
	height: 1080,
	zoom:   1.2,
}

type LoadedLevel struct {
	Tiles     []Tile `json:"tiles"`
	Colliders []Tile `json:"colliders"`
}

var loadedLevel LoadedLevel = loadJson("level-tilesheets/level1.json", &LoadedLevel{})

var currentLevel Level = &Level1{
	tiles: loadedLevel.Tiles,
	sprites: map[string]*ebiten.Image{
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
	colliders: loadedLevel.Colliders,
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

		if isKeyJustPressed(ebiten.KeyK) {
			invincible = !invincible
		}

		currentLevel.UpdateLevel()
		checkCollisions(&player.transform, Vector2{playerX, playerY})

		moveCamera()
	}

	UpdateLevelEditor(currentLevel)

	updateKeyState()
	updateMouseState()

	return nil
}

var debugRect *ebiten.Image = ebiten.NewImage(220, 200)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	DrawLevel(screen, currentLevel)

	// for _, gameObject := range gameObjects {
	// 	if levelEditorActivated && hideGameobjects {
	// 		break
	// 	}
	// 	gameObject.Draw(screen)
	// }

	DrawLevelEditor(screen, currentLevel)

	// rect := ebiten.NewImage(200, 100)
	rect := debugRect
	rect.Fill(color.Black)
	ebitenutil.DebugPrintAt(rect, "Operation Billstedt Prev. Version 2", 0, 0)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("%.2f", ebiten.ActualFPS()), 0, 20)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Current gun: %s", player.currentGun.Name()), 0, 40)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Cooldown: %.2f", player.currentGun.GetCooldownTimer()), 0, 60)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Camera pos: %.2f %.2f", camera.x, camera.y), 0, 80)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Camera Zoom: %.2f", camera.zoom), 0, 100)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Player pos: %.2f %.2f", player.transform.x, player.transform.y), 0, 120)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Remaining Enemies: %d", len(getGameobjectsOfType[*Enemy]())), 0, 140)
	ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Invincibility: %t", invincible), 0, 160)

	screen.DrawImage(rect, nil)
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

func moveCamera() {
	diff := Vector2{
		x: 750 - camera.x,
		y: player.transform.y - camera.y,
	}

	zoomDiff := camera.zoom - 1.2

	camera.zoom -= zoomDiff * 0.1
	camera.x += diff.x * 0.1
	camera.y += diff.y * 0.1
}

func checkCollisions(tr *Transform, startPosition Vector2) {

	xCenterCircle := Circle{
		Center: Vector2{tr.x, startPosition.y},
		Radius: tr.width / 2,
	}

	yCenterCircle := Circle{
		Center: Vector2{startPosition.x, tr.y},
		Radius: tr.width / 2,
	}

	for _, gameObj := range gameObjects {

		if gameObj == nil {
			continue
		}

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
				tr.x = startPosition.x
			}

			if CircleRotatedRectColliding(yCenterCircle, rect) {
				tr.y = startPosition.y
			}
		}
	}

	for _, col := range currentLevel.GetColliders() {
		rect := Rect{
			Center: Vector2{
				col.X * 100,
				col.Y * 100,
			},
			Width:  col.Width * 100,
			Height: col.Height * 100,
		}

		if CircleRotatedRectColliding(xCenterCircle, rect) {
			tr.x = startPosition.x
		}

		if CircleRotatedRectColliding(yCenterCircle, rect) {
			tr.y = startPosition.y
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


	gameObjects = append(gameObjects, NewHealthBar())

	gameObjects = append(gameObjects, &player)

	ebiten.SetWindowSize(1920, 1080)
	// ebiten.SetWindowSize(2000, 1700)
	// ebiten.SetWindowSize(2000 / 2, 1700 / 2)
	ebiten.SetWindowTitle("Operation Billstedt")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	currentLevel.StartLevel()

	if err := ebiten.RunGame(&Game{}); err != nil {
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
