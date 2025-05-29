package main

type shootingType int

const (
    HalfAutomatic shootingType = iota
    FullAutomatic
)

type GunStats struct {
    name          string
    cooldown      int
    offset       Vector2
    spread       float64
    hasCasing bool
    damage      int
    casingPoint Vector2
    cooldownTimer float64
    isEnemy       bool
    carrier     GameObject
    maxAmmo   int
    currentAmmo int

    locked bool
    shootingType shootingType

    firingRate int
    firintCooldown float64

    shootBehavior func(transform *Transform, gun *GunStats)
}

var pistolStats = GunStats{
    name:     "Pistol",
    cooldown: 2,
    offset:   Vector2{
        x: 97,
        y: 7,
    },
    spread: 0.01,
    damage: 35,
    hasCasing:     true,
    maxAmmo: 12,

    shootingType: HalfAutomatic,

    casingPoint: Vector2{77, 18},
    shootBehavior: PistolShoot,
}

var shotgunStats = GunStats{
    name:     "Shotgun",
    cooldown: 5,
    // cooldown: 1,
    offset:   Vector2{
        x: 78,
        y: 19,
    },
    spread: 0.05,
    damage: 50,
    hasCasing: false,

    shootingType: HalfAutomatic,

    maxAmmo: 2,
    shootBehavior: ShotgunShoot,
}

var rifleStats = GunStats{
    name:     "Rifle",
    cooldown: 10,
    // cooldown: 5000,
    offset:   Vector2{
        x: 94,
        y: 17,
    },
    hasCasing: true,
    damage: 20,
    casingPoint: Vector2{30, 19},
    spread: 0.025,
    maxAmmo: 30,

    shootingType: FullAutomatic,

    firingRate: 100,
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
    damage: 15,
    maxAmmo: 200, 

    shootingType: FullAutomatic,

    firingRate: 500,
    casingPoint: Vector2{0, 27},
    shootBehavior: MinigunShoot,
}
