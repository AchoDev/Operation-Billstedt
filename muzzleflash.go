package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type MuzzleFlash struct {
	transform Transform
    path     string
	target GameObject
	spawnTime time.Time
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

	flash.transform = newPos
	flash.transform.rotation = newPos.rotation + math.Pi/2
}

func NewMuzzleFlash(target GameObject, offset Vector2) *MuzzleFlash {

	randomNum := rand.IntN(4)

	flash := MuzzleFlash{
		path: "/sprites/muzzle-flash/" + fmt.Sprint(randomNum + 1),
		target: target,
		offset: offset,
		spawnTime: time.Now(),
	}

	flash.MoveToTarget()

	return &flash
}

func (flash *MuzzleFlash) Update() {
	flash.MoveToTarget()
	
	if time.Since(flash.spawnTime) > 40*time.Millisecond {
		removeGameObject(flash)
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