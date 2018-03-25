package main

import (
	"log"
	"sync"
)

// Person is a person
type Person struct {
	ID   int
	Name string
}

// Status is an Image status representation
type Status int

const (
	// PENDING -- not yet send to face recognition service
	PENDING Status = iota
	// PROCESSED -- processed by face recognition service; even if no person was found for the image
	PROCESSED
	// FAILEDPROCESSING -- for whatever reason the processing failed and this image is flagged for a retry
	FAILEDPROCESSING
)

func init() {
	log.Println("Initiating environment...")
	initiateEnvironment()
	configuration = new(Configuration)
	configuration.loadConfiguration()
}

// Response potential error response from the image processor routine.
// The caller decides what to do in case of an error rather than the routine.
// Error handling is the responsibility of the one who creates the channel.
type Response struct {
	Error error
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	log.Println("Starting image processing routine...")
	go processImages()
	log.Println("Starting queue consumer...")
	go consume()
	wg.Wait()
}
