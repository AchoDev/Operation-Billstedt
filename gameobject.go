package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameObject interface {
	Update()
	Draw(screen *ebiten.Image)
	GetTransform() Transform
	SetTransform(Transform)
}

var rectCache = make(map[string]*ebiten.Image)
var imageCache = make(map[string]*ebiten.Image)

func getCachedRect(width, height int, color color.Color) *ebiten.Image {
    key := fmt.Sprintf("%dx%d-%v", width, height, color)
	if rect, ok := rectCache[key]; ok {
		return rect
	}

	rect := ebiten.NewImage(width, height)
	rect.Fill(color)
	rectCache[key] = rect
	return rect
}

func getCachedImage(path string) *ebiten.Image {
	if image, ok := imageCache[path]; ok {
		return image
	}

	image, _, err := ebitenutil.NewImageFromFile("assets/" + path + ".png")
	if err != nil {
		panic(err)
	}
	imageCache[path] = image
	return image
}

func drawRect(screen *ebiten.Image, transform Transform, color color.Color) {
	rect := getCachedRect(int(transform.width), int(transform.height), color)
	drawImageWithOptions(screen, rect, transform, defaultImageOptions())
}

func drawImage(screen *ebiten.Image, image *ebiten.Image, transform Transform) {
	drawImageWithOptions(screen, image, transform, defaultImageOptions())
}

func drawImageWithOptions(screen *ebiten.Image, image *ebiten.Image, transform Transform, options ImageOptions) {
	if options.OriginalImageSize {
		transform.width = float64(image.Bounds().Dx())
		transform.height = float64(image.Bounds().Dy())
	}



	if transform.x+(transform.width * options.Scale)/2 < camera.x-camera.width/camera.zoom/2 || transform.x-(transform.width * options.Scale)/2 > camera.x+camera.width/camera.zoom/2 {
		return
	}

	if transform.y+(transform.height * options.Scale)/2 < camera.y-camera.height/camera.zoom/2 || transform.y-(transform.height * options.Scale)/2 > camera.y+camera.height/camera.zoom/2 {
		return
	}

	transform.x -= camera.x
	transform.y -= camera.y

	transform.x -= camera.offset.x
	transform.y -= camera.offset.y

	transform.x *= camera.zoom
	transform.y *= camera.zoom

	transform.x += camera.width / 2
	transform.y += camera.height / 2

	transform.width *= camera.zoom
	transform.height *= camera.zoom

	drawAbsoluteImageWithOptions(screen, image, transform, options)
}

func drawAbsoluteRect(screen *ebiten.Image, transform Transform, color color.Color) {
	rect := getCachedRect(int(transform.width), int(transform.height), color)
	drawAbsoluteImage(screen, rect, transform)
}

func drawAbsoluteImage(screen *ebiten.Image, image *ebiten.Image, transform Transform) {
	drawAbsoluteImageWithOptions(screen, image, transform, defaultImageOptions())
}

func drawAbsoluteImageWithOptions(screen *ebiten.Image, image *ebiten.Image, transform Transform, options ImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(image.Bounds().Dx())/2, -float64(image.Bounds().Dy())/2) // Center the sprite
	op.GeoM.Translate(-options.Anchor.x, -options.Anchor.y)                             // Center the sprite
	if options.FlipX {
		op.GeoM.Scale(-1, 1) // Flip the sprite horizontally
	}
	if options.FlipY {
		op.GeoM.Scale(1, -1)
	}
	op.GeoM.Rotate(transform.rotation)

	op.GeoM.Scale(transform.width/float64(image.Bounds().Dx()), transform.height/float64(image.Bounds().Dy()))
	op.GeoM.Scale(options.Scale, options.Scale)              // Scale the sprite
	op.GeoM.Translate(transform.x, transform.y)              // Offset the sprite position
	op.ColorScale.ScaleAlpha(float32(options.Alpha) / 255.0) // Set the alpha value
	screen.DrawImage(image, op)
}

type ImageOptions struct {
	Anchor            Vector2
	Alpha             float64
	Scale             float64
	OriginalImageSize bool
	FlipX 		  bool
	FlipY bool
}

func defaultImageOptions() ImageOptions {
	return ImageOptions{
		Anchor:            Vector2{0, 0},
		Alpha:             255,
		Scale:             1,
		OriginalImageSize: false,
		FlipX: false,
		FlipY: false,
	}
}
