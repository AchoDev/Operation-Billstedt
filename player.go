package main

import (
	"fmt"
	"math"

	"strings"

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

func CreatePlayer() *Player {
	p := Player{
		transform: Transform{
			x:      650,
			y:      0,
			width:  30,
			height: 30,
		},
		velocity:   Vector2{0, 0},
		shooting:   false,
		sprites: map[string]*ebiten.Image{
			"minigun": loadImage("assets/leo/minigun.png"),
			"rifle":   loadImage("assets/leo/rifle.png"),
			"pistol":  loadImage("assets/leo/pistol.png"),
			"shotgun": loadImage("assets/leo/shotgun.png"),
		},
		health:   100,
	}

	p.guns = []GunBase{
		*NewGun(pistolStats, &p),
		// *NewGun("Shotgun", 3000, &p, ShotgunShoot),
		*NewGun(shotgunStats, &p),
		*NewGun(rifleStats, &p),
		*NewGun(minigunStats, &p),
	}

	p.currentGun = p.guns[0]

	return &p
}

var invincible bool = false


type Player struct {
	transform  Transform
	velocity   Vector2
	shooting   bool
	currentGun GunBase
	guns       []GunBase
	sprites    map[string]*ebiten.Image
	health    int
}

func (player *Player) Update() {
	move(player)

	mousePos := getMousePosition()

	angle := math.Atan2(
		mousePos.y-player.transform.y,
		mousePos.x-player.transform.x,
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
		player.currentGun = player.guns[0]
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		player.currentGun = player.guns[1]
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		player.currentGun = player.guns[2]
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		player.currentGun = player.guns[3]
	}

	if isKeyJustPressed(ebiten.Key6) {
		gameObjects = append(gameObjects, createEnemy(500, 100, EnemyTypeNick))
	}

	if isKeyJustPressed(ebiten.Key7) {
		gameObjects = append(gameObjects, createEnemy(500, 100, EnemyTypeEvren))
	}

	if isKeyJustPressed(ebiten.Key8) {
		gameObjects = append(gameObjects, createEnemy(500, 100, EnemyTypeEmran))
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		camera.zoom += 0.01
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		camera.zoom -= 0.01
	}


	// if ebiten.IsKeyPressed(ebiten.KeyRight) {
	// 	casingPoint.y += 1
	// }
	
	// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
	// 	casingPoint.y -= 1
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	casingPoint.x -= 1
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	casingPoint.x += 1
	// }

	// fmt.Println("Casing Point:", casingPoint.x, casingPoint.y)
}

func move(player *Player) {
	acceleration := 4.0
	movement := Vector2{0, 0}
	max_vel := 7.5

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

	
	// player.velocity = clampVector(player.velocity, -max_vel, max_vel)
	
	if math.Abs(player.velocity.x) > max_vel {
		movement.x = 0
	}
	
	if math.Abs(player.velocity.y) > max_vel {
		movement.y = 0
	}
		
	player.velocity.x -= movement.x * acceleration
	player.velocity.y -= movement.y * acceleration

	player.transform.x += player.velocity.x
	player.transform.y += player.velocity.y
		
	if movement.x == 0 {
		player.velocity.x /= 1.15
	}
	if movement.y == 0 {
		player.velocity.y /= 1.15
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	sprite := player.sprites[strings.ToLower(player.currentGun.Name())]

	if sprite != nil {

		offset := Vector2{
			-110,
			500,
		}

		op := defaultImageOptions()
		op.Anchor = offset
		op.Scale = 4

		tr := player.GetTransform()
		tr.rotation += math.Pi / 2

		drawImageWithOptions(
			screen,
			sprite,
			tr,
			op,
		)
	} else {
		fmt.Println("Sprite not found for gun:", strings.ToLower(player.currentGun.Name()), player.sprites,)
	}
}

func (player *Player) GetTransform() Transform {
	if player == nil {
		return Transform{}
	}
	return player.transform
}

func (player *Player) SetTransform(transform Transform) {
	player.transform = transform
}
