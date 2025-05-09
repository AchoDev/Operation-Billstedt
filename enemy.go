package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func createEnemy(x, y int) *Enemy {
    return &Enemy{
        transform: Transform{
            x:      float64(x),
            y:      float64(y),
            width:  50,
            height: 50,
        },
        gun: &Pistol{},
    }
}

type Enemy struct {
    transform Transform
    gun Gun
}

func (enemy *Enemy) Update() {
    
    enemy.transform.rotation = math.Atan2(
        enemy.transform.y - player.transform.y,
        enemy.transform.x - player.transform.x,
    )

    direction := Vector2{
        x: math.Cos(enemy.transform.rotation),
        y: math.Sin(enemy.transform.rotation),
    }

    enemy.transform.x -= direction.x * 2
    enemy.transform.y -= direction.y * 2

    distance := math.Sqrt(
        math.Pow(enemy.transform.x-player.transform.x, 2) +
        math.Pow(enemy.transform.y-player.transform.y, 2),
    )

    if distance < 100 {
        enemy.gun.Shoot(&enemy.transform)
    }

}
func (enemy *Enemy) Draw(screen *ebiten.Image) {
    drawRotatedRect(
        screen,
        enemy.transform.x,
        enemy.transform.y,
        enemy.transform.width,
        enemy.transform.height,
        enemy.transform.rotation,
        color.RGBA{0, 255, 0, 255},
    )
}

func (enemy *Enemy) GetTransform() Transform {
    return enemy.transform
}