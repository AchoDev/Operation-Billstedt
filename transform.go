package main

import (
	"math"
)

type Transform struct {
    x float64
    y float64
    z float64
    width float64
    height float64

    rotation float64
}

func (t *Transform) GetPosition() Vector2 {
    return Vector2{t.x, t.y}
}

func (t *Transform) RotateAround(angle float64, pivot Vector2) {

    offset := Vector2{
        x: t.x - pivot.x,
        y: t.y - pivot.y,
    }

    // Rotate the offset
    rotatedX := offset.x*math.Cos(angle) - offset.y*math.Sin(angle)
    rotatedY := offset.x*math.Sin(angle) + offset.y*math.Cos(angle)

    // Update the position by adding the pivot back
    t.x = pivot.x + rotatedX
    t.y = pivot.y + rotatedY
}