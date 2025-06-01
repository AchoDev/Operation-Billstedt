package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type GameObject interface {
	Update()
	Draw(screen *ebiten.Image)
	GetTransform() Transform
	SetTransform(Transform)
}

type Alignment int

const (
	AlignNone Alignment = iota
	AlignCenter
	AlignStart
	AlignEnd
)

var gameObjectsMutex = &sync.Mutex{}
var rectCache = make(map[string]*ebiten.Image)
var imageCache = make(map[string]*ebiten.Image)

func removeGameObject(obj GameObject) {
    gameObjectsMutex.Lock()
    defer gameObjectsMutex.Unlock()

    for i, gameObject := range gameObjects {
        if gameObject == obj {
            gameObjects = append(gameObjects[:i], gameObjects[i+1:]...)
            return
        }
    }
}

func addGameObject(obj GameObject) {
    gameObjectsMutex.Lock()
    defer gameObjectsMutex.Unlock()

    gameObjects = append(gameObjects, obj)
}

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

	if transform.x+(transform.width*options.Scale.x)/2 < camera.x-camera.width/camera.zoom/2 || transform.x-(transform.width*options.Scale.x)/2 > camera.x+camera.width/camera.zoom/2 {
		return
	}

	if transform.y+(transform.height*options.Scale.y)/2 < camera.y-camera.height/camera.zoom/2 || transform.y-(transform.height*options.Scale.y)/2 > camera.y+camera.height/camera.zoom/2 {
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

	options.Scale.Multiply(camera.zoom)

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
	
	if options.OriginalImageSize {
		transform.width = float64(image.Bounds().Dx())
		transform.height = float64(image.Bounds().Dy())
	}

	currentSize := transform.GetSize()
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(image.Bounds().Dx())/2, -float64(image.Bounds().Dy())/2) // Center the sprite
	op.GeoM.Translate(-options.Anchor.x, -options.Anchor.y)                             // Center the sprite
	if options.FlipX {
		op.GeoM.Scale(-1, 1) // Flip the sprite horizontally
	}
	if options.FlipY {
		op.GeoM.Scale(1, -1)
	}

	if options.ScaleColor {
		op.ColorScale.SetR(float32(options.ColorScale.R))
		op.ColorScale.SetG(float32(options.ColorScale.G))
		op.ColorScale.SetB(float32(options.ColorScale.B))
	}

	if killscreen && !options.KillScreenExclusion {
		if options.KillScreenBlack {
			op.ColorScale.Scale(0, 0, 0, 1)
		} else {
			op.ColorScale.SetR(255)
			op.ColorScale.SetG(255)
			op.ColorScale.SetB(255)
		}
	}


	scaleX := transform.width / float64(image.Bounds().Dx())
	scaleY := transform.height / float64(image.Bounds().Dy())
	currentSize.Multiply2(scaleX, scaleY)
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.SetElement(0, 1, options.Skew.x) // Apply skew in x direction
	op.GeoM.SetElement(1, 0, options.Skew.y) // Apply skew in y direction
	currentSize.Multiply2(options.Scale.x, options.Scale.y)
	op.GeoM.Scale(options.Scale.x, options.Scale.y) // Scale the sprite
	op.GeoM.Rotate(transform.rotation)

	switch options.JustifyContent {
	case AlignNone:
		op.GeoM.Translate(transform.x, 0)             // Offset the sprite position
	case AlignCenter:
		op.GeoM.Translate(float64(screen.Bounds().Dx()) / 2, 0)
	case AlignStart:
		op.GeoM.Translate(0, 0)
		op.GeoM.Translate(currentSize.x / 2, 0)
	case AlignEnd:
		op.GeoM.Translate(float64(screen.Bounds().Dx()), 0)
		op.GeoM.Translate(-currentSize.x / 2, 0)
	}

	switch options.AlignItems {
	case AlignNone:
		op.GeoM.Translate(0, transform.y) // Offset the sprite position
	case AlignCenter:
		op.GeoM.Translate(0, float64(screen.Bounds().Dy())/2)
	case AlignStart:
		op.GeoM.Translate(0, 0)
		op.GeoM.Translate(0, currentSize.y/2) // Center the sprite vertically
	case AlignEnd:
		op.GeoM.Translate(0, float64(screen.Bounds().Dy()))
		op.GeoM.Translate(0, -currentSize.y/2) // Center the sprite vertically
	}

	op.GeoM.Translate(options.Margin.x, options.Margin.y) // Apply margin

	op.ColorScale.ScaleAlpha(float32(options.Alpha) / 255.0) // Set the alpha value
	op.Filter = ebiten.FilterLinear

	screen.DrawImage(image, op)
}

func drawText(screen *ebiten.Image, characters string, transform Transform) {
	if customFont == nil {
		loadFont()
	}

	op := text.DrawOptions{}
	op.PrimaryAlign = text.AlignCenter
	op.GeoM.Translate(transform.x, transform.y)


	text.Draw(screen, characters, customFont, &op)
}

type ImageOptions struct {
	Anchor            Vector2
	Alpha             float64
	Scale             Vector2
	OriginalImageSize bool
	FlipX             bool
	FlipY             bool
	ScaleColor bool
	ColorScale        color.RGBA
	KillScreenBlack bool
	KillScreenExclusion bool
	JustifyContent Alignment
	AlignItems Alignment
	Margin Vector2
	Skew Vector2 // Skew is not used in this example, but can be implemented if needed
}

func defaultImageOptions() ImageOptions {
	return ImageOptions{
		Anchor:            Vector2{0, 0},
		Alpha:             255,
		Scale:             Vector2{1, 1},
		OriginalImageSize: false,
		FlipX:             false,
		FlipY:             false,
		ColorScale:        color.RGBA{255, 255, 255, 255},
		JustifyContent: AlignNone,
		AlignItems: AlignNone,
		Margin: 		  Vector2{0, 0},
		KillScreenBlack: false,
		KillScreenExclusion: false,
	}
}

var customFont text.Face

func loadFont() {
    // Load the font file
    fontBytes, err := os.ReadFile("assets/font.ttf")

    if err != nil {
        log.Fatal(err)
    }

    // Parse the font
    parsedFont, err := opentype.Parse(fontBytes)
    if err != nil {
        log.Fatal(err)
    }

    // Create a font face with a specific size
    const fontSize = 24
	fontFace, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	customFont = text.NewGoXFace(fontFace)
    if err != nil {
        log.Fatal(err)
    }
}
