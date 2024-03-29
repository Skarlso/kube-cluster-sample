package main

import (
	"log"
)

// Status is an Image status representation
type Status int

const (
	// PENDING -- not yet send to face recognition service
	_ Status = iota
	// PROCESSED -- processed by face recognition service; even if no person was found for the image
	PROCESSED
)

func init() {
	log.Println("Initiating environment...")
	initiateEnvironment()
	configuration = new(Configuration)
	configuration.loadConfiguration()
}

func main() {
	log.Println("Starting redeeming failed images...")
	db := DbConnection{}
	if err := db.open(); err != nil {
		log.Fatal("unable to make db connection: ", err)
	}
	ids, err := db.getAllNonProcessedImages()
	if err != nil {
		log.Fatal("unable to get all failed images: ", err)
	}

	log.Printf("found %d failed images to redeem\n", len(ids))
	log.Println("starting parallel processing of failed images...")
	for _, id := range ids {
		err := sendImage(id)
		if err != nil {
			log.Println("failed to send image with id: ", id)
		}
	}
}
