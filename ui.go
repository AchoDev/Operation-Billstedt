package main

import (
	"fmt"
	"image/color"
	"strings"

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

    pistol := &player.guns[0]
    shotgun := &player.guns[1]
    rifle := &player.guns[2]

    x := 50.0
    y := -40.0
    for _, gun := range []*GunStats{pistol, shotgun, rifle} {
        op.Margin.y = y
        op.Margin.x = x + ui.velocity.x * 0.3
        op.Margin.y += ui.velocity.y * 0.3
        op.Scale.Set(0.06)
        sprite := getCachedImage("ui/" + strings.ToLower(gun.name))
        drawAbsoluteImageWithOptions(
            screen,
            sprite,
            Transform{
                x: x,
            },
            op,
        )

        if gun.cooldownTimer != -1 {
            rect := getCachedRect(2000, 2000, color.Black)
            op.Anchor.y = 200
            op.Alpha = 100
            percentage := float64(gun.cooldownTimer) / float64(gun.cooldown)
            op.Scale.y *= percentage
            op.Skew.Multiply(percentage)
            drawAbsoluteImageWithOptions(
                screen,
                rect,
                Transform{
                    x: x,
                },
                op,
            )

            op.Skew.SetVector(skew)
            op.Anchor.y = 0
            op.Alpha = 255
        }

        x += 120
        y += 15
    }

    textOp := &text.DrawOptions{}
    textOp.GeoM.Translate(350, 800)
    textOp.ColorScale.Scale(0, 0, 0, 255)
    textOp.GeoM.Translate(ui.velocity.x * 0.4, ui.velocity.y * 0.4)
    text.Draw(
        screen, 
        fmt.Sprintf("%d", player.currentGun.currentAmmo),
        customFont,
        textOp,
    )

    textOp.GeoM.Translate(1, 1)
    textOp.ColorScale.SetR(255)
    textOp.ColorScale.SetG(255)
    textOp.ColorScale.SetB(255)
    text.Draw(
        screen, 
        fmt.Sprintf("%d", player.currentGun.currentAmmo),
        customFont,
        textOp,
    )

    op.Scale.Set(0.2)
    op.Margin.Set2(350, -100)
    op.Margin.x += ui.velocity.x * 0.4
    op.Margin.y += ui.velocity.y * 0.4
    drawAbsoluteImageWithOptions(
        screen,
        getCachedImage("sprites/bullet"),
        Transform{},
        op,
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
        getCachedImage("ui/" + strings.ToLower(bigGun.name)),
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
