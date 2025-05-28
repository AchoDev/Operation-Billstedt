package main

import "github.com/hajimehoshi/ebiten/v2"

func getArrowDirection() Vector2 {
    // This function returns the direction of the arrow based on the current input state.
    if ebiten.IsKeyPressed(ebiten.KeyUp) {
        return Vector2{0, -1}
    } else if ebiten.IsKeyPressed(ebiten.KeyDown) {
        return Vector2{0, 1}
    } else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        return Vector2{-1, 0}
    } else if ebiten.IsKeyPressed(ebiten.KeyRight) {
        return Vector2{1, 0}
    }
    return Vector2{0, 0}
}