package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func createEnemy(x, y int, enemyType EnemyType) *Enemy {

	var gun *GunBase

	switch enemyType {
	case EnemyTypeEvren:
		gun = NewGun(pistolStats, nil)
		gun.cooldown = 750
	case EnemyTypeEmran:
		gun = NewGun(rifleStats, nil)
		gun.cooldown = 2000
	case EnemyTypeNick:
		gun = NewGun(shotgunStats, nil)
		gun.cooldown = 2500
	}

	enemy := Enemy{
		transform: Transform{
			x:      float64(x),
			y:      float64(y),
			z: 0.5,
			width:  30,
			height: 30,
		},
		gun:       *gun,
		enemyType: enemyType,
	}

	enemy.gun.carrier = &enemy
	enemy.gun.isEnemy = true

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
	gun         GunBase
	enemyType   EnemyType
	currentPath []Vector2
	currentGoal Vector2
	velocity    Vector2
	pathChan    <-chan []Vector2
}

var pathFindingGridSize = 20

func runPathfindingAlgorithmAsync(start, end Transform, colliders []*Collider, gridSize int, worldSize Vector2) <-chan []Vector2 {
	resultChan := make(chan []Vector2, 1)

	go func() {
		result := runPathfindingAlgorithm(start, end, colliders, gridSize, worldSize)
		resultChan <- result
	}()

	return resultChan
}

