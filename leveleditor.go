package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"os"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var pos = Vector2{300, 1000}
var selectedSprite int
var currentScale Vector2 = Vector2{1, 1}
var levelEditorActivated bool = false
var selectedTool int = 0
var currentRotation float64 = 0
var currentZEditor int = 0
var onionSkin bool = false
var hideGameobjects bool = false

func DrawLevelEditor(screen *ebiten.Image, level Level) {
	if !levelEditorActivated {
		return
	}

	drawAbsoluteRect(screen, Transform{
		x:        pos.x,
		y:        pos.y,
		width:    5000,
		height:   150,
		rotation: 0,
	}, color.RGBA{0, 0, 0, 200})

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Z layer: %.0f", float64(currentZEditor)), 10, 950)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rotation: %.2f", float64(currentRotation)), 10, 970)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Scale: %.2f", currentScale), 10, 990)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Onion: %t", onionSkin), 10, 1010)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hide Gameobjects: %t", hideGameobjects), 10, 1030)

	counter := 0
	sprites := level.GetSprites()
	keys := make([]string, 0, len(sprites))
	for k := range sprites {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// y := 1080.0
	for _, name := range keys {
		sprite := sprites[name]
		transform := Transform{
			x:        float64(counter*120) + pos.x,
			y:        pos.y,
			width:    100,
			height:   100,
			rotation: 0,
		}

		if counter == selectedSprite {
			size := 10.0
			drawAbsoluteRect(screen, Transform{
				x:        transform.x,
				y:        transform.y,
				width:    transform.width + size*2,
				height:   transform.height + size*2,
				rotation: 0,
			}, color.RGBA{255, 129, 129, 255})
		}

		drawAbsoluteImage(screen, sprite, transform)

		counter++
	}

	gridPos := getMouseGridPosition()
	sprite := sprites[keys[selectedSprite]]
	op := defaultImageOptions()
	op.Alpha = 100

	if selectedTool == 1 {
		sprite = ebiten.NewImage(100, 100)
		sprite.Fill(color.RGBA{255, 100, 200, 255})
	}
	drawImageWithOptions(screen, sprite, Transform{
		x:        float64(gridPos.x * 100),
		y:        float64(gridPos.y * 100),
		width:    100 * currentScale.x,
		height:   100 * currentScale.y,
		rotation: currentRotation,
	}, op)
}

