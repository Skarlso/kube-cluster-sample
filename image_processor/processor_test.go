package main

import "testing"

func TestProcessImages(t *testing.T) {
	imageQueue.imageQueue = append(imageQueue.imageQueue, 1, 2, 3, 4, 5)
	ch, con := processImages()
	con.Broadcast()
	var err error
	select {
	case r := <-ch:
		err = r.Error
	default:
		err = nil
	}
	if err != nil {
		t.Fatal(err)
	}
}
