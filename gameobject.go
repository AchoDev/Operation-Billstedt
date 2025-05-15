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

func drawRect(screen *ebiten.Image, transform Transform, color color.Color) {
	rect := ebiten.NewImage(int(transform.width), int(transform.height))
	rect.Fill(color)
	drawImageWithOptions(screen, rect, transform, defaultImageOptions())
}

func drawImage(screen *ebiten.Image, image *ebiten.Image, transform Transform) {
	drawImageWithOptions(screen, image, transform, defaultImageOptions())
}

func drawImageWithOptions(screen *ebiten.Image, image *ebiten.Image, transform Transform, options ImageOptions) {
    if transform.x + transform.width / 2 < camera.x - camera.width / camera.zoom / 2 || transform.x - transform.width / 2 > camera.x + camera.width / camera.zoom / 2  {
        return
    }

    if transform.y + transform.height / 2 < camera.y - camera.height / camera.zoom / 2 || transform.y - transform.height / 2 > camera.y + camera.height / camera.zoom / 2 {
        return
    }


	transform.x -= camera.x
	transform.y -= camera.y
    
    transform.x *= camera.zoom
    transform.y *= camera.zoom


	transform.x += camera.width / 2
	transform.y += camera.height / 2


    transform.width *= camera.zoom
    transform.height *= camera.zoom
    

	drawAbsoluteImageWithOptions(screen, image, transform, options)
}

func drawAbsoluteRect(screen *ebiten.Image, transform Transform, color color.Color) {
	rect := ebiten.NewImage(int(transform.width), int(transform.height))
	rect.Fill(color)
	drawAbsoluteImage(screen, rect, transform)
}

func drawAbsoluteImage(screen *ebiten.Image, image *ebiten.Image, transform Transform) {
	drawAbsoluteImageWithOptions(screen, image, transform, defaultImageOptions())
}

func drawAbsoluteImageWithOptions(screen *ebiten.Image, image *ebiten.Image, transform Transform, options ImageOptions) {    
    op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(image.Bounds().Dx())/2, -float64(image.Bounds().Dy())/2) // Center the sprite
	op.GeoM.Translate(-options.Anchor.x, -options.Anchor.y)                             // Center the sprite
	op.GeoM.Rotate(transform.rotation)                                                  // Rotate the sprite
	op.GeoM.Scale(transform.width/float64(image.Bounds().Dx()), transform.height/float64(image.Bounds().Dy()))
	op.GeoM.Scale(options.Scale, options.Scale)              // Scale the sprite
	op.GeoM.Translate(transform.x, transform.y)              // Offset the sprite position
	op.ColorScale.ScaleAlpha(float32(options.Alpha) / 255.0) // Set the alpha value
	screen.DrawImage(image, op)
}

type ImageOptions struct {
	Anchor Vector2
	Alpha  float64
	Scale  float64
}

func defaultImageOptions() ImageOptions {
	return ImageOptions{
		Anchor: Vector2{0, 0},
		Alpha:  255,
		Scale:  1,
	}
}
