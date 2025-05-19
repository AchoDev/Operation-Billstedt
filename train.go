package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Train struct {
	transform Transform
	parts     map[string]*ebiten.Image
	length    int
}

func (train *Train) Update() {
}

func (train *Train) Draw(screen *ebiten.Image) {

	op := defaultImageOptions()
	op.OriginalImageSize = true
	op.Scale = 0.2

	middle := train.parts["middle"]
	front := train.parts["front"]
	drawImageWithOptions(screen, front, train.transform, op)

	y := train.transform.y - 150

	for i := 0; i < train.length; i++ {
		if i != 0 {
			drawRect(screen, Transform{
				x:      train.transform.x + 2.5,
				y:      y + 350,
				width: 60,
				height: 40,
			}, color.Black)
		}

		y -= 680
	}

	y = train.transform.y - 150

	for i := 0; i < train.length; i++ {	
		drawImageWithOptions(screen, middle, Transform{
			x:      train.transform.x,
			y:      y,
			width:  train.transform.width,
			height: train.transform.height,
		}, op)



		y -= 680
	}
}

func (train *Train) GetTransform() Transform {
	return train.transform
}

func (train *Train) Drive(distance float64, speed float64) {
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
	}()
}
