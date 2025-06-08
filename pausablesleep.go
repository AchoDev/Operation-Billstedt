package main

import (
	"math/rand/v2"
	"sync"
	"time"
)

var isPaused bool
var pauseMutex sync.Mutex
var sleepMutex sync.Mutex
var currentSleeps []int

func Pause() {
    pauseMutex.Lock()
    isPaused = true
    pauseMutex.Unlock()
}

func UnPause() {
    pauseMutex.Lock()
    isPaused = false
    pauseMutex.Unlock()
}

func resetSleeps() {
    sleepMutex.Lock()
    defer sleepMutex.Unlock()
    currentSleeps = []int{}
}

func pausableSleep(duration time.Duration) {
    id := rand.IntN(10000000)
    currentSleeps = append(currentSleeps, id)

    defer func() {
        sleepMutex.Lock()
        defer sleepMutex.Unlock()
        for i, v := range currentSleeps {
            if v == id {
                currentSleeps = append(currentSleeps[:i], currentSleeps[i+1:]...)
                break
            }
        }
    }()

    start := time.Now()
    scaledDuration := time.Duration(float64(duration) / globalTimeScale)
    for time.Since(start) < scaledDuration {

        isInside := func () bool {
            sleepMutex.Lock()
            for _, v := range currentSleeps {
                if v == id {
                    sleepMutex.Unlock()
                    return true 
                }
            }

            sleepMutex.Unlock()

            return false
        }()

        if !isInside {
            return
        }

        pauseMutex.Lock()
        if isPaused {
            pauseMutex.Unlock()
            for isPaused || killscreen{
                time.Sleep(10 * time.Millisecond) // Check pause state periodically
                start = start.Add(10 * time.Millisecond)
            }
            continue
        }
        pauseMutex.Unlock()
        time.Sleep(10 * time.Millisecond) // Sleep in scaled increments
    }
}