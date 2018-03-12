package main

import (
	"log"
	"sync"
)

func init() {
	log.Println("Initiating environment...")
	initiateEnvironment()
	configuration = new(Configuration)
	configuration.loadConfiguration()
	log.Println(configuration)
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
