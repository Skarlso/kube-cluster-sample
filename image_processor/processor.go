package main

import (
	"log"
	"sync"
)

var imageQueue = make([]int, 0)
var c = sync.NewCond(&sync.Mutex{})

func processImages() {
	for {
		c.L.Lock()
		for len(imageQueue) == 0 {
			c.Wait()
		}
		for len(imageQueue) > 0 {
			processImage(imageQueue[0])
			imageQueue = imageQueue[1:]
		}
		c.L.Unlock()
	}
}

func processImage(i int) {
	log.Println("Processing image id: ", i)
}
