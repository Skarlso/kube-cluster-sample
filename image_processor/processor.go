package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
	"google.golang.org/grpc"
)

var imageQueue = make([]int, 0)
var c = sync.NewCond(&sync.Mutex{})
var circuitBreaker *CircuitBreaker

func processImages() {
	for {
		c.L.Lock()
		for len(imageQueue) == 0 {
			c.Wait()
		}
		circuitBreaker = NewCircuitBreaker()
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
	c := facerecog.NewIdentifyClient(conn)
	h := facerecog.NewHealthCheckClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	path, err := db.getPath(i)
	if err != nil {
		log.Printf("error while getting path for image id: %d", i)
		return
	}
	circuitBreaker.F = func() (*facerecog.IdentifyResponse, error) {
		r, err := c.Identify(ctx, &facerecog.IdentifyRequest{
			ImagePath: path,
		})
		return r, err
	}
	circuitBreaker.Ping = func() bool {
		_, err := h.HealthCheck(ctx, &facerecog.Empty{})
		if err != nil {
			return false
		}
		return true
	}
	r, err := circuitBreaker.Call()
	if err != nil {
		dbErr := db.updateImageWithFailedStatus(i)
		if dbErr != nil {
			log.Printf("could not update image to failed status: %v", dbErr)
		}
		return
	}
	p, err := db.getPersonFromImage(r.GetImageName())
	if err != nil {
		log.Printf("warning: could not retrieve person: %v", err)
		return
	}
	log.Println("got person: ", p.Name)
	log.Println("updating record with person id")
	err = db.updateImageWithPerson(p.ID, i)
	if err != nil {
		log.Printf("warning: could not update image record: %v", err)
	}
	log.Println("done")
}
