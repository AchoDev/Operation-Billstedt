package main

import "github.com/hajimehoshi/ebiten/v2"

func getMousePosition() Vector2 {
	x, y := ebiten.CursorPosition()

	x += int(camera.x)
	y += int(camera.y)

	x -= int(camera.width / 2)
	y -= int(camera.height / 2)

	return Vector2{
		x: float64(x),
		y: float64(y),
	}
}
