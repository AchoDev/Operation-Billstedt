package main

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	transform Transform
	angle     float64
	speed     float64
	fromEnemy bool
}

func (bullet *Bullet) Update() {
	bullet.transform.x += bullet.speed * math.Cos(bullet.angle)
	bullet.transform.y += bullet.speed * math.Sin(bullet.angle)

	for i := 0; i < len(gameObjects); i++ {
		gameObject := gameObjects[i]
		var target GameObject

		hitboxMultiplier := 1.0

		if bullet.fromEnemy {
			if player, ok := gameObject.(*Player); ok && !invincible {
				target = player
			} else {
				continue
			}
		} else {
			if enemy, ok := gameObject.(*Enemy); ok {
				target = enemy
				hitboxMultiplier = 3
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
				Angle:  bullet.angle,
			},
			Rect{
				Center: Vector2{tr.x, tr.y},
				Width:  tr.width * hitboxMultiplier,
				Height: tr.height * hitboxMultiplier,
				Angle:  tr.rotation,
			},
		) {
			if player, ok := target.(*Player); ok {
				player.health -= 10
				if player.health <= 0 {
					// Remove the gameObject from the list
					gameObjects = append(gameObjects[:i], gameObjects[i+1:]...)
				}
			} else {
				gameObjects = append(gameObjects[:i], gameObjects[i+1:]...)
			}

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

	for _, collider := range currentLevel.GetColliders() {
		if RotatedRectsColliding(
			Rect{
				Center: Vector2{bullet.transform.x, bullet.transform.y},
				Width:  bullet.transform.width,
				Height: bullet.transform.height,
				Angle:  bullet.angle,
			},
			Rect{
				Center: Vector2{
					collider.X * 100,
					collider.Y * 100,
				},
				Width:  collider.Width * 100,
				Height: collider.Height * 100,
			},
		) {
			for i := 0; i < len(gameObjects); i++ {
				if gameObjects[i] == bullet {
					gameObjects = append(gameObjects[:i], gameObjects[i+1:]...)
					break
				}
			}
			break
		}
	}
}

func CreateBullet(transform *Transform, gun *GunBase) *Bullet {
	bullet := Bullet{
		transform: Transform{
			x:      transform.x,
			y:      transform.y,
			width:  25,
			height: 10,
		},
		angle:     transform.rotation,
		speed:     15,
		fromEnemy: gun.isEnemy,
	}

	bullet.angle += float64(rand.IntN(10) - 5) * gun.spread

	bullet.transform.x += gun.offset.x
	bullet.transform.y += gun.offset.y
	bullet.transform.RotateAround(transform.rotation, transform.GetPosition())

	if gun.hasCasing {
		tr := Transform{
			x:      gun.carrier.GetTransform().x,
			y:      gun.carrier.GetTransform().y,
			rotation: bullet.angle,
		}

		speed := Vector2{}

		player, ok := gun.carrier.(*Player)
		if ok {
			speed = player.velocity
		}

		if enemy, ok := gun.carrier.(*Enemy); ok {
			speed = enemy.velocity
		}

		casing := CreateCasing(&tr, gun.casingPoint, speed)
		gameObjects = append(gameObjects, casing)
	}

	return &bullet
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
	bullet.transform.rotation = bullet.angle + math.Pi/2
	sprite := getCachedImage("sprites/bullet")
	op := defaultImageOptions()
	op.OriginalImageSize = true
	op.Scale = 0.04

	drawImageWithOptions(screen, sprite, bullet.transform, op)
}

func (bullet *Bullet) GetTransform() Transform {
	return bullet.transform
}

func (bullet *Bullet) SetTransform(transform Transform) {
	bullet.transform = transform
}


type Casing struct {
	transform Transform
	angularVelocity float64
	velocity Vector2
	spawnHeight float64
	spawnTime time.Time
}

func (casing *Casing) Update() {
	casing.transform.rotation += casing.angularVelocity
	
	startPos := casing.transform.GetPosition()

	casing.transform.y += casing.velocity.y
	casing.transform.x += casing.velocity.x

	casing.spawnHeight -= 1

	if casing.spawnHeight < 0 {
		casing.velocity.y *= 0.9
		casing.velocity.x *= 0.9
		casing.angularVelocity *= 0.9
	}

	xCollided, yCollided := checkCollisions(&casing.transform, startPos)

	if xCollided {
		casing.velocity.x = -casing.velocity.x
		casing.angularVelocity *= -1

		casing.velocity.x *= 0.3
		casing.angularVelocity *= 0.3
	}

	if yCollided {
		casing.velocity.y = -casing.velocity.y
		casing.angularVelocity *= -1
		casing.velocity.y *= 0.3
		casing.angularVelocity *= 0.3
	}

	if time.Since(casing.spawnTime) > 1*time.Second {
		removeGameObject(casing)
	}
}

func CreateCasing(transform *Transform, offset, currentSpeed Vector2) *Casing {

	xVel, yVel := currentSpeed.x, currentSpeed.y

	xVel += -math.Sin(transform.rotation) * 10
	yVel += math.Cos(transform.rotation) * 10

	xVel += float64(rand.IntN(10) - 5) * 0.25
	yVel += float64(-rand.IntN(10) - 5) * 0.25

	casing := Casing{
		transform: Transform{
			x:      transform.x + offset.x,
			y:      transform.y + offset.y,
			width:  25,
			height: 10,
		},
		angularVelocity: float64(rand.IntN(10) - 5) * 0.5,
		spawnTime: time.Now(),
		velocity: Vector2{
			x: xVel,
			y: yVel,
		},
		spawnHeight: 10,
	}

	casing.transform.RotateAround(transform.rotation, transform.GetPosition())

	return &casing
}

func (casing *Casing) Draw(screen *ebiten.Image) {
	sprite := getCachedImage("sprites/casing")
	op := defaultImageOptions()
	op.OriginalImageSize = true
	op.Scale = 0.04

	drawImageWithOptions(screen, sprite, casing.transform, op)
}
func (casing *Casing) GetTransform() Transform {
	return casing.transform
}

func (casing *Casing) SetTransform(transform Transform) {
	casing.transform = transform
}
