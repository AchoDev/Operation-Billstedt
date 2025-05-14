package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
		gun:       createGun(gun, true),
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
	transform   Transform
	gun         Gun
	enemyType   EnemyType
	currentPath []Vector2
	currentGoal Vector2
}

var pathFindingGridSize = 10

func (enemy *Enemy) Update() {
	colliders := getGameobjectsOfType[*Collider]()

	// path will be set below based on the currentGoal

	var path []Vector2

	if math.Round(enemy.currentGoal.x) == math.Round(player.transform.x) && math.Round(enemy.currentGoal.y) == math.Round(player.transform.y) {
		path = enemy.currentPath
	} else {
		path = runPathfindingAlgorithm(enemy.transform, player.transform, colliders, pathFindingGridSize, Vector2{2000, 2000})
		fmt.Println(enemy.currentGoal, player.transform.x, player.transform.y)
	}

	enemy.currentGoal = Vector2{
		x: player.transform.x,
		y: player.transform.y,
	}
	enemy.currentPath = path
	var target Vector2
	enemyGridPos := Vector2{
		x: float64(int(enemy.transform.x) / pathFindingGridSize),
		y: float64(int(enemy.transform.y) / pathFindingGridSize),
	}

	if len(path) > 1 {

		pathPos := Vector2{
			x: float64(int(path[0].x) / pathFindingGridSize),
			y: float64(int(path[0].y) / pathFindingGridSize),
		}

		fmt.Println(pathPos, enemyGridPos)
		if enemyGridPos.x == pathPos.x && enemyGridPos.y == pathPos.y {
			fmt.Println("deleting latest link")
			enemy.currentPath = enemy.currentPath[1:]
			path = enemy.currentPath
		}

		target = path[0]

		// fmt.Println(target)

	} else {
		target = Vector2{
			x: player.transform.x,
			y: player.transform.y,
		}
	}

	enemy.transform.rotation = math.Atan2(
		(enemyGridPos.y*float64(pathFindingGridSize))-target.y,
		(enemyGridPos.x*float64(pathFindingGridSize))-target.x,
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

	for _, point := range enemy.currentPath {
		vector.DrawFilledRect(screen, float32(point.x-float64(pathFindingGridSize/2)), float32(point.y-float64(pathFindingGridSize/2)), float32(pathFindingGridSize), float32(pathFindingGridSize), color.RGBA{255, 0, 0, 50}, true)
	}

	textX := enemy.transform.x - enemy.transform.width/2
	textY := enemy.transform.y - enemy.transform.height

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", enemy.gun.GetCooldownTimer()), int(textX), int(textY))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", enemy.transform.rotation), int(textX), int(textY-20))
}

func (enemy *Enemy) GetTransform() Transform {
	return enemy.transform
}
