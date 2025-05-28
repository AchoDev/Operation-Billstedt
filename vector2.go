package main

import (
	"math"
)

type Vector2 struct {
    x float64
    y float64
}

func (vec *Vector2) normalize() {
    x := float64(vec.x)
    y := float64(vec.y)
    len := math.Sqrt(x * x + y * y)

    if vec.x != 0 {
        vec.x /= len
    }
    if vec.y != 0 {
        vec.y /= len
    }   
}

func (vec *Vector2) Multiply(number float64) {
    vec.x *= number
    vec.y *= number
}

func (vec *Vector2) Multiply2(x, y float64) {
    vec.x *= x
    vec.y *= y
}

func (vec *Vector2) Set(number float64) {
    vec.x = number
    vec.y = number
}

func (vec *Vector2) Set2(x, y float64) {
    vec.x = x
    vec.y = y
}

func (vec *Vector2) SetVector(other Vector2) {
    vec.x = other.x
    vec.y = other.y
}