package main

import "github.com/hajimehoshi/ebiten/v2"

type UI struct {
    transform Transform
}

func NewUI() *UI {
    return &UI{}
}

func (ui *UI) Draw(screen *ebiten.Image) {

    op := defaultImageOptions()
    op.AlignItems = AlignEnd
    op.JustifyContent = AlignStart
    op.Margin = Vector2{30, -30}
    op.OriginalImageSize = true
    op.Scale.Set(0.1)
    // op.Skew.x = 0.1

    drawAbsoluteImageWithOptions(
        screen,
        getCachedImage("ui/gunbackground"),
        ui.transform,
        op,
    )
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
