package main

import (
	"time"
)

type Gun interface {
    Name() string
    Shoot(transform *Transform)
    GetCooldown() int
    GetCooldownTimer() float64
    SetCooldownTimer(timer float64)
    IsEnemy() bool
    SetIsEnemy(isEnemy bool)
}

type Pistol struct {
    cooldownTimer float64
    isEnemy bool
}

func createGun(gun Gun, isEnemy bool) Gun {
    gun.SetCooldownTimer(-1)
    gun.SetIsEnemy(isEnemy)
    return gun
}

func (g *Pistol) Shoot(transform *Transform) {

    if g.cooldownTimer != -1 {
        return
    }

    // Create a bullet
    bullet := Bullet{
        transform: Transform{
            x: transform.x,
            y: transform.y,
            width: 25,
            height: 10,
        },
        angle: transform.rotation, // Set the angle based on the player's rotation
        speed: 10,
        fromEnemy: g.isEnemy,
    }

    StartGunCooldown(g)

    // Update the bullet's position
    gameObjects = append(gameObjects, &bullet)
}

func (g *Pistol) GetCooldown() int {
    return 500 // Cooldown in milliseconds
}

func (g *Pistol) GetCooldownTimer() float64 {
    return g.cooldownTimer
}

func (g *Pistol) SetCooldownTimer(timer float64) {
    g.cooldownTimer = timer
}

func (g *Pistol) IsEnemy() bool {
    return g.isEnemy
}

func (g *Pistol) SetIsEnemy(isEnemy bool) {
    g.isEnemy = isEnemy
}

func (g *Pistol) Name() string {
    return "Pistol"
}


type Shotgun struct {
    cooldownTimer float64
    isEnemy bool
}

func (g *Shotgun) Shoot(transform *Transform) {

    if g.cooldownTimer != -1 {
        return
    }

    // Create multiple bullets for the shotgun
    for i := -2; i <= 2; i++ {
        angleOffset := float64(i) * 0.2 // Adjust the spread of the shotgun
        bullet := Bullet{
            transform: Transform{
                x: transform.x,
                y: transform.y,
                width: 25,
                height: 10,
            },
            angle: transform.rotation + angleOffset, // Set the angle based on the player's rotation
            speed: 10,
            fromEnemy: g.isEnemy,
        }

        // Update the bullet's position
        gameObjects = append(gameObjects, &bullet)
    }

    StartGunCooldown(g)
}

func (g *Shotgun) GetCooldown() int {
    return 5000 // Cooldown in milliseconds
}

func (g *Shotgun) GetCooldownTimer() float64 {
    return g.cooldownTimer
}

func (g *Shotgun) SetCooldownTimer(timer float64) {
    g.cooldownTimer = timer
}

func (g *Shotgun) IsEnemy() bool {
    return g.isEnemy
}

func (g *Shotgun) SetIsEnemy(isEnemy bool) {
    g.isEnemy = isEnemy
}

func (g *Shotgun) Name() string {
    return "Shotgun"
}


type Rifle struct {
    cooldownTimer float64
    isEnemy bool
}

func (g *Rifle) Shoot(transform *Transform) {

    if g.cooldownTimer != -1 {
        return
    }

    go func () {
        for i := 0; i < 5; i++ {
            // Create a bullet
            bullet := Bullet{
                transform: Transform{
                    x: transform.x,
                    y: transform.y,
                    width: 25,
                    height: 10,
                },
                angle: transform.rotation, // Set the angle based on the player's rotation
                speed: 10,
                fromEnemy: g.isEnemy,
            }
            gameObjects = append(gameObjects, &bullet)
            
            time.Sleep(100 * time.Millisecond) // Sleep for 100 milliseconds
        }
    }()

    StartGunCooldown(g)
}

func (g *Rifle) GetCooldown() int {
    return 1000 // Cooldown in milliseconds
}

func (g *Rifle) GetCooldownTimer() float64 {
    return g.cooldownTimer
}

func (g *Rifle) SetCooldownTimer(timer float64) {
    g.cooldownTimer = timer
}

func (g *Rifle) IsEnemy() bool {
    return g.isEnemy
}

func (g *Rifle) SetIsEnemy(isEnemy bool) {
    g.isEnemy = isEnemy
}

func (g *Rifle) Name() string {
    return "Rifle"
}

func StartGunCooldown(gun Gun) {
    go func() {
        gun.SetCooldownTimer(float64(gun.GetCooldown()))
        start := time.Now()
        
        for gun.GetCooldownTimer() > 0 {
            time.Sleep(1 * time.Millisecond)
            gun.SetCooldownTimer(float64(gun.GetCooldown()) - float64(time.Since(start).Milliseconds()))
        }

        gun.SetCooldownTimer(-1)
    }()
}
