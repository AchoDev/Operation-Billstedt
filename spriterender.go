package main

import "github.com/hajimehoshi/ebiten/v2"

type SpriteRenderer struct {
    transform Transform
    sprite *ebiten.Image
    options ImageOptions
}

func (renderer *SpriteRenderer) Update() {
    // Update the transform if needed
}

func (renderer *SpriteRenderer) Draw(screen *ebiten.Image) {
    drawImageWithOptions(screen, renderer.sprite, renderer.transform, renderer.options)
}

func (renderer *SpriteRenderer) GetTransform() Transform {
    return renderer.transform
}

func (renderer *SpriteRenderer) SetTransform(transform Transform) {
    renderer.transform = transform
}