package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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

type Circle struct {
    Center Vector2
    Radius float64
}

// Check collision between a circle and a rotated rectangle
func CircleRotatedRectColliding(circle Circle, rect Rect) bool {
    corners := getCorners(rect)
    closestPoint := closestPointOnRect(circle.Center, corners)

    dx := closestPoint.x - circle.Center.x
    dy := closestPoint.y - circle.Center.y
    distanceSquared := dx*dx + dy*dy

    return distanceSquared <= circle.Radius*circle.Radius
}

func closestPointOnRect(point Vector2, corners []Vector2) Vector2 {
    closest := corners[0]
    minDistanceSquared := distanceSquared(point, closest)

    for _, corner := range corners[1:] {
        distSquared := distanceSquared(point, corner)
        if distSquared < minDistanceSquared {
            closest = corner
            minDistanceSquared = distSquared
        }
    }

    return closest
}

func distanceSquared(v1, v2 Vector2) float64 {
    dx := v1.x - v2.x
    dy := v1.y - v2.y
    return dx*dx + dy*dy
}




type Collider struct {
    transform Transform
}

func (collider *Collider) Update() {

}

func (collider *Collider) Draw(screen *ebiten.Image) {
    rect := ebiten.NewImage(int(collider.transform.width), int(collider.transform.height))
    rect.Fill(color.RGBA{0, 100, 0, 255})
    opts := &ebiten.DrawImageOptions{}
    opts.GeoM.Translate(collider.transform.x, collider.transform.y)
    screen.DrawImage(rect, opts)
}

func (collider *Collider) GetTransform() Transform {
    return collider.transform
}