package main

import (
	"errors"
	"fmt"
	"log"
)

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

// Image we'll let the DB assign an ID to an image
type Image struct {
	ID       int
	PersonID int
	Path     []byte
	Status   Status
}

func (i *Image) saveImage() error {
	log.Println("Saving image path")
	if ok, _ := i.searchPath(); ok {
		e := fmt.Sprintf("image with path '%s' already exists", string(i.Path))
		return errors.New(e)
	}
	dc := new(DbConnection)
	return dc.saveImage(i)
}

func (i *Image) loadImage(ID int) {
	log.Println("Loading image with ID: ", ID)
	dc := new(DbConnection)
	var err error
	*i, err = dc.loadImage(ID)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Image) searchPath() (bool, error) {
	dc := new(DbConnection)
	return dc.searchPath(string(i.Path))
}
