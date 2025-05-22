package main

type GunStats struct {
	name          string
	cooldown      int
	offset        Vector2
	shootBehavior func(transform *Transform, gun *GunBase)
}

var pistolStats = GunStats{
	name:     "Pistol",
	cooldown: 100,
	offset: Vector2{
		x: 97,
		y: 7,
	},
	shootBehavior: PistolShoot,
}

var shotgunStats = GunStats{
	name:     "Shotgun",
	cooldown: 3000,
	offset: Vector2{
		x: 78,
		y: 19,
	},
	shootBehavior: ShotgunShoot,
}

var rifleStats = GunStats{
	name:     "Rifle",
	cooldown: 5000,
	offset: Vector2{
		x: 94,
		y: 17,
	},
	shootBehavior: RifleShoot,
}

var minigunStats = GunStats{
	name:     "Minigun",
	cooldown: 5000,
	offset: Vector2{
		x: 101,
		y: 20,
	},
	shootBehavior: MinigunShoot,
}