func UpdateLevelEditor(level Level) {

	if isKeyJustPressed(ebiten.KeyP) {
		levelEditorActivated = !levelEditorActivated
		pauseMutex.Lock()
		isPaused = levelEditorActivated
		pauseMutex.Unlock()
	}

	if !levelEditorActivated {
		return
	}

	// Detect mouse scroll to change selected sprite
	_, yoff := ebiten.Wheel()
	if yoff != 0 {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			var direction float64
			if yoff > 0 {
				direction = 1
			} else {
				direction = -1
			}

			currentScale.x += direction * 0.5
			currentScale.y += direction * 0.5

			currentScale.x = math.Round(currentScale.x*2) / 2
			currentScale.y = math.Round(currentScale.y*2) / 2
		} else {
			currentScale.x += yoff * 0.5
			currentScale.y += yoff * 0.5
		}
			
		currentScale.x = math.Max(currentScale.x, 0.1)
		currentScale.y = math.Max(currentScale.y, 0.1)
	}

	if isKeyJustPressed(ebiten.Key1) {
		selectedTool = 0
	}
	if isKeyJustPressed(ebiten.Key2) {
		selectedTool = 1
	}

	if isKeyJustPressed(ebiten.KeyTab) {

		direction := 1

		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			direction = -1
		}

		selectedSprite += direction
		if selectedSprite >= len(level.GetSprites()) {
			selectedSprite = 0
		}
	}

	if isKeyJustPressed(ebiten.KeyM) {
		onionSkin = !onionSkin
	}

	if isKeyJustPressed(ebiten.KeyN) {
		hideGameobjects = !hideGameobjects
	}

	if isKeyJustPressed(ebiten.KeyR) {
		currentRotation += math.Pi / 2
		if currentRotation >= 2*math.Pi {
			currentRotation = 0
		}
	}

	if isKeyJustPressed(ebiten.KeyUp) {
		currentScale.x += 0.25
		currentScale.y += 0.25
		currentScale.x = math.Round(currentScale.x*4) / 4
		currentScale.y = math.Round(currentScale.y*4) / 4

		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			currentScale.x -= 0.25
		}
	}

	if isKeyJustPressed(ebiten.KeyDown) {
		currentScale.x -= 0.25
		currentScale.y -= 0.25
		currentScale.x = math.Round(currentScale.x*4) / 4
		currentScale.y = math.Round(currentScale.y*4) / 4

		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			currentScale.x += 0.25
		}
	}

	if isMouseButtonJustPressed(ebiten.MouseButtonRight) {
		mousePosition := getMousePosition()

		fmt.Println("Searching for tile to delete at:", mousePosition)

		tiles := level.GetTiles()

		// Reverse the tiles list

		if selectedTool == 1 {
			tiles = level.GetColliders()
		}

		for i := len(tiles) - 1; i >= 0; i-- {
			tile := tiles[i]
			
			if tile.Z != float64(currentZEditor) {
				continue
			}


			tileTr := Transform{
				x:      tile.X * 100,
				y:      tile.Y * 100,
				width:  tile.Width * 100,
				height: tile.Height * 100,
			}
			halfWidth := tileTr.width / 2
			halfHeight := tileTr.height / 2
			if mousePosition.x >= tileTr.x-halfWidth && mousePosition.x <= tileTr.x+halfWidth &&
				mousePosition.y >= tileTr.y-halfHeight && mousePosition.y <= tileTr.y+halfHeight {

				tiles = append(tiles[:i], tiles[i+1:]...)

				if selectedTool == 1 {
					level.SetColliders(tiles)
					break
				} else {
					level.SetTiles(tiles)
				}
				break
			}
		}
	}

	if isKeyJustPressed(ebiten.KeyLeft) {
		currentZEditor -= 1
		if currentZEditor <= 0 {
			currentZEditor = 0
		}
	}

	if isKeyJustPressed(ebiten.KeyRight) {
		currentZEditor += 1
	}

	if isKeyJustPressed(ebiten.KeyO) {
		data := map[string]interface{}{
			"tiles":     level.GetTiles(),
			"colliders": level.GetColliders(),
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling tiles to JSON:", err)
			return
		}

		err = os.WriteFile("level-output.json", jsonData, 0644)

		if err != nil {
			fmt.Println("Error writing JSON to file:", err)
			return
		} else {
			fmt.Println("Level saved to level-output.json")
		}
	}

	speed := 10.0

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		speed = 30
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		camera.x -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		camera.x += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		camera.y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		camera.y += speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		camera.zoom -= 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		camera.zoom += 0.01
	}

	if isMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		sprites := level.GetSprites()
		keys := make([]string, 0, len(sprites))
		for k := range sprites {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		counter := 0
		for _, name := range keys {
			if counter == selectedSprite {

				gridPosition := getMouseGridPosition()

				if selectedTool == 1 {

					fmt.Println("Creating collider at", gridPosition)
					colliders := level.GetColliders()
					for _, collider := range colliders {
						if collider.X == gridPosition.x && collider.Y == gridPosition.y {
							fmt.Println("Collider already exists at", gridPosition)
							return
						}
					}
					colliders = append(colliders, Tile{
						X:      gridPosition.x,
						Y:      gridPosition.y,
						Z: 	float64(currentZEditor),
						Width:  currentScale.x,
						Height: currentScale.y,
					})
					level.SetColliders(colliders)
					fmt.Println("Added collider at", gridPosition)

					break
				}

				tiles := level.GetTiles()

				// for i, tile := range tiles {
				// 	if tile.X == gridPosition.x && tile.Y == gridPosition.y {
				// 		tiles = append(tiles[:i], tiles[i+1:]...)
				// 		level.SetTiles(tiles)
				// 		fmt.Println("Replaced tile at", gridPosition)
				// 	}
				// }

				tiles = append(tiles, Tile{
					X:      gridPosition.x,
					Y:      gridPosition.y,
					Z: 	float64(currentZEditor),
					Width:  currentScale.x,
					Height: currentScale.y,
					Rotation: currentRotation,
					Sprite: name,
				})

				level.SetTiles(tiles)

				fmt.Println("Added tile at", gridPosition)
				break
			}

			counter++
		}
	}
}

func getMouseGridPosition() Vector2 {

	gridStep := 1.0

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		gridStep = 0.5
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) {
		gridStep = 0.1
	}

	cursorX, cursorY := ebiten.CursorPosition()
	worldX := float64(cursorX) / camera.zoom
	worldY := float64(cursorY) / camera.zoom

	worldX += camera.x
	worldY += camera.y

	worldX -= camera.width / 2 / camera.zoom
	worldY -= camera.height / 2 / camera.zoom

	worldX /= 100
	worldY /= 100

	gridPosition := Vector2{
		x: math.Round(worldX/gridStep) * gridStep,
		y: math.Round(worldY/gridStep) * gridStep,
	}

	return gridPosition
}

func getOrderedSprites(level Level) []string {
	sprites := level.GetSprites()
	keys := make([]string, 0, len(sprites))
	for k := range sprites {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
