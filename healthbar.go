package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type HealthBar struct {
    transform Transform
    maxHealth float64
    currentHealth float64
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
        maxHealth: 100,
        currentHealth: 100,
    }
}

func (hb *HealthBar) Update() {
    hb.currentHealth = float64(player.health)
}

func (hb *HealthBar) Draw(screen *ebiten.Image) {
    // Draw the background bar
    bgColor := color.White
    fgColor := color.RGBA{255, 0, 0, 200} // Red foreground

    drawAbsoluteRect(screen, hb.transform, bgColor)

    // Calculate the width of the foreground bar based on current health
    healthPercentage := hb.currentHealth / hb.maxHealth
    foregroundWidth := hb.transform.width * healthPercentage

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
        foreground := getCachedRect(int(foregroundTransform.width), int(foregroundTransform.height), fgColor)

        drawAbsoluteImageWithOptions(screen, foreground, foregroundTransform, op)
    }
}

func (hb *HealthBar) GetTransform() Transform {
    return hb.transform
}

func (hb *HealthBar) SetTransform(transform Transform) {
    hb.transform = transform
}
