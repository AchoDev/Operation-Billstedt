package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Train struct {
	transform Transform
	parts     map[string]*ebiten.Image
	length    int
	currentAlpha float64
}

func CreateTrain(transform Transform) *Train {
	return &Train{
		transform: transform,
		parts: map[string]*ebiten.Image{
			"front":  loadImage("assets/sprites/u2-front.png"),
			"middle": loadImage("assets/sprites/u2-middle.png"),
			"door-light": loadImage("assets/sprites/train-door-light.png"),
		},
		length: 3,
	}
}

func (train *Train) Update() {

}

func (train *Train) calculateEnemySpawnPoints() []Vector2 {
	result := []Vector2{}

	distance := 690.0
	startDistance := 50.0

	y := train.transform.y - (distance + startDistance)

	for i := 0; i < train.length; i++ {
		result = append(result, Vector2{
			x: train.transform.x + 135,
			y: y + 200,
		})
		result = append(result, Vector2{
			x: train.transform.x + 135,
			y: y - 200,
		})

		y -= distance
	}

	result = append(result, Vector2{
		x: train.transform.x + 135,
		y: train.transform.y + 150,
	})
	result = append(result, Vector2{
		x: train.transform.x + 135,
		y: train.transform.y - 250,
	})

	return result
}

func (train *Train) Draw(screen *ebiten.Image) {

	op := defaultImageOptions()
	op.OriginalImageSize = true
	op.Scale = 0.2

	middle := train.parts["middle"]
	front := train.parts["front"]
	
	distance := 690.0
	startDistance := 50.0
	
	y := train.transform.y - (distance + startDistance)

	lightOp := defaultImageOptions()
	lightOp.OriginalImageSize = true
	lightOp.Scale = 0.175
	lightOp.Alpha = train.currentAlpha * 255
	
	for i := 0; i < train.length; i++ {
		
		drawRect(screen, Transform{
			x:      train.transform.x + 2.5,
			y:      y + 350,
			width: 60,
			height: 40,
			}, color.Black)	

		y -= distance
	}
	drawImageWithOptions(screen, front, train.transform, op)

	drawImageWithOptions(screen, train.parts["door-light"], Transform{
		x:      train.transform.x + 135,
		y:      train.transform.y + 150,
		width:  1,
		height: 1,
		rotation: math.Pi / 2,
	}, lightOp)

	drawImageWithOptions(screen, train.parts["door-light"], Transform{
		x:      train.transform.x + 135,
		y:      train.transform.y - 250,
		width:  1,
		height: 1,
		rotation: math.Pi / 2,
	}, lightOp)

	y = train.transform.y - (distance + startDistance)

	for i := 0; i < train.length; i++ {	
		drawImageWithOptions(screen, middle, Transform{
			x:      train.transform.x,
			y:      y,
			width:  train.transform.width,
			height: train.transform.height,
		}, op)


		drawImageWithOptions(screen, train.parts["door-light"], Transform{
			x:      train.transform.x + 135,
			y:      y + 200,
			width:  1,
			height: 1,
			rotation: math.Pi / 2,
		}, lightOp)

		drawImageWithOptions(screen, train.parts["door-light"], Transform{
			x:      train.transform.x + 135,
			y:      y - 200,
			width:  1,
			height: 1,
			rotation: math.Pi / 2,
		}, lightOp)

		y -= distance
	}
}

func (train *Train) GetTransform() Transform {
	return train.transform
}

func (train *Train) Drive(distance float64, speed float64) {

	times := 3

	go func() {
		startY := train.transform.y
		targetY := startY + distance
		for {
			remaining := targetY - train.transform.y

			// If close enough, snap to target and stop
			if math.Abs(remaining) < 1 {
				train.transform.y = targetY
				break
			}

			// Move towards the target with smoothing
			train.transform.y += remaining * speed * 0.15

			time.Sleep(time.Second / 60)
		}

		// time.Sleep(time.Second)
	
		cycles := 2.0
		for i := 0.0; i < (cycles * 1.5 + 0.75) * 60; i++ {
			
			train.currentAlpha = (-math.Cos(float64(i) * (math.Pi / 180) * 4) + 1) / 2 

			time.Sleep(time.Second / 60)
		}

		for _, spawnPoint  := range train.calculateEnemySpawnPoints() {
			// get random enemy

			random := rand.Intn(3)
			var enemyType EnemyType

			switch random {
			case 0:
				enemyType = EnemyTypeEvren
			case 1:
				enemyType = EnemyTypeEmran
			case 2:
				enemyType = EnemyTypeNick
			}

			enemy := createEnemy(
				int(spawnPoint.x),
				int(spawnPoint.y),
				enemyType,
			)
			gameObjects = append(gameObjects, enemy)
		}

		for i := 0; i < .75 * 60; i++ {
			train.currentAlpha = (math.Cos(float64(i) * (math.Pi / 180) * 4) + 1) / 2 
			time.Sleep(time.Second / 60)
		}

		// drive away again

		velocity := 0.0
		midY := train.transform.y

		for {
			velocity += 0.1
			train.transform.y += velocity
			if train.transform.y > midY + 5000 {
				break
			}

			time.Sleep(time.Second / 60)
		}

		times -= 1

		if times > 0 {
			time.Sleep(time.Second * time.Duration(rand.Intn(15)))
			train.transform.y = startY
			train.Drive(distance, speed)
		}
	}()
}
