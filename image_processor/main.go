package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	log.Println("Initiating configuration...")
	configuration = new(Configuration)
	ex, _ := os.Executable()
	configuration.loadConfiguration(filepath.Dir(ex))
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	log.Println("Starting image processing routine...")
	log.Println("Processor running...")

	log.Println("Starting queue consumer...")
	go consume()
	log.Println("Queue consumer started...")
	wg.Wait()
}
