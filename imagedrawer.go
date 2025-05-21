package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageDrawer struct {
	transform Transform
    path     string
}

func (collider *ImageDrawer) Update() {

}

func (collider *ImageDrawer) Draw(screen *ebiten.Image) {
	image := getCachedImage(collider.path)
	
    drawImage(screen, image, collider.transform)
}

func (collider *ImageDrawer) GetTransform() Transform {
	return collider.transform
}

func (collider *ImageDrawer) SetTransform(transform Transform) {
	collider.transform = transform
}