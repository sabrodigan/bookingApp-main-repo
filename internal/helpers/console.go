// /home/sabrodigan/go/src/bookingApp/internal/helpers/console.go
// Package helpers provides utility functions for the booking application
package helpers

import (
	"bytes"
	_ "embed" // Use a blank import to satisfy the compiler for the //go:embed directive.
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/sabrodigan/bookings-app/internal/config"
)

var app *config.AppConfig

//go:embed typewriter-key.mp3
var typewriterSoundData []byte

var (
	otoCtx      *oto.Context
	soundPcm    []byte
	isAudioInit bool
)

// NewHelpers sets the config for the helpers package
func NewHelpers(a *config.AppConfig) {
	app = a
}

// InitAudio initializes the audio system for the typewriter effect.
// It decodes an embedded MP3 file into raw PCM data to be played.
// It should be called once at application startup.
func InitAudio() {
	if len(typewriterSoundData) == 0 {
		log.Println("INFO: Typewriter sound file not found or empty. Sound will be disabled.")
		isAudioInit = false
		return
	}

	// Decode the embedded MP3 file.
	decodedMp3, err := mp3.NewDecoder(bytes.NewReader(typewriterSoundData))
	if err != nil {
		log.Printf("ERROR: decoding MP3: %v. Typewriter sound will be disabled.", err)
		isAudioInit = false
		return
	}

	// Prepare the oto context for playback.
	var ready <-chan struct{}
	// Use '=' to assign to the package-level 'otoCtx' and re-assign 'err'.
	// This avoids shadowing the global variable.
	otoCtx, ready, err = oto.NewContext(decodedMp3.SampleRate(), 2, 2)
	if err != nil {
		log.Printf("ERROR: creating oto context: %v. Typewriter sound will be disabled.", err)
		isAudioInit = false
		return
	}
	// Wait for the audio context to be ready.
	<-ready

	// Read the entire decoded stream into memory. This is efficient
	// because we can replay the sound without re-decoding the MP3 every time.
	soundPcm, err = io.ReadAll(decodedMp3)
	if err != nil {
		log.Printf("ERROR: reading decoded MP3: %v. Typewriter sound will be disabled.", err)
		isAudioInit = false
		return
	}

	isAudioInit = true
	log.Println("INFO: Audio system initialized for typewriter effect.")
}

// playKeySound plays the loaded typewriter key press sound once.
// It's run in a goroutine to prevent it from blocking the text display.
func playKeySound() {
	if !isAudioInit || otoCtx == nil {
		return
	}
	// Create a new Player. This is cheap and safe to do concurrently.
	player := otoCtx.NewPlayer(bytes.NewReader(soundPcm))
	defer player.Close()

	player.Play()

	// Wait for the sound to finish before the function returns and the player is closed.
	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}
}

// TypeWriter simulates a typewriter effect with a custom delay and sound.
func TypeWriter(phrase string, delayMs int) {
	for _, char := range phrase {
		fmt.Print(string(char))
		// Play a sound for printable characters, but not for spaces or newlines.
		if isAudioInit && char != ' ' && char != '\n' && char != '\t' {
			go playKeySound()
		}
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}

// Spinner displays a spinning animation until the done channel receives a value
func Spinner(delay time.Duration, done chan bool) {
	// Spinner characters
	for {
		for _, r := range `|/-\` {
			select {
			case <-done:
				fmt.Print("\r") // Clear the spinner
				return
			default:
				fmt.Printf("\r%c", r) // Print spinner character
				time.Sleep(delay)
			}
		}
	}
}

// ClientError logs a client error and sends a standard response.
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError logs a server error with a stack trace and sends a standard response.
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Server errors should go to the error log
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
