package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

var skew Vector2 = Vector2{0.11, 0.12}


func (ui *UI) Draw(screen *ebiten.Image) {

    bg := getCachedImage("ui/gunbackground")

    op := defaultImageOptions()
    op.AlignItems = AlignEnd
    op.JustifyContent = AlignStart
    op.Margin = Vector2{-30, -30}

    op.Margin.x += ui.velocity.x * 0.2
    op.Margin.y += ui.velocity.y * 0.2

    op.OriginalImageSize = true
    
    op.Scale.Set(0.2)
    op.Skew.SetVector(skew)

    drawAbsoluteImageWithOptions(
        screen,
        bg,
        ui.transform,
        op,
    )

    pistol := getCachedImage("ui/pistol")
    shotgun := getCachedImage("ui/shotgun")
    rifle := getCachedImage("ui/rifle")

    x := 50.0
    y := -40.0
    for _, gun := range []*ebiten.Image{pistol, shotgun, rifle} {
        op.Margin.y = y
        op.Margin.x = x + ui.velocity.x * 0.3
        op.Margin.y += ui.velocity.y * 0.3
        op.Scale.Set(0.06)
        drawAbsoluteImageWithOptions(
            screen,
            gun,
            Transform{
                x: x,
            },
            op,
        )
        x += 120
        y += 15
    }

    textOp := &text.DrawOptions{}
    textOp.GeoM.Translate(100, 600)
    text.Draw(
        screen, 
        fmt.Sprintf("%d", player.currentGun.currentAmmo),
        customFont,
        textOp,
    )

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

    op.Scale.Set(0.13)
    op.Margin.Set2(30, -100)
    op.Margin.x += ui.velocity.x * 0.5
    op.Margin.y += ui.velocity.y * 0.5

    drawAbsoluteImageWithOptions(
        screen,
        bigGun,
        Transform{
            x: 500,
            y: 100,
        },
        op,
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
