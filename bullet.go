package main

import (
	"image/color"
	"math"

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

		if bullet.fromEnemy {
			if player, ok := gameObject.(*Player); ok && !invincible {
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
				Angle:  bullet.angle,
			},
			Rect{
				Center: Vector2{tr.x, tr.y},
				Width:  tr.width,
				Height: tr.height,
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

	bullet.transform.x += gun.offset.x
	bullet.transform.y += gun.offset.y
	bullet.transform.RotateAround(transform.rotation, transform.GetPosition())

	return &bullet
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
	bullet.transform.rotation = bullet.angle
	drawRect(
		screen,
		bullet.transform,
		color.RGBA{255, 238, 66, 255},
	)
}

func (bullet *Bullet) GetTransform() Transform {
	return bullet.transform
}

func (bullet *Bullet) SetTransform(transform Transform) {
	bullet.transform = transform
}
