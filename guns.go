package main

import "time"

type Gun interface {
    Name() string
    Shoot(transform *Transform)
    GetCooldown() int
    IsCoolingDown() bool
}

type Pistol struct {
    coolingDown bool
}

func (g *Pistol) Shoot(transform *Transform) {
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
    }

    // Update the bullet's position
    gameObjects = append(gameObjects, &bullet)
}

func (g *Pistol) GetCooldown() int {
    return 500 // Cooldown in milliseconds
}

func (g *Pistol) Name() string {
    return "Pistol"
}


type Shotgun struct {}

func (g *Shotgun) Shoot(transform *Transform) {
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
        }

        // Update the bullet's position
        gameObjects = append(gameObjects, &bullet)
    }


}

func (g *Shotgun) GetCooldown() int {
    return 1000 // Cooldown in milliseconds
}


func (g *Shotgun) Name() string {
    return "Shotgun"
}


type Rifle struct {}

func (g *Rifle) Shoot(transform *Transform) {
    go func () {
        for i := 0; i < 10; i++ {
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
            }
            gameObjects = append(gameObjects, &bullet)
            
            time.Sleep(100 * time.Millisecond) // Sleep for 100 milliseconds
        }
    }()
}

func (g *Rifle) GetCooldown() int {
    return 1000 // Cooldown in milliseconds
}

func (g *Rifle) Name() string {
    return "Rifle"
}