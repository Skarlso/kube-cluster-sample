package main

import "log"

// Image we'll let the DB assign an ID to an image
type Image struct {
	ID       int
	PersonID int
	Path     []byte
}

func (i *Image) saveImage() error {
	log.Println("Saving image path")
	dc := new(DbConnection)
	return dc.saveImage(i)
}

func (i *Image) loadImage(ID int) {
	log.Println("Loading image with ID: ", ID)
	dc := new(DbConnection)
	i, err := dc.loadImage(ID)
	if err != nil {
		log.Fatal(err)
	}
}
