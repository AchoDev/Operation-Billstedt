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

func drawRect(screen *ebiten.Image, transform Transform, color color.Color) {
    rect := ebiten.NewImage(int(transform.width), int(transform.height))
    rect.Fill(color)
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(-transform.width / 2, -transform.height / 2)
    op.GeoM.Translate(transform.x, transform.y)

    screen.DrawImage(rect, op)
}

func drawImage(screen *ebiten.Image, image *ebiten.Image, transform Transform) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(-float64(image.Bounds().Dx())/2, -float64(image.Bounds().Dy())/2) // Center the sprite
    op.GeoM.Translate(transform.x, transform.y) // Offset the sprite position
    op.GeoM.Rotate(transform.rotation + 3.14/2) // Rotate the sprite
    op.GeoM.Scale(transform.width / float64(image.Bounds().Dx()), transform.height / float64(image.Bounds().Dy()))
    screen.DrawImage(image, op)
}
