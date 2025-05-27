package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Particle struct {
    transform Transform
    velocity Vector2
    torque float64
    spawned time.Time
    lifetime time.Duration
}

type HitParticle struct {
    transform Transform
    particles []Particle
}

func (hitParticle *HitParticle) Update() {
    for i := len(hitParticle.particles) - 1; i >= 0; i-- {
        particle := &hitParticle.particles[i]
        particle.transform.x += particle.velocity.x
        particle.transform.y += particle.velocity.y

        // Remove particles that are too far away
        if time.Since(particle.spawned) > particle.lifetime {
            particle.transform.width -= 5
            particle.transform.height -= 5

            if particle.transform.width <= 0 || particle.transform.height <= 0 {
                hitParticle.particles = append(hitParticle.particles[:i], hitParticle.particles[i+1:]...)
            }
        }
    }    
}

func NewHitParticle(position Vector2, angle float64) *HitParticle {
    particles := []Particle{}
    for i := 0; i < 10; i++ {

        angle += math.Pi

        randomAngle := angle + (rand.Float64()-0.5)*(math.Pi / 2)

        velocity := Vector2{
            x: math.Cos(randomAngle),
            y: math.Sin(randomAngle),
        }
        velocity.normalize()
        velocity.Multiply(10 + rand.Float64()*5)

        torque := (rand.Float64() - 0.5) * 0.1

        particles = append(particles, Particle{
            transform: Transform{
                x:      position.x,
                y:      position.y,
                z:      0.1,
                width:  10,
                height: 10,
                rotation: angle + (rand.Float64()-0.5)*math.Pi/4,
            },
            torque:  torque,
            velocity: velocity,
            spawned: time.Now(),
            lifetime: time.Duration(rand.Intn(20)+10) * time.Millisecond,
        })
    }

    return &HitParticle{
        transform: Transform{
            x:      position.x,
            y:      position.y,
            z:      0.1,
            width:  50,
            height: 50,
        },
        particles: particles,
    }
}

func (hitParticle *HitParticle) Draw(screen *ebiten.Image) {
    rect := getCachedRect(50, 50, color.RGBA{255, 195, 0, 255})
    for _, particle := range hitParticle.particles {
        drawImage(screen, rect, particle.transform)
    }
}

func (hitParticle *HitParticle) GetTransform() Transform {
    return hitParticle.transform
}

func (hitParticle *HitParticle) SetTransform(transform Transform) {
    hitParticle.transform = transform
}