package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func clamp(num float64, min float64, max float64) float64 {
	if num < min {
		return min
	} else if num > max {
		return max
	}

	return num
}

func clampVector(vector Vector2, min float64, max float64) Vector2 {
	return Vector2{
		x: clamp(vector.x, min, max),
		y: clamp(vector.y, min, max),
	}
}

func CreatePlayer() Player {
	return Player{
		transform: Transform{
			x:      500,
			y:      500,
			width:  30,
			height: 30,
		},
		velocity: Vector2{0, 0},
        shooting: false,
		currentGun: &Pistol{},
	}
}

type Player struct {
	transform Transform
	velocity  Vector2
    shooting bool
	currentGun Gun
}

func (player *Player) Update() {
	move(player)

	mouseX, mouseY := ebiten.CursorPosition()

	angle := math.Atan2(
		float64(mouseY)-player.transform.y,
		float64(mouseX)-player.transform.x,
	)

	player.transform.rotation = angle

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !player.shooting {
		player.shooting = true
        player.currentGun.Shoot(&player.transform)
	}

    if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        player.shooting = false
    }


	if ebiten.IsKeyPressed(ebiten.Key1) {
		player.currentGun = createGun(&Pistol{}, false)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		player.currentGun = createGun(&Shotgun{}, false)
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		player.currentGun = createGun(&Rifle{}, false)
	}

	if isKeyJustPressed(ebiten.Key7) {
		gameObjects = append(gameObjects, createEnemy(100, 100))
	}

}

func move(player *Player) {
	acceleration := 8.0
	movement := Vector2{0, 0}
	max_vel := 10.0

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		movement.x = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		movement.x = -1
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		movement.y = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		movement.y = -1
	}

	movement.normalize()

	player.velocity.x -= movement.x * acceleration
	player.velocity.y -= movement.y * acceleration

	player.velocity = clampVector(player.velocity, -max_vel, max_vel)

	player.transform.x += player.velocity.x
	player.transform.y += player.velocity.y

	if movement.x == 0 {
		player.velocity.x /= 1.5
	}
	if movement.y == 0 {
		player.velocity.y /= 1.5
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	drawRotatedRect(
		screen,
		player.transform.x-camera.x,
		player.transform.y-camera.y,
		player.transform.width,
		player.transform.height,
		player.transform.rotation,
		color.Color(color.RGBA{255, 0, 0, 255}),
	)
}

func (player *Player) GetTransform() Transform {
	if player == nil {
		return Transform{}
	}
	return player.transform
}
