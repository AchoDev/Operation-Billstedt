package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var (
    audioContext *audio.Context
    soundCache   = make(map[string][]byte) // Cache for decoded sound data
    playerCache = make(map[string]*audio.Player)
    cacheMutex   sync.Mutex                // Mutex to protect the cache
)

func init() {
    // Initialize the audio context with a sample rate of 44100 Hz
    audioContext = audio.NewContext(44100)
}

// LoadSound loads and caches an MP3 sound from the given file path
func LoadSound(filePath string) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()

    // Check if the sound is already cached
    if _, exists := soundCache[filePath]; exists {
        return
    }

    // Open the MP3 file
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Failed to open sound file: %v", err)
    }
    defer file.Close()

    // Decode the MP3 file
    decoded, err := mp3.Decode(audioContext, file)
    if err != nil {
        log.Fatalf("Failed to decode MP3 file: %v", err)
    }

    // Read the decoded data into memory
    data, err := io.ReadAll(decoded)
    if err != nil {
        log.Fatalf("Failed to read decoded MP3 data: %v", err)
    }

    // Store the decoded data in the cache
    soundCache[filePath] = data
}

// PlaySound plays a cached MP3 sound from the given file path
// PlaySound plays a cached MP3 sound from the given file path
func PlaySound(filePath string) {

    filePath = "assets/sounds/" + filePath + ".mp3"

    cacheMutex.Lock()
    data, exists := soundCache[filePath]
    cacheMutex.Unlock()

    if !exists {
        // Automatically load the sound if it's not cached
        LoadSound(filePath)

        // Retrieve the sound data again after loading
        cacheMutex.Lock()
        data = soundCache[filePath]
        cacheMutex.Unlock()
    }

    // Create a player from the cached data
    player, err := audioContext.NewPlayer(bytes.NewReader(data))
    if err != nil {
        log.Fatalf("Failed to create audio player: %v", err)
    }

    cacheMutex.Lock()
    playerCache[filePath] = player
    cacheMutex.Unlock()

    // Play the sound
    player.Play()
}

func StopSound(filePath string) {
    filePath = "assets/sounds/" + filePath + ".mp3"

    cacheMutex.Lock()
    player, exists := playerCache[filePath]
    cacheMutex.Unlock()

    if exists {
        player.Pause() // Stop the sound by pausing the player
        cacheMutex.Lock()
        delete(playerCache, filePath) // Remove the player from the cache
        cacheMutex.Unlock()
    }
}