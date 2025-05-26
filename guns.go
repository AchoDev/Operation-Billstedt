package main

import (
	"math"
	"time"
)

func createMuzzleFlash(gun *GunStats) {    
    flash := NewMuzzleFlash(gun.carrier, gun.offset)

    addGameObject(flash)
}

func (g *GunStats) Shoot(transform *Transform) {
    if g.cooldownTimer != -1 {
        return
    }
    if g.shootBehavior != nil {
        g.shootBehavior(transform, g)
        createMuzzleFlash(g)
    }
    StartGunCooldown(g)
}

func (g *GunStats) GetCooldown() int {
    return g.cooldown
}

func (g *GunStats) GetCooldownTimer() float64 {
    return g.cooldownTimer
}

func (g *GunStats) SetCooldownTimer(timer float64) {
    g.cooldownTimer = timer
}

func (g *GunStats) IsEnemy() bool {
    return g.isEnemy
}

func (g *GunStats) SetIsEnemy(isEnemy bool) {
    g.isEnemy = isEnemy
}

func (g *GunStats) Name() string {
    return g.name
}

func NewGun(stats GunStats, carrier GameObject) *GunStats {
    fromEnemy := false

    if _, ok := carrier.(*Enemy); ok {
        fromEnemy = true
    }
    
    stats.cooldownTimer = -1
    stats.isEnemy = fromEnemy
    stats.carrier = carrier

    return &stats
}


// Example shoot behaviors
func PistolShoot(transform *Transform, gun *GunStats) {
    bullet := CreateBullet(transform, gun)
    addGameObject(bullet)

    pushBack(gun.carrier, 2.0)
    PlaySound("pistol")
    if !gun.isEnemy {
        camera.Shake(transform.rotation, 5.0)
    }
}

func ShotgunShoot(transform *Transform, gun *GunStats) {
    for i := -2; i <= 2; i++ {
        angleOffset := float64(i) * 0.1
        tr := Transform{
            x:      transform.x,
            y:      transform.y,
            rotation: transform.rotation + angleOffset, 
        }
        bullet := CreateBullet(&tr, gun)
        addGameObject(bullet)
    }
    
    PlaySound("shotgun")
    pushBack(gun.carrier, 20.0)

    if !gun.isEnemy {
        camera.Shake(transform.rotation, 10.0)
        // camera.MotionBlur(400)
    }

}

func RifleShoot(transform *Transform, gun *GunStats) {
    go func() {
        for i := 0; i < 5; i++ {
            bullet := CreateBullet(transform, gun)
            addGameObject(bullet)

            if !gun.isEnemy {
                camera.Shake(transform.rotation, 5.0)
            }

            pushBack(gun.carrier, 2.0)
            PlaySound("rifle")

            pausableSleep(100 * time.Millisecond)
            createMuzzleFlash(gun)
        }
    }()
}

func MinigunShoot(transform *Transform, gun *GunStats) {
    go func() {
        PlaySound("minigun")
        for i := 0; i < 20; i++ {
            bullet := CreateBullet(transform, gun)
            addGameObject(bullet)

            pushBack(gun.carrier, 3.0)

            pausableSleep(50 * time.Millisecond)

            if !gun.isEnemy {
                camera.Shake(transform.rotation, 7.5)
            }

            createMuzzleFlash(gun)
        }
        StopSound("minigun")
    }()
}

func StartGunCooldown(gun *GunStats) {
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

