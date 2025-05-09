package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject interface {
    Update()
    Draw(screen *ebiten.Image)
    GetTransform() Transform
}

func drawRotatedRect(screen *ebiten.Image, x, y, width, height, angle float64, color color.Color) {
    // Create a rectangle image
    rect := ebiten.NewImage(int(width), int(height))
    rect.Fill(color)

    // Create a GeoM (geometric matrix) for transformations
    op := &ebiten.DrawImageOptions{}

    // Translate the rectangle to its center for rotation
    op.GeoM.Translate(-width/2, -height/2)

    // Rotate the rectangle
    op.GeoM.Rotate(angle)

    // Translate the rectangle back to its position
    op.GeoM.Translate(x, y)

    // Draw the rotated rectangle onto the screen
    screen.DrawImage(rect, op)
}