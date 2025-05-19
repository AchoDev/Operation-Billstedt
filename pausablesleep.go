package main

import (
	"sync"
	"time"
)

var isPaused bool
var pauseMutex sync.Mutex

func pausableSleep(duration time.Duration) {
    start := time.Now()
    for time.Since(start) < duration {
        pauseMutex.Lock()
        if isPaused {
            pauseMutex.Unlock()
            for isPaused {
                time.Sleep(10 * time.Millisecond) // Check pause state periodically
                start = start.Add(10 * time.Millisecond)
            }
            continue
        }
        pauseMutex.Unlock()
        time.Sleep(10 * time.Millisecond) // Sleep in small increments
    }
}