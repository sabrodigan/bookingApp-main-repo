package main

import (
	"github.com/sabrodigan/bookings-app/internal/helpers"
	"testing"
	"time"
)

func TestRunApp(t *testing.T) {
	err := run()
	if err != nil {
		t.Errorf("Error starting app: %v", err)
	}

}
func TestSpinner(t *testing.T) {
	done := make(chan bool)
	go helpers.Spinner(100*time.Millisecond, done)

	// Simulate some work
	time.Sleep(1 * time.Second)

	// Stop the spinner
	done <- true
}
func TestTypeWriter(t *testing.T) {
	// Test the TypeWriter function
	phrase := "Testing the typewriter effect"
	delayMs := 100

	// Call the function
	helpers.TypeWriter(phrase, delayMs)

	//// Check if the output is as expected (this is a simple test, you can improve it)
	//if len(phrase) == 0 {
	//	t.Errorf("TypeWriter did not produce expected output")
	//}
}
