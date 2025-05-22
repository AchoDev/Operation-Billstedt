package main

type GunStats struct {
    name          string
    cooldown      int
    offset       Vector2
    spread       float64
    hasCasing bool
    casingPoint Vector2
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
    hasCasing:     true,
    casingPoint: Vector2{77, 18},
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
    hasCasing: false,
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
    hasCasing: true,
    casingPoint: Vector2{30, 19},
    spread: 0.025,
    shootBehavior: RifleShoot,
}

var minigunStats = GunStats{
    name:     "Minigun",
    cooldown: 1,
    // cooldown: 5000,
    offset:   Vector2{
        x: 101,
        y: 20,
    },
    spread: 0.05,
    hasCasing: true,
    casingPoint: Vector2{0, 27},
    shootBehavior: MinigunShoot,
}
