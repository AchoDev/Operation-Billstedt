package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func createEnemy(x, y int, enemyType EnemyType) *Enemy {

    var gun Gun

    switch enemyType {
    case EnemyTypeEvren:
        gun = &Pistol{}
    case EnemyTypeEmran:
        gun = &Shotgun{}
    case EnemyTypeNick:
        gun = &Rifle{}
    }


    enemy := Enemy{
        transform: Transform{
            x:      float64(x),
            y:      float64(y),
            width:  50,
            height: 50,
        },
        gun: createGun(gun, true),
        enemyType: enemyType,
    }

    return &enemy
}

type EnemyType int

const (
    EnemyTypeEvren EnemyType = iota
    EnemyTypeEmran
    EnemyTypeNick
)

type Enemy struct {
    transform Transform
    gun Gun
    enemyType EnemyType
}

func (enemy *Enemy) Update() {
    
    enemy.transform.rotation = math.Atan2(
        enemy.transform.y - player.transform.y,
        enemy.transform.x - player.transform.x,
    ) + math.Pi

    direction := Vector2{
        x: math.Cos(enemy.transform.rotation),
        y: math.Sin(enemy.transform.rotation),
    }

    var speed float64

    switch enemy.enemyType {
    case EnemyTypeEvren:
        speed = 3.5
    case EnemyTypeEmran:
        speed = 2
    case EnemyTypeNick:
        speed = 5
    }

    enemy.transform.x += direction.x * speed
    enemy.transform.y += direction.y * speed

    distance := math.Sqrt(
        math.Pow(enemy.transform.x-player.transform.x, 2) +
        math.Pow(enemy.transform.y-player.transform.y, 2),
    )

    var attackDistance float64

    switch enemy.enemyType {
    case EnemyTypeEvren:
        attackDistance = 500
    case EnemyTypeEmran:
        attackDistance = 200
    case EnemyTypeNick:
        attackDistance = 750
    }

    if distance < attackDistance {
        enemy.gun.Shoot(&enemy.transform)
    }

}
func (enemy *Enemy) Draw(screen *ebiten.Image) {

    var col color.RGBA

    switch enemy.enemyType {
    case EnemyTypeEvren:
        col = color.RGBA{0, 255, 0, 255}
    case EnemyTypeEmran:
        col = color.RGBA{255, 0, 255, 255}
    case EnemyTypeNick:
        col = color.RGBA{0, 0, 255, 255}
    }

    drawRotatedRect(
        screen,
        enemy.transform.x,
        enemy.transform.y,
        enemy.transform.width,
        enemy.transform.height,
        enemy.transform.rotation,
        col,
    )

    textX := enemy.transform.x - enemy.transform.width/2
    textY := enemy.transform.y - enemy.transform.height


    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", enemy.gun.GetCooldownTimer()), int(textX), int(textY))
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", enemy.transform.rotation), int(textX), int(textY - 20))
}

func (enemy *Enemy) GetTransform() Transform {
    return enemy.transform
}