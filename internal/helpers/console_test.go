package helpers

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// TestTypeWriter tests the TypeWriter function
func TestTypeWriter(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout
	
	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Call TypeWriter with a very small delay to speed up the test
	go TypeWriter("Hello", 1)
	
	// Allow time for TypeWriter to complete
	time.Sleep(20 * time.Millisecond)
	
	// Close the writer to flush the buffer
	w.Close()
	
	// Restore original stdout
	os.Stdout = oldStdout
	
	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	// Check the output
	if buf.String() != "Hello" {
		t.Errorf("TypeWriter output incorrect, got: %s, want: %s", buf.String(), "Hello")
	}
}

// TestSpinner tests the Spinner function
func TestSpinner(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout
	
	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Create a done channel
	done := make(chan bool)
	
	// Start the spinner with a very small delay
	go Spinner(1*time.Millisecond, done)
	
	// Allow time for at least one cycle
	time.Sleep(10 * time.Millisecond)
	
	// Signal the spinner to stop
	done <- true
	
	// Allow time for the spinner to process the done signal
	time.Sleep(5 * time.Millisecond)
	
	// Close the writer to flush the buffer
	w.Close()
	
	// Restore original stdout
	os.Stdout = oldStdout
	
	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	// Check that some output was produced
	if buf.Len() == 0 {
		t.Error("Spinner did not produce any output")
	}
	
	// Check that the output contains at least one of the spinner characters
	spinnerChars := []string{"|", "/", "-", "\\"}
	found := false
	for _, char := range spinnerChars {
		if strings.Contains(buf.String(), char) {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Spinner output does not contain any spinner characters: %s", buf.String())
	}
}