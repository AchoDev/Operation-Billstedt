package main

type Transform struct {
    x float64
    y float64
    z float64
    width float64
    height float64

    rotation float64
}

func (t *Transform) GetPosition() Vector2 {
    return Vector2{t.x, t.y}
}