package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UI struct {
    transform Transform
    velocity Vector2
}

func NewUI() *UI {
    return &UI{
        transform: Transform{
            z: 200,
        },
    }
}

func (ui *UI) Draw(screen *ebiten.Image) {

    bg := getCachedImage("ui/gunbackground")

    op := defaultImageOptions()
    op.AlignItems = AlignEnd
    op.JustifyContent = AlignStart
    op.Margin = Vector2{-30, -30}

    op.Margin.x += ui.velocity.x * 0.2
    op.Margin.y += ui.velocity.y * 0.2

    op.OriginalImageSize = true
    
    op.Scale.Set2(0.20, 0.20)
    op.Skew.Set2(0.11, 0.12)

    drawAbsoluteImageWithOptions(
        screen,
        bg,
        ui.transform,
        op,
    )

    pistol := getCachedImage("ui/pistol")
    shotgun := getCachedImage("ui/shotgun")
    rifle := getCachedImage("ui/rifle")

    x := 500.0
    for _, gun := range []*ebiten.Image{pistol, shotgun, rifle} {
        op := defaultImageOptions()
        op.AlignItems = AlignEnd
        op.JustifyContent = AlignStart
        op.Margin = Vector2{0, -150}
        op.OriginalImageSize = true
        op.Scale.Set(0.2)
        drawAbsoluteImageWithOptions(
            screen,
            gun,
            Transform{
                x: x,
            },
            op,
        )
        x += 400
    }

    bigGun := pistol
    switch player.currentGun.name {
    case "Pistol":
        bigGun = pistol
    case "Shotgun":
        bigGun = shotgun
    case "Rifle":
        bigGun = rifle
    case "Minigun":
        bigGun = rifle
    }

    bOp := defaultImageOptions()
    bOp.OriginalImageSize = true
    bOp.Scale.Set(0.6)
    drawAbsoluteImageWithOptions(
        screen,
        bigGun,
        Transform{
            x: 1000,
            y: 500,
        },
        bOp,
    )

    
}

func (ui *UI) Update() {
    ui.velocity.x -= player.velocity.x
    ui.velocity.y -= player.velocity.y

    ui.velocity.Multiply(0.8)
}

func (ui *UI) GetTransform() Transform {
    return ui.transform
}

func (ui *UI) SetTransform(t Transform) {
    ui.transform = t
}
