package main

import (
	"fmt"
	"math"
	"time"
)

func createMuzzleFlash(gun *GunStats) {    
    flash := NewMuzzleFlash(gun.carrier, gun.offset)

    addGameObject(flash)
}

func (g *GunStats) Shoot(transform *Transform) {

    fmt.Println(g.locked)

    if g.cooldownTimer != -1 || g.currentAmmo <= 0 || g.locked {
        return
    }

    if g.firintCooldown != -1 {
        return
    }

    if g.shootBehavior != nil {
        g.shootBehavior(transform, g)
        createMuzzleFlash(g)
        g.currentAmmo -= 1
        WaitForNextShot(g)
    }
    
    if g.currentAmmo <= 0 {
        StartGunCooldown(g)
    }
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
    stats.firintCooldown = -1
    stats.isEnemy = fromEnemy
    stats.carrier = carrier
    stats.currentAmmo = stats.maxAmmo

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
        bullet := CreateBullet(transform, gun)
        addGameObject(bullet)

        if !gun.isEnemy {
            camera.Shake(transform.rotation, 5.0)
        }

        pushBack(gun.carrier, 2.0)
        PlaySound("rifle")


        createMuzzleFlash(gun)

    }()
}
    
func MinigunShoot(transform *Transform, gun *GunStats) {
    PlaySound("minigun")
    
    bullet := CreateBullet(transform, gun)
    addGameObject(bullet)
    
    pushBack(gun.carrier, 3.0)

    if !gun.isEnemy {
        camera.Shake(transform.rotation, 7.5)
    }

    createMuzzleFlash(gun)
        
    StopSound("minigun")
}

func StartGunCooldown(gun *GunStats) {
    go func() {
        gun.cooldownTimer = float64(gun.cooldown)
        start := time.Now()
        
        for gun.cooldownTimer > 0 {
            pausableSleep(1 * time.Millisecond)
            gun.cooldownTimer = float64(gun.cooldown) - float64(time.Since(start).Milliseconds()) / 1000.0
        }
        gun.cooldownTimer = -1
        gun.currentAmmo = gun.maxAmmo
    }()
}

func WaitForNextShot(gun *GunStats) {
    go func() {

        if gun.firingRate != 0 {
            gun.firintCooldown = 1 / float64(gun.firingRate)
    
            for gun.firintCooldown > 0 {
                pausableSleep(1 * time.Millisecond)
                gun.firintCooldown -= (1 / 1000.0) 
                fmt.Println("Firing Cooldown:", gun.firintCooldown)
            }
        }


        gun.firintCooldown = -1
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

