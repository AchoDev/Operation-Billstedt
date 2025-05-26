package main

type LivingThing interface {
    TakeDamage(damage int, direction float64)
}