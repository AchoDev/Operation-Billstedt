package main

import "math"

type Camera struct {
	x      float64
	y      float64
    offset Vector2
    velocity Vector2
	width  float64
	height float64
	zoom   float64
}

var camera Camera = Camera{
	x:      0,
	y:      0,
    offset: Vector2{
        x: 0,
        y: 0,
    },
	width:  1920,
	height: 1080,
	zoom:   1.2,
}

func (c *Camera) Shake(angle float64, intensity float64) {

    direction := Vector2{
        x: math.Cos(angle + math.Pi),
        y: math.Sin(angle + math.Pi),
    }

    camera.velocity.x += direction.x * intensity
    camera.velocity.y += direction.y * intensity
}

func (c *Camera) Update() {

    c.velocity.x -= c.offset.x * 0.1
    c.velocity.y -= c.offset.y * 0.1

    c.offset.x += c.velocity.x
    c.offset.y += c.velocity.y

    c.velocity.x *= 0.9
    c.velocity.y *= 0.9


    target := Vector2{
        x: 750,
        y: player.transform.y,
    }

    direction := Vector2{
        x: math.Cos(player.transform.rotation),
        y: math.Sin(player.transform.rotation),
    }

    target.x += direction.x * 100
    target.y += direction.y * 100

    diff := Vector2{
		x: target.x - camera.x,
		y: target.y - camera.y,
	}

	zoomDiff := camera.zoom - 1.2

	camera.zoom -= zoomDiff * 0.1
	camera.x += diff.x * 0.1
	camera.y += diff.y * 0.1
}