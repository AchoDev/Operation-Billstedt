package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
    transform Transform
    text string
    action func()
}

func (b *Button) Update() {
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        mouseXInt, mouseYInt := ebiten.CursorPosition()
        mouseX := float64(mouseXInt)
        mouseY := float64(mouseYInt)
        if b.transform.x <= mouseX && mouseX <= b.transform.x+b.transform.width &&
           b.transform.y <= mouseY && mouseY <= b.transform.y+b.transform.height {
            b.action()
        }
    }
}

func (b *Button) Draw(screen *ebiten.Image) {
    rect := getCachedRect(200, 50, color.RGBA{0, 0, 0, 100})
    op := defaultImageOptions()
    op.OriginalImageSize = true

    tr := Transform {
        x: b.transform.x + screenPosition.x * 1920,
        y: b.transform.y + screenPosition.y * 1080,
    }

    drawAbsoluteImageWithOptions(screen, rect, tr, op)

    drawText(screen, b.text, tr)
}

var screenPosition Vector2 = Vector2{0, 0}

var mainButtons []*Button = []*Button{
    {
        transform: Transform{x: 100, y: 100},
        text: "Play",
        action: func() {
            moveScreenPosition(0, 1)
        },
    },
}

func moveScreenPosition(dx, dy float64) {}

func UpdateMainMenu() {
    for _, button := range mainButtons {
        button.Update()
    }
}

func DrawMainMenu(screen *ebiten.Image) {
    screen.Fill(color.White)

    for _, button := range mainButtons {
        button.Draw(screen)
    }
}