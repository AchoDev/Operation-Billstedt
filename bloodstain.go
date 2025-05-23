package main

import (
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bloodstain struct {
    transform Transform
    sprite *ebiten.Image
    currentHeight float64
    velocity Vector2
    collidedWithWall bool
    scale float64
}

func NewBloodstain(position Vector2, angle float64) *Bloodstain {
    sprite := getCachedImage("sprites/blood/splash-" + strconv.Itoa(rand.IntN(5) + 1))
    angle += (rand.Float64() - 1) * 0.5
    
    velocity := Vector2{
        x: math.Cos(angle),
        y: math.Sin(angle),
    }

    scale := 0.1 + rand.Float64()*(0.15-0.1)
    speed := 6 + rand.Float64()*(10-6)
    height := 10 + rand.Float64()*(15-10)

    velocity.normalize()
    velocity.Multiply(speed)

    transform := Transform{
        x: position.x,
        y: position.y,
        z: 0.1,
        width: 10,
        height: 10,
        rotation: angle,
    }

    return &Bloodstain{
        transform: transform,
        sprite: sprite,
        currentHeight: height,
        velocity: velocity,
        scale: scale,
    }
}

func (b *Bloodstain) Update() {
    
    startPos := Vector2{
        x: b.transform.x,
        y: b.transform.y,
    }

    b.transform.x += b.velocity.x
    b.transform.y += b.velocity.y
    b.currentHeight -= 1

    if b.currentHeight <= 0 {
        b.velocity.x = 0
        b.velocity.y = 0
    }

    xCollision, yCollision := false, false
    if !b.collidedWithWall {
        xCollision, yCollision = checkCollisions(&b.transform, startPos)
    }


    if xCollision || yCollision {
        b.currentHeight = 0
        b.collidedWithWall = true
        
        rot := 0.0
        if xCollision {
            rot = math.Pi * 3 / 2
            if b.velocity.x < 0 {
                rot += math.Pi
            }
        } else {
            rot = 0
            if b.velocity.y < 0 {
                rot += math.Pi
            }
        }

        b.sprite = getCachedImage("sprites/blood/splash-wall-" + strconv.Itoa(rand.IntN(2) + 1))
        b.transform.rotation = rot
        b.velocity.x = 0
        b.velocity.y = 0
    }
}

func (b *Bloodstain) Draw(screen *ebiten.Image) {
    s := b.sprite
    op := defaultImageOptions()
    op.OriginalImageSize = true
    op.Scale = b.scale

    if b.currentHeight > 0 {
        op.Scale = b.scale * 0.5
        s = getCachedImage("sprites/blood/ball")
    }

    drawImageWithOptions(screen, s, b.transform, op)
}


func (b *Bloodstain) GetTransform() Transform {
    return b.transform
}

func (b *Bloodstain) SetTransform(transform Transform) {
    b.transform = transform
}