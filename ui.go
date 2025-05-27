package main

import "github.com/hajimehoshi/ebiten/v2"

type UI struct {
    transform Transform
}

func NewUI() *UI {
    return &UI{}
}

func (ui *UI) Draw(screen *ebiten.Image) {
    
}

func (ui *UI) Update() {
    // Update the UI state if necessary
}

func (ui *UI) GetTransform() Transform {
    return ui.transform
}

func (ui *UI) SetTransform(t Transform) {
    ui.transform = t
}
