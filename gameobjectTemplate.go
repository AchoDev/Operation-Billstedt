package main

import "github.com/hajimehoshi/ebiten/v2"

type Template struct {
    transform Transform
}

func (template *Template) Update() {
    // Update the transform if needed
}

func (template *Template) Draw(screen *ebiten.Image) {

}

func (template *Template) GetTransform() Transform {
    return template.transform
}

func (template *Template) SetTransform(transform Transform) {
    template.transform = transform
}