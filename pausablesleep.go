package main

import (
	"sync"
	"time"
)

var isPaused bool
var pauseMutex sync.Mutex

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

func pausableSleep(duration time.Duration) {
    start := time.Now()
    scaledDuration := time.Duration(float64(duration) / globalTimeScale)
    for time.Since(start) < scaledDuration {
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