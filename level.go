package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
    x, y int
    width, height int
    sprite string
    collide bool
}

type Level interface {
    GetTiles() []Tile
    GetSprites() map[string]*ebiten.Image
    UpdateLevel()
}

type Level1 struct {
    tiles []Tile
    sprites map[string]*ebiten.Image
}

var levelEditorActivated bool = false

func DrawLevel(screen *ebiten.Image, level Level) {
    tiles := level.GetTiles()
    sprites := level.GetSprites()

    for _, tile := range tiles {
        tileImage := sprites[tile.sprite]
        if tileImage == nil {
            continue
        }

        scale := Vector2{
            x: float64(tile.width * 20) / float64(tileImage.Bounds().Dx()),
            y: float64(tile.height * 20) / float64(tileImage.Bounds().Dy()),
        }

        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(float64(tile.x), float64(tile.y))
        op.GeoM.Scale(scale.x, scale.y)
    }


    
}

func (level *Level1) GetTiles() []Tile {
    return level.tiles
}

func (level *Level1) GetSprites() map[string]*ebiten.Image {
    return level.sprites
}

func (level *Level1) UpdateLevel() {



}


