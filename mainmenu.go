package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
    transform Transform
    text string
    action func()
}

func (b *Button) Update() {
    if isMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        mouseXInt, mouseYInt := ebiten.CursorPosition()

        mouseX := float64(mouseXInt)
        mouseY := float64(mouseYInt)

        x := b.transform.x - b.transform.width/2 + actualScreenPosition.x
        y := b.transform.y - b.transform.height/2 + actualScreenPosition.y

        fmt.Println("Button position:", x, y, "Mouse position:", mouseX, mouseY)

        if x <= mouseX && mouseX <= x+b.transform.width &&
           y <= mouseY && mouseY <= y+b.transform.height {
            fmt.Println("Button clicked:", b.text)
            b.action()
        }
    }
}

func (b *Button) Draw(screen *ebiten.Image) {
    rect := getCachedRect(int(b.transform.width), int(b.transform.height), color.RGBA{0, 0, 0, 100})
    op := defaultImageOptions()
    op.OriginalImageSize = true

    tr := Transform {
        x: b.transform.x + actualScreenPosition.x,
        y: b.transform.y + actualScreenPosition.y,
    }

    drawAbsoluteImageWithOptions(screen, rect, tr, op)

    drawText(screen, b.text, tr)
}

var screenPosition Vector2 = Vector2{0, 0}
var actualScreenPosition Vector2 = Vector2{0, 0}

var mainButtons []*Button = []*Button{
    {
        transform: Transform{x: 100, y: 100, width: 200, height: 50},
        text: "Play",
        action: func() {
            moveScreenPosition(0, -1)
        },
    },

    {
        transform: Transform{x: 500, y: 1600, width: 200, height: 50},
        text: "Level 1",
        action: func() {
            loadLevel1()
            mainMenuActivated = false
        },
    },

    {
        transform: Transform{x: 500, y: 1900, width: 200, height: 50},
        text: "Level 2",
        action: func() {
            loadLevel2()
            mainMenuActivated = false
        },
    },
}



func moveScreenPosition(dx, dy int) {
    screenPosition.x += float64(dx)
    screenPosition.y += float64(dy)
}

func UpdateMainMenu() {
    for _, button := range mainButtons {
        button.Update()
    }

    target := Vector2{}
    target.SetVector(screenPosition)
    target.Multiply2(1920, 1080)
    actualScreenPosition.SmoothMove(target, 0.1)
}

func DrawMainMenu(screen *ebiten.Image) {
    screen.Fill(color.White)

    for _, button := range mainButtons {
        button.Draw(screen)
    }
}