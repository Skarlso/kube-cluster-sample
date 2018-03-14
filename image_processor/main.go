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

func init() {
	log.Println("Initiating environment...")
	initiateEnvironment()
	configuration = new(Configuration)
	configuration.loadConfiguration()
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
