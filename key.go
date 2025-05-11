package main

import "github.com/hajimehoshi/ebiten/v2"

// Map to track the previous state of keys
var previousKeyState = map[ebiten.Key]bool{}

// Map to track the previous state of mouse buttons
var previousMouseState = map[ebiten.MouseButton]bool{}

// Detect if a key was just pressed (button down event)
func isKeyJustPressed(key ebiten.Key) bool {
    // Check if the key is pressed in the current frame but was not pressed in the previous frame
    if ebiten.IsKeyPressed(key) && !previousKeyState[key] {
        return true
    }
    return false
}

// Detect if a mouse button was just pressed (button down event)
func isMouseButtonJustPressed(button ebiten.MouseButton) bool {
    // Check if the mouse button is pressed in the current frame but was not pressed in the previous frame
    if ebiten.IsMouseButtonPressed(button) && !previousMouseState[button] {
        return true
    }
    return false
}

// Update the key state at the end of each frame
func updateKeyState() {
    for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
        previousKeyState[key] = ebiten.IsKeyPressed(key)
    }
}

// Update the mouse state at the end of each frame
func updateMouseState() {
    for button := ebiten.MouseButtonLeft; button <= ebiten.MouseButtonRight; button++ {
        previousMouseState[button] = ebiten.IsMouseButtonPressed(button)
    }
}

// Update both key and mouse states
func updateInputState() {
    updateKeyState()
    updateMouseState()
}

