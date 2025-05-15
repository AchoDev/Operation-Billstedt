package main

import "github.com/hajimehoshi/ebiten/v2"

func getMousePosition() Vector2 {
	x, y := ebiten.CursorPosition()


	worldX := float64(x)
	worldY := float64(y)

	worldX /= camera.zoom
	worldY /= camera.zoom

	worldX += camera.x
	worldY += camera.y

	worldX -= camera.width / 2 / camera.zoom
	worldY -= camera.height / 2 / camera.zoom

	return Vector2{
		x: worldX,
		y: worldY,
	}
}
