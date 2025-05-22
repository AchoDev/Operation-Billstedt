package main

type GunStats struct {
    name          string
    cooldown      int
    offset       Vector2
    spread       float64
    shootBehavior func(transform *Transform, gun *GunBase)
}

var pistolStats = GunStats{
    name:     "Pistol",
    cooldown: 100,
    offset:   Vector2{
        x: 97,
        y: 7,
    },
    spread: 0.01,
    shootBehavior: PistolShoot,
}

var shotgunStats = GunStats{
    name:     "Shotgun",
    cooldown: 1,
    // cooldown: 1,
    offset:   Vector2{
        x: 78,
        y: 19,
    },
    spread: 0.05,
    shootBehavior: ShotgunShoot,
}

var rifleStats = GunStats{
    name:     "Rifle",
    cooldown: 1,
    // cooldown: 5000,
    offset:   Vector2{
        x: 94,
        y: 17,
    },
    spread: 0.025,
    shootBehavior: RifleShoot,
}

var minigunStats = GunStats{
    name:     "Minigun",
    cooldown: 5000,
    offset:   Vector2{
        x: 101,
        y: 20,
    },
    spread: 0.05,
    shootBehavior: MinigunShoot,
}
