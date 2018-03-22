package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
	"google.golang.org/grpc"
)

// ImageQueue is a thread safe slice of image IDs.
type ImageQueue struct {
	imageQueue []int
	sync.RWMutex
}

func (i *ImageQueue) add(image int) {
	i.Lock()
	defer i.Unlock()
	i.imageQueue = append(i.imageQueue, image)
}

func (i *ImageQueue) drain() int {
	i.Lock()
	defer i.Unlock()
	ret := i.imageQueue[0]
	i.imageQueue = i.imageQueue[1:]
	return ret
}

var imageQueue = &ImageQueue{imageQueue: make([]int, 0)}
var c = sync.NewCond(&sync.Mutex{})
var circuitBreaker *CircuitBreaker

// Return a result channel
func processImages(response chan Response) {
	for {
		c.L.Lock()
		for len(imageQueue.imageQueue) == 0 {
			c.Wait()
		}
		circuitBreaker = NewCircuitBreaker()
		for len(imageQueue.imageQueue) > 0 {
			processImage(imageQueue.drain(), response)
		}
		c.L.Unlock()
	}
}

func processImage(i int, response chan Response) {
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
		message := fmt.Sprintf("error while getting path for image id: %d", i)
		resp := Response{Error: errors.New(message)}
		response <- resp
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
			message := fmt.Sprintf("could not update image to failed status: %v", dbErr)
			resp := Response{Error: errors.New(message)}
			response <- resp
		}
		resp := Response{Error: err}
		response <- resp
		return
	}
	p, err := db.getPersonFromImage(r.GetImageName())
	if err != nil {
		message := fmt.Sprintf("warning: could not retrieve person: %v", err)
		resp := Response{Error: errors.New(message)}
		response <- resp
		return
	}
	log.Println("got person: ", p.Name)
	log.Println("updating record with person id")
	err = db.updateImageWithPerson(p.ID, i)
	if err != nil {
		message := fmt.Sprintf("warning: could not update image record: %v", err)
		resp := Response{Error: errors.New(message)}
		response <- resp
	}
	log.Println("done")
}
