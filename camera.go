package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

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

const (
    maxBlurFrames = 5
    scaleFactor   = 0.5
)


var (
    blurFrames []*ebiten.Image
    blurIndex  int
    w, h       int
    blurMotion bool
    blurAlpha  float64
)

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

    target.x += direction.x * 150
    target.y += direction.y * 150

    diff := Vector2{
		x: target.x - camera.x,
		y: target.y - camera.y,
	}

	zoomDiff := camera.zoom - 1.2

	camera.zoom -= zoomDiff * 0.1
	camera.x += diff.x * 0.1
	camera.y += diff.y * 0.1
}

func ApplyMotionBlur(screen *ebiten.Image) {
    if !blurMotion {
        return
    }

    for i := 0; i < maxBlurFrames; i++ {
        idx := (blurIndex + i) % maxBlurFrames
        weight := float32(i+1) / float32(maxBlurFrames+1)
        alpha := weight * float32(blurAlpha)

        op := &ebiten.DrawImageOptions{}
        op.GeoM.Scale(1.0/scaleFactor, 1.0/scaleFactor)
        op.ColorScale.ScaleAlpha(alpha)
        screen.DrawImage(blurFrames[idx], op)
    }

    curr := blurFrames[blurIndex]
    curr.Clear()
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(scaleFactor, scaleFactor)
    curr.DrawImage(screen, op)

    blurIndex = (blurIndex + 1) % maxBlurFrames


}
func InitMotionBlur(screenW, screenH int) {
    w = int(float64(screenW))
    h = int(float64(screenH))
    blurFrames = make([]*ebiten.Image, maxBlurFrames)
    for i := range blurFrames {
        blurFrames[i] = ebiten.NewImage(w, h)
    }
}

var lastMotionBlur time.Time

func (*Camera) MotionBlur(duration float64) {
    lastMotionBlur = time.Now()
    start := lastMotionBlur
    blurAlpha = 0.2
    blurMotion = true

    go func() {
        pausableSleep(time.Duration(duration) * time.Millisecond)
        if lastMotionBlur.Equal(start) {
            for blurAlpha > 0 {
                blurAlpha -= 0.01 // Decrease alpha gradually
                pausableSleep(16 * time.Millisecond) // Approx. 60 FPS
            }
            
            blurMotion = false
        }
    }()
}
