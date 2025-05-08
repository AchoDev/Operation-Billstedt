package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var gameObjects []GameObject = []GameObject{}


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


    for _, gameObject := range gameObjects {
        gameObject.Draw(screen)
    }
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
    return 1920, 1080
}

func main() {

    gameObjects = append(gameObjects, CreatePlayer())

    ebiten.SetWindowSize(2000, 1700)
    ebiten.SetWindowTitle("Operation Billstedt")
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
    err := ebiten.RunGame(&Game{}); 
    
    if err != nil {
        log.Fatal(err)
    }
}