func (enemy *Enemy) Update() {
	colliders := getGameobjectsOfType[*Collider]()

	for _, col := range currentLevel.GetColliders() {
		collider := &Collider{
			transform: Transform{
				x:      col.X * 100,
				y:      col.Y * 100,
				width:  col.Width * 100,
				height: col.Height * 100,
			},
		}
		colliders = append(colliders, collider)
	}

	for _, e := range getGameobjectsOfType[*Enemy]() {
		if enemy == e {
			continue
		}

		collider := &Collider{
			transform: Transform{
				x:      e.transform.x,
				y:      e.transform.y,
				width:  e.transform.width,
				height: e.transform.height,
			},
		}

		colliders = append(colliders, collider)
	}

	for _, bullet := range getGameobjectsOfType[*Bullet]() {
		if bullet.fromEnemy {
			continue
		}

		predictionSteps := 10
		stepSize := 50.0
		colliderWidth := 10.0
		colliderHeight := 10.0

		// Predict the bullet's trajectory
		for i := 1; i <= predictionSteps; i++ {
			angle := bullet.angle
			predictedX := bullet.transform.x + math.Cos(angle)*bullet.speed*stepSize*float64(i)
			predictedY := bullet.transform.y + math.Sin(angle)*bullet.speed*stepSize*float64(i)

			// Create a small collider at the predicted position
			collider := &Collider{
				transform: Transform{
					x:      predictedX,
					y:      predictedY,
					width:  colliderWidth,
					height: colliderHeight,
				},
			}

			// Add the collider to the list
			colliders = append(colliders, collider)
		}
	}

	path := enemy.currentPath

	if (enemy.currentGoal.x != player.transform.x || enemy.currentGoal.y != player.transform.y) && enemy.pathChan == nil {
		enemy.currentGoal = Vector2{
			x: player.transform.x,
			y: player.transform.y,
		}
		enemy.pathChan = runPathfindingAlgorithmAsync(enemy.transform, player.transform, colliders, pathFindingGridSize, Vector2{14000, 14000})
	}

	select {
	case newPath := <-enemy.pathChan:
		enemy.pathChan = nil
		enemy.currentPath = newPath
		if len(newPath) != 0 {
			enemy.currentPath = enemy.currentPath[1:]
		}

		path = enemy.currentPath

	default:
	}

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
		if enemyGridPos.x == pathPos.x && enemyGridPos.y == pathPos.y {
			enemy.currentPath = enemy.currentPath[1:]
			path = enemy.currentPath
		}

		forwardVec := Vector2{
			x: path[0].x - enemy.transform.x,
			y: path[0].y - enemy.transform.y,
		}
		dotProduct := forwardVec.x*enemy.velocity.x + forwardVec.y*enemy.velocity.y

		if dotProduct < 0 && len(enemy.currentPath) > 1 {
			enemy.currentPath = enemy.currentPath[1:]
			path = enemy.currentPath
		}

		target = path[0]
	} else {
		target = Vector2{
			x: player.transform.x,
			y: player.transform.y,
		}
	}

	targetRotation := math.Atan2(
		(enemyGridPos.y*float64(pathFindingGridSize))-target.y,
		(enemyGridPos.x*float64(pathFindingGridSize))-target.x,
	) + math.Pi

	// rotationDiff := targetRotation - enemy.transform.rotation

	// // Normalize the angle to the range [-Pi, Pi]
	// for rotationDiff > math.Pi {
	// 	rotationDiff -= 2 * math.Pi
	// }
	// for rotationDiff < -math.Pi {
	// 	rotationDiff += 2 * math.Pi
	// }

	// enemy.transform.rotation += rotationDiff * 0.1 // Smoothing factor

	enemy.transform.rotation = math.Atan2(
		enemy.transform.y-player.transform.y,
		enemy.transform.x-player.transform.x,
	) + math.Pi

	direction := Vector2{
		x: math.Cos(targetRotation),
		y: math.Sin(targetRotation),
	}

	var speed float64

	switch enemy.enemyType {
	case EnemyTypeEvren:
		speed = 5
	case EnemyTypeEmran:
		speed = 2.5
	case EnemyTypeNick:
		speed = 6
	}

	enemy.velocity.x += direction.x
	enemy.velocity.y += direction.y

	if enemy.velocity.x > speed {
		enemy.velocity.x = speed
	} else if enemy.velocity.x < -speed {
		enemy.velocity.x = -speed
	}
	if enemy.velocity.y > speed {
		enemy.velocity.y = speed
	} else if enemy.velocity.y < -speed {
		enemy.velocity.y = -speed
	}

	startX := enemy.transform.x
	startY := enemy.transform.y

	// Apply velocity to position
	enemy.transform.x += enemy.velocity.x
	enemy.transform.y += enemy.velocity.y

	// Optional: friction (if you want smoothing, otherwise remove these lines)
	enemy.velocity.x *= 0.8
	enemy.velocity.y *= 0.8

	distance := math.Sqrt(
		math.Pow(enemy.transform.x-player.transform.x, 2) +
			math.Pow(enemy.transform.y-player.transform.y, 2),
	)

	checkCollisions(&enemy.transform, Vector2{
		x: startX,
		y: startY,
	})

	var attackDistance float64

	switch enemy.enemyType {
	case EnemyTypeEvren:
		attackDistance = 500
	case EnemyTypeEmran:
		attackDistance = 600
	case EnemyTypeNick:
		attackDistance = 200
	}

	if distance < attackDistance {
		enemy.gun.Shoot(&enemy.transform)
	}

}
func (enemy *Enemy) Draw(screen *ebiten.Image) {

	// var col color.RGBA

	// switch enemy.enemyType {
	// case EnemyTypeEvren:
	// 	col = color.RGBA{0, 255, 0, 255}
	// case EnemyTypeEmran:
	// 	col = color.RGBA{255, 0, 255, 255}
	// case EnemyTypeNick:
	// 	col = color.RGBA{0, 0, 255, 255}
	// }

	// drawRect(
	// 	screen,
	// 	enemy.transform,
	// 	col,
	// )

	var image *ebiten.Image

	switch enemy.enemyType {
	case EnemyTypeEvren:
		image = getCachedImage("enemies/evren")
	case EnemyTypeEmran:
		image = getCachedImage("enemies/emran")
	case EnemyTypeNick:
		image = getCachedImage("enemies/nick")
	}

	offset := Vector2{
		-110,
		500,
	}

	op := defaultImageOptions()
	op.Anchor = offset
	op.Scale = 4

	tr := enemy.GetTransform()
	tr.rotation += math.Pi / 2

	
	drawImageWithOptions(
		screen,
		image,
		tr,
		op,
	)

	// for _, point := range enemy.currentPath {
	// 	drawRect(screen, Transform{
	// 		x:        point.x - float64(pathFindingGridSize/2),
	// 		y:        point.y - float64(pathFindingGridSize/2),
	// 		width:    float64(pathFindingGridSize),
	// 		height:   float64(pathFindingGridSize),
	// 		rotation: 0,
	// 	}, color.RGBA{255, 0, 0, 50})
	// }
}

func (enemy *Enemy) GetTransform() Transform {
	return enemy.transform
}

func (enemy *Enemy) SetTransform(transform Transform) {
	enemy.transform = transform
}
