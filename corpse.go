package main

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Corpse struct {
    transform Transform
    velocity Vector2
    sprite *ebiten.Image
}

func NewCorpse(position Vector2, bulletAngle float64, enemyName string) *Corpse {
    num := rand.Intn(2) + 1
    sprite := getCachedImage("sprites/corpses/" + enemyName + "-" + strconv.Itoa(num))

    velocity := Vector2{
        x: math.Cos(bulletAngle),
        y: math.Sin(bulletAngle),
    }

    velocity.normalize()

    velocity.Multiply(10)

    return &Corpse{
        transform: Transform{
            x:      position.x,
            y:      position.y,
            z: 0.1,
            width:  200,
            height: 200,
            rotation: bulletAngle + math.Pi / 2,
        },
        velocity: velocity,
        sprite:   sprite,
    }
}

func (c *Corpse) Update() {
    c.transform.x += c.velocity.x
    c.transform.y += c.velocity.y

    c.velocity.x *= 0.8
    c.velocity.y *= 0.8
}

func (c *Corpse) Draw(screen *ebiten.Image) {
    drawImage(screen, c.sprite, c.transform)
}


func (c *Corpse) GetTransform() Transform {
    return c.transform
}

func (c *Corpse) SetTransform(transform Transform) {
    c.transform = transform
}
