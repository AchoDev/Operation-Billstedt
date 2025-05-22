package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

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
		drag: 1.15,
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
	drag float64
	shooting   bool
	currentGun GunBase
	guns       []GunBase
	sprites    map[string]*ebiten.Image
	health    int
	dashing bool
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
		addGameObject(createEnemy(500, 100, EnemyTypeNick))
	}

	if isKeyJustPressed(ebiten.Key7) {
		addGameObject(createEnemy(500, 100, EnemyTypeEvren))
	}

	if isKeyJustPressed(ebiten.Key8) {
		addGameObject(createEnemy(500, 100, EnemyTypeEmran))
	}

	if isKeyJustPressed(ebiten.KeyE) {
		player.dashing = true

		direction := Vector2{}

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			direction.x = -1
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			direction.x = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			direction.y = -1
		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			direction.y = 1
		}

		player.velocity.x = direction.x * 30
		player.velocity.y = direction.y * 30

		angle := math.Atan2(direction.y, direction.x)

		camera.Shake(angle, 5.0)

		player.drag = 1

		go func () {
			for i := 0; i < 12; i++ {
				spriteRender := SpriteRenderer{
					transform: player.transform,
					sprite: player.GetSprite(),
				}
				spriteRender.options = defaultImageOptions()
				spriteRender.options.Alpha = 0.75
				spriteRender.options.ColorScale = color.RGBA{255, 255, 255, 100}
				spriteRender.options.ScaleColor = true
				spriteRender.options.Scale = 4

				spriteRender.transform.rotation += math.Pi / 2

				addGameObject(&spriteRender)

				go func() {
					steps := 60.0
					duration := 0.2
					for j := 0; j < int(steps); j++ {
						spriteRender.options.Alpha -= 0.75 / steps
						pausableSleep(time.Second / time.Duration(steps / duration))
					}

					removeGameObject(&spriteRender)
				}()

				pausableSleep(25 * time.Millisecond)
			}
		}()
			
		go func() {
			pausableSleep(250 * time.Millisecond)
			player.dashing = false
			player.drag = 1.15
		}()
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
		player.velocity.x /= player.drag
	}
	if movement.y == 0 {
		player.velocity.y /= player.drag
	}
}

func (player *Player) GetSprite() *ebiten.Image {
	sprite := player.sprites[strings.ToLower(player.currentGun.Name())]

	if sprite != nil {
		return sprite
	}

	return nil
}

func (player *Player) Draw(screen *ebiten.Image) {
	sprite := player.GetSprite()

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

		if player.dashing {
			op.ColorScale = color.RGBA{255, 255, 255, 255}
			op.ScaleColor = true
		}

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
