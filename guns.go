package main

import (
	"fmt"
	"math"
	"time"
)
type GunBase struct {
    cooldownTimer float64
    isEnemy       bool
    carrier     GameObject
    cooldown      int
    name          string
    offset       Vector2
    shootBehavior func(transform *Transform, gun *GunBase)
}

func createMuzzleFlash(gun *GunBase) {    
    flash := NewMuzzleFlash(gun.carrier, gun.offset)

    fmt.Println(gun)

    gameObjects = append(gameObjects, flash)

    go func() {
        pausableSleep(50 * time.Millisecond)
        removeGameObject(flash)
    }()
}

func (g *GunBase) Shoot(transform *Transform) {
    if g.cooldownTimer != -1 {
        return
    }
    if g.shootBehavior != nil {
        g.shootBehavior(transform, g)
        createMuzzleFlash(g)
    }
    StartGunCooldown(g)
}

func (g *GunBase) GetCooldown() int {
    return g.cooldown
}

func (g *GunBase) GetCooldownTimer() float64 {
    return g.cooldownTimer
}

func (g *GunBase) SetCooldownTimer(timer float64) {
    g.cooldownTimer = timer
}

func (g *GunBase) IsEnemy() bool {
    return g.isEnemy
}

func (g *GunBase) SetIsEnemy(isEnemy bool) {
    g.isEnemy = isEnemy
}

func (g *GunBase) Name() string {
    return g.name
}

func NewGun(stats GunStats, carrier GameObject) *GunBase {
    fromEnemy := false

    getCachedImage("/sprites/muzzle-flash")

    if _, ok := carrier.(*Enemy); ok {
        fromEnemy = true
    }
    
    return &GunBase{
        cooldownTimer: -1,
        isEnemy:       fromEnemy,
        carrier:      carrier,
        offset:         stats.offset,
        cooldown:      stats.cooldown,
        name:          stats.name,
        shootBehavior: stats.shootBehavior,
    }
}

// Example shoot behaviors
func PistolShoot(transform *Transform, gun *GunBase) {
    bullet := Bullet{
        transform: Transform{
            x:      transform.x,
            y:      transform.y,
            width:  25,
            height: 10,
        },
        angle:     transform.rotation,
        speed:     10,
        fromEnemy: gun.isEnemy,
    }
    gameObjects = append(gameObjects, &bullet)
}

func ShotgunShoot(transform *Transform, gun *GunBase) {
    for i := -2; i <= 2; i++ {
        angleOffset := float64(i) * 0.2
        bullet := Bullet{
            transform: Transform{
                x:      transform.x,
                y:      transform.y,
                width:  25,
                height: 10,
            },
            angle:     transform.rotation + angleOffset,
            speed:     10,
            fromEnemy: gun.isEnemy,
        }
        gameObjects = append(gameObjects, &bullet)
    }


    pushbackForce := 20.0
    push := Vector2{
        x: pushbackForce * math.Cos(transform.rotation),
        y: pushbackForce * math.Sin(transform.rotation),
    }
    
    if player, ok := gun.carrier.(*Player); ok {
        player.velocity.x -= push.x
        player.velocity.y -= push.y
    }
    if enemy, ok := gun.carrier.(*Enemy); ok {
        enemy.velocity.x -= push.x
        enemy.velocity.y -= push.y
    }

    if !gun.isEnemy {
        camera.Shake(transform.rotation, 5.0)
    }

}

func RifleShoot(transform *Transform, gun *GunBase) {
    go func() {
        for i := 0; i < 5; i++ {
            bullet := Bullet{
                transform: Transform{
                    x:      transform.x,
                    y:      transform.y,
                    width:  25,
                    height: 10,
                },
                angle:     transform.rotation,
                speed:     10,
                fromEnemy: gun.isEnemy,
            }
            gameObjects = append(gameObjects, &bullet)
            pausableSleep(100 * time.Millisecond)
        }
    }()
}

func MinigunShoot(transform *Transform, gun *GunBase) {
    go func() {
        for i := 0; i < 20; i++ {
            bullet := Bullet{
                transform: Transform{
                    x:      transform.x,
                    y:      transform.y,
                    width:  25,
                    height: 10,
                },
                angle:     transform.rotation,
                speed:     15,
                fromEnemy: gun.isEnemy,
            }
            gameObjects = append(gameObjects, &bullet)
            pausableSleep(50 * time.Millisecond)
        }
    }()
}

func StartGunCooldown(gun *GunBase) {
    go func() {
        gun.SetCooldownTimer(float64(gun.GetCooldown()))
        start := time.Now()
        
        for gun.GetCooldownTimer() > 0 {
            pausableSleep(1 * time.Millisecond)
            gun.SetCooldownTimer(float64(gun.GetCooldown()) - float64(time.Since(start).Milliseconds()))
        }

        gun.SetCooldownTimer(-1)
    }()
}

