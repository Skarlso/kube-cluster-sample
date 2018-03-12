package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Skarlso/kube-cluster-sample/facerecog"
	"google.golang.org/grpc"
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
	db := DbConnection{}
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Println("Processing image id: ", i)
	c := facerecog.NewIdentityClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.IdentifyRequest{
		ImagePath: db.getPath(i),
	}
	if err != nil {
		log.Fatalf("could not send image: %v", err)
	}
	// TODO: Resposne Perons id will have the Id I have to update the database record with.
}
