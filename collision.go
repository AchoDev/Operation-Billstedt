package main

import (
	"math"
)


type Rect struct {
    Center Vector2
    Width  float64
    Height float64
    Angle  float64 // Rotation angle in radians
}

// Separating Axis Theorem (SAT) for rotated rectangles
func RotatedRectsColliding(rect1, rect2 Rect) bool {
    axes := []Vector2{
        rotate(Vector2{x: 1, y: 0}, rect1.Angle),
        rotate(Vector2{x: 0, y: 1}, rect1.Angle),
        rotate(Vector2{x: 1, y: 0}, rect2.Angle),
        rotate(Vector2{x: 0, y: 1}, rect2.Angle),
    }

    for _, axis := range axes {
        if !overlapOnAxis(rect1, rect2, axis) {
            return false
        }
    }

    return true
}

func rotate(v Vector2, angle float64) Vector2 {
    cos := math.Cos(angle)
    sin := math.Sin(angle)
    return Vector2{
        x: v.x*cos - v.y*sin,
        y: v.x*sin + v.y*cos,
    }
}

func getCorners(rect Rect) []Vector2 {
    halfWidth := rect.Width / 2
    halfHeight := rect.Height / 2

    corners := []Vector2{
        {-halfWidth, -halfHeight},
        {halfWidth, -halfHeight},
        {halfWidth, halfHeight},
        {-halfWidth, halfHeight},
    }

    for i := range corners {
        corners[i] = rotate(corners[i], rect.Angle)
        corners[i].x += rect.Center.x
        corners[i].y += rect.Center.y
    }

    return corners
}

func project(rect Rect, axis Vector2) (float64, float64) {
    corners := getCorners(rect)
    min := dot(corners[0], axis)
    max := min

    for _, corner := range corners[1:] {
        proj := dot(corner, axis)
        if proj < min {
            min = proj
        }
        if proj > max {
            max = proj
        }
    }

    return min, max
}

func overlapOnAxis(rect1, rect2 Rect, axis Vector2) bool {
    min1, max1 := project(rect1, axis)
    min2, max2 := project(rect2, axis)

    return !(max1 < min2 || max2 < min1)
}

func dot(v1, v2 Vector2) float64 {
    return v1.x*v2.x + v1.y*v2.y
}