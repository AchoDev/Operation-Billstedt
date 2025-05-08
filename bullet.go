package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
    transform Transform
    angle float64
    speed float64
}

func (bullet *Bullet) Update() {
    bullet.transform.x += bullet.speed * math.Cos(bullet.angle)
    bullet.transform.y += bullet.speed * math.Sin(bullet.angle)
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
    drawRotatedRect(
        screen,
        bullet.transform.x,
        bullet.transform.y,
        bullet.transform.width,
        bullet.transform.height,
        bullet.transform.rotation,
        color.RGBA{255, 0, 0, 255},
    )
}