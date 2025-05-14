package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var pos = Vector2{0, 0}

func DrawLevelEditor(screen *ebiten.Image, level Level) {
    if !levelEditorActivated {
        return
    }
    
    counter := 0
    // y := 1080.0
    for _, sprite := range level.GetSprites() {
        transform := Transform{
            x:      float64(counter * 20) + pos.x,
            y:      pos.y,
            width:  300,
            height: 300,
            rotation: 0,
        }

        drawImage(screen, sprite, transform)

        counter++
    }

    ebitenutil.DebugPrintAt(screen, "Level Editor", 10, 10)
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", pos), 10, 30)
}

func UpdateLevelEditor(level Level) {
    
    if isKeyJustPressed(ebiten.KeyP) {
        levelEditorActivated = !levelEditorActivated
    }

    if isKeyJustPressed(ebiten.KeyUp){
        pos.y -= 10
    }
    if isKeyJustPressed(ebiten.KeyDown){
        pos.y += 10
    }
    if isKeyJustPressed(ebiten.KeyLeft){
        pos.x -= 10

    }
    if isKeyJustPressed(ebiten.KeyRight){
        pos.x += 10
    }

    if !levelEditorActivated {
        return
    }
}