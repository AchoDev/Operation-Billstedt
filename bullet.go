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
    fromEnemy bool
}

func (bullet *Bullet) Update() {
    bullet.transform.x += bullet.speed * math.Cos(bullet.angle)
    bullet.transform.y += bullet.speed * math.Sin(bullet.angle)

    for i := 0; i < len(gameObjects); i++ {
        gameObject := gameObjects[i]
        var target GameObject

        if bullet.fromEnemy {
            if player, ok := gameObject.(*Player); ok {
                target = player
            } else {
                continue
            }
        } else {

            if enemy, ok := gameObject.(*Enemy); ok {
                target = enemy
            } else {
                continue
            }
            
        }

        tr := target.GetTransform()
        if RotatedRectsColliding(
            Rect{
                Center: Vector2{bullet.transform.x, bullet.transform.y},
                Width:  bullet.transform.width,
                Height: bullet.transform.height,
                Angle: bullet.angle,
            },
            Rect{
                Center: Vector2{tr.x, tr.y},
                Width:  tr.width,
                Height: tr.height,
                Angle: tr.rotation,
            },
        ) {
            // Remove the gameObject from the list
            gameObjects = append(gameObjects[:i], gameObjects[i+1:]...)

            // Remove the bullet from the list
            for j := 0; j < len(gameObjects); j++ {
                if gameObjects[j] == bullet {
                    gameObjects = append(gameObjects[:j], gameObjects[j+1:]...)
                    break
                }
            }

            break
        }
    }
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
    drawRotatedRect(
        screen,
        bullet.transform.x,
        bullet.transform.y,
        bullet.transform.width,
        bullet.transform.height,
        bullet.angle,
        color.RGBA{255, 238, 66, 255},
    )
}

func (bullet *Bullet) GetTransform() Transform {
    return bullet.transform
}
