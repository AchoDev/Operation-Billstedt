package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var pos = Vector2{300, 900}
var selectedSprite int


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
            width:  100,
            height: 100,
            rotation: 0,
        }

        if counter == selectedSprite {
            size := 10.0
            drawRect(screen, Transform{
                x:      transform.x - size,
                y:      transform.y - size,
                width:  transform.width + size*2,
                height: transform.height + size*2,
                rotation: 0,
            }, color.RGBA{100, 0, 0, 255})
        }

        drawImage(screen, sprite, transform)


        counter++
    }

    drawRect(screen, Transform{
        x:      1200,
        y: 100,
        width:  300,
        height: 100,  
        rotation: 0,
    }, color.RGBA{0, 0, 0, 255})

    ebitenutil.DebugPrintAt(screen, "Level Editor", 1050, 100)
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", pos), 1050, 130)
}

func UpdateLevelEditor(level Level) {
    
    if isKeyJustPressed(ebiten.KeyP) {
        levelEditorActivated = !levelEditorActivated
    }

    if isKeyJustPressed(ebiten.KeyTab) {
        selectedSprite++
        if selectedSprite >= len(level.GetSprites()) {
            selectedSprite = 0
        }
    }

    if isMouseButtonJustPressed(ebiten.MouseButtonLeft) {

        sprites := level.GetSprites()

        counter := 0
        for name := range sprites {
            if counter == selectedSprite {

                cursorX, cursorY := ebiten.CursorPosition()

                gridPosition := Vector2{
                    x: math.Round(float64(cursorX) / 100),
                    y: math.Round(float64(cursorY) / 100),
                }

                tiles := level.GetTiles()

                tiles = append(tiles, Tile{
                    x:      int(gridPosition.x),
                    y:      int(gridPosition.y),
                    width:  1,
                    height: 1,
                    sprite: name,
                    collide: false,
                })

                level.SetTiles(tiles)

                fmt.Println("Added tile at", gridPosition)
                break
            }

            counter++
        }
    }

    if !levelEditorActivated {
        return
    }
}