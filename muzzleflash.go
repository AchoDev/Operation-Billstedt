package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type MuzzleFlash struct {
	transform Transform
    path     string
	target GameObject
	offset Vector2
}

func (flash *MuzzleFlash) MoveToTarget() {
	newPos := flash.target.GetTransform()

	// rotate around the target based on the offset

	angle := newPos.rotation
	sin, cos := math.Sin(angle), math.Cos(angle)

	rotatedX := flash.offset.x*cos - flash.offset.y*sin
	rotatedY := flash.offset.x*sin + flash.offset.y*cos

	newPos.x += rotatedX
	newPos.y += rotatedY


	fmt.Println(rotatedX, rotatedY, flash.offset)

	flash.transform = newPos
	flash.transform.rotation = newPos.rotation
}

func NewMuzzleFlash(target GameObject, offset Vector2) *MuzzleFlash {
	flash := MuzzleFlash{
		path: "/sprites/muzzle-flash",
		target: target,
		offset: offset,
	}

	flash.MoveToTarget()

	return &flash
}

func (flash *MuzzleFlash) Update() {
	flash.MoveToTarget()


	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		flash.offset.y += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		flash.offset.y -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		flash.offset.x += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		flash.offset.x -= 1
	}
}

func (flash *MuzzleFlash) Draw(screen *ebiten.Image) {
	image := getCachedImage(flash.path)
	
    drawImage(screen, image, flash.transform)
}

func (flash *MuzzleFlash) GetTransform() Transform {
	return flash.transform
}

func (flash *MuzzleFlash) SetTransform(transform Transform) {
	flash.transform = transform
}