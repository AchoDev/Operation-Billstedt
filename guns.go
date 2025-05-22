package main

import (
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
    spread      float64
    hasCasing bool
    casingPoint Vector2
    shootBehavior func(transform *Transform, gun *GunBase)
}

func createMuzzleFlash(gun *GunBase) {    
    flash := NewMuzzleFlash(gun.carrier, gun.offset)

    gameObjects = append(gameObjects, flash)
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
        spread:       stats.spread,
        hasCasing: stats.hasCasing,
        casingPoint: stats.casingPoint,
        shootBehavior: stats.shootBehavior,
    }
}


// Example shoot behaviors
func PistolShoot(transform *Transform, gun *GunBase) {
    bullet := CreateBullet(transform, gun)
    gameObjects = append(gameObjects, bullet)

    pushBack(gun.carrier, 2.0)
    if !gun.isEnemy {
        camera.Shake(transform.rotation, 2.0)
    }
}

func ShotgunShoot(transform *Transform, gun *GunBase) {
    for i := -2; i <= 2; i++ {
        angleOffset := float64(i) * 0.1
        tr := Transform{
            x:      transform.x,
            y:      transform.y,
            rotation: transform.rotation + angleOffset, 
        }
        bullet := CreateBullet(&tr, gun)
        gameObjects = append(gameObjects, bullet)
    }
    
    pushBack(gun.carrier, 20.0)

    if !gun.isEnemy {
        camera.Shake(transform.rotation, 5.0)
    }

}

func RifleShoot(transform *Transform, gun *GunBase) {
    go func() {
        for i := 0; i < 5; i++ {
            bullet := CreateBullet(transform, gun)
            gameObjects = append(gameObjects, bullet)

            if !gun.isEnemy {
                camera.Shake(transform.rotation, 2.5)
            }

            pushBack(gun.carrier, 2.0)

            pausableSleep(100 * time.Millisecond)
            createMuzzleFlash(gun)
        }
    }()
}

func MinigunShoot(transform *Transform, gun *GunBase) {
    go func() {
        for i := 0; i < 20; i++ {
            bullet := CreateBullet(transform, gun)
            gameObjects = append(gameObjects, bullet)

            pushBack(gun.carrier, 3.0)

            pausableSleep(50 * time.Millisecond)

            if !gun.isEnemy {
                camera.Shake(transform.rotation, 5.0)
            }

            createMuzzleFlash(gun)
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

func pushBack(target GameObject, amount float64) {
    if player, ok := target.(*Player); ok {
        player.velocity.x -= amount * math.Cos(player.transform.rotation)
        player.velocity.y -= amount * math.Sin(player.transform.rotation)
    }
    if enemy, ok := target.(*Enemy); ok {
        enemy.velocity.x -= amount * math.Cos(enemy.transform.rotation)
        enemy.velocity.y -= amount * math.Sin(enemy.transform.rotation)
    }
}

