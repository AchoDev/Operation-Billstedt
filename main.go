package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var gameObjects []GameObject = []GameObject{}
var player Player = CreatePlayer()

var camera Vector2 = Vector2{
    x: 0,
    y: 0,
}

type Game struct {}

func (g *Game) Update() error {
    for _, gameObject := range gameObjects {
        gameObject.Update()
    }

    return nil
}

func (g *Game) Draw (screen *ebiten.Image) {
    screen.Fill(color.White)

    rect := ebiten.NewImage(200, 100)
    rect.Fill(color.RGBA{0, 0, 0, 255})
    ebitenutil.DebugPrintAt(rect, fmt.Sprintf("%.2f", ebiten.ActualFPS()), 0, 0)
    ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Current gun: %s", player.currentGun.Name()), 0, 20)
    screen.DrawImage(rect, nil)

    for _, gameObject := range gameObjects {
        gameObject.Draw(screen)
    }
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
    return 1920, 1080
}

func main() {
    gameObjects = append(gameObjects, &player)

    gameObjects = append(gameObjects, createEnemy(100, 100))
    gameObjects = append(gameObjects, createEnemy(200, 200))
    gameObjects = append(gameObjects, createEnemy(300, 300))
    gameObjects = append(gameObjects, createEnemy(400, 400))

    ebiten.SetWindowSize(2000, 1700)
    ebiten.SetWindowTitle("Operation Billstedt")
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
    err := ebiten.RunGame(&Game{}); 
    
    if err != nil {
        log.Fatal(err)
    }
}
