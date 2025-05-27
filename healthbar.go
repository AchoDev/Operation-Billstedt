package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type HealthBar struct {
    transform Transform
    maxHealth float64
    currentHealth float64
    bloodOverlayPulse float64
}

func NewHealthBar() *HealthBar {
    return &HealthBar{
        transform: Transform{
            x: 200,
            y: 1000,
            z: 100,
            width: 300,
            height: 30,
        },
        maxHealth: 300,
        currentHealth: 300,
    }
}

func (hb *HealthBar) Update() {
    hb.currentHealth = float64(player.health)
    hb.bloodOverlayPulse += 0.02
}

func (hb *HealthBar) Draw(screen *ebiten.Image) {
    // Draw the background bar
    bgColor := color.White
    fgColor := color.RGBA{255, 0, 0, 200} // Red foreground

    drawAbsoluteRect(screen, hb.transform, bgColor)

    foregroundWidth := hb.transform.width

    // Ensure the foreground width is not negative
    if hb.currentHealth <= 0 {
        foregroundWidth = 0
    }

    // Draw the foreground bar
    foregroundTransform := hb.transform
    foregroundTransform.width = foregroundWidth

    border := 5.0

    if foregroundWidth > 0 {
        foregroundTransform.width -= border * 2
        foregroundTransform.height -= border * 2
        foregroundTransform.x -= hb.transform.width/2 - border

        op := defaultImageOptions()
        op.Anchor.x = -foregroundTransform.width / 2
        op.Scale.x = hb.currentHealth / hb.maxHealth
        foreground := getCachedRect(int(foregroundTransform.width), int(foregroundTransform.height), fgColor)

        drawAbsoluteImageWithOptions(screen, foreground, foregroundTransform, op)
    }

    hb.DrawBloodOverlay(screen)
}

func (hb *HealthBar) DrawBloodOverlay(screen *ebiten.Image) {
    op := defaultImageOptions()
    op.OriginalImageSize = true
    op.Alignment = AlignCenter
    op.Alpha = (hb.maxHealth - hb.currentHealth) / hb.maxHealth
    op.Alpha *= 255

    op.Alpha -= math.Max(math.Sin(hb.bloodOverlayPulse) * 50, 0)

    drawAbsoluteImageWithOptions(screen, getCachedImage("sprites/blood/overlay"), Transform{
        // width: 1920,
        // height: 1080,
        z: 1000,
    }, op)
}

func (hb *HealthBar) GetTransform() Transform {
    return hb.transform
}

func (hb *HealthBar) SetTransform(transform Transform) {
    hb.transform = transform
}
