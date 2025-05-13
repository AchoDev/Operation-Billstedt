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

    playerX := player.transform.x
    playerY := player.transform.y

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

    checkPlayerCollision(playerX, playerY)


    updateKeyState()
    updateMouseState()

    return nil
}

func (g *Game) Draw (screen *ebiten.Image) {
    screen.Fill(color.White)

    rect := ebiten.NewImage(200, 100)
    rect.Fill(color.RGBA{0, 0, 0, 255})
    ebitenutil.DebugPrintAt(rect, fmt.Sprintf("%.2f", ebiten.ActualFPS()), 0, 0)
    ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Current gun: %s", player.currentGun.Name()), 0, 20)
    ebitenutil.DebugPrintAt(rect, fmt.Sprintf("Cooldown: %.2f", player.currentGun.GetCooldownTimer()), 0, 40)
    screen.DrawImage(rect, nil)

    for _, gameObject := range gameObjects {
        gameObject.Draw(screen)
    }
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
    return 1920, 1080
}

func checkPlayerCollision(playerX, playerY float64) {
    xCenterCircle := Circle{
        Center: Vector2{player.transform.x, playerY},
        Radius: player.transform.width,
    }

    yCenterCircle := 
}
 
func main() {
    gameObjects = append(gameObjects, &player)

    collider := Collider{
        transform: Transform{
            500, 500, 300, 300, 0,
        },
    }
    gameObjects = append(gameObjects, &collider)

    // ebiten.SetWindowSize(2000, 1700)
    ebiten.SetWindowSize(2000 / 2, 1700 / 2)
    ebiten.SetWindowTitle("Operation Billstedt")
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
    err := ebiten.RunGame(&Game{}); 
    
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