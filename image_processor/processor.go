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
func processImages() {
	for {
		c.L.Lock()
		for len(imageQueue.imageQueue) == 0 {
			c.Wait()
		}
		circuitBreaker = NewCircuitBreaker()
		for len(imageQueue.imageQueue) > 0 {
			processImage(imageQueue.drain())
		}
		c.L.Unlock()
	}
}

func processImage(i int) error {
	db := DbConnection{}
	conn, err := grpc.Dial(configuration.GRPCAddress, grpc.WithInsecure())
	if err != nil {
		message := fmt.Sprintf("could not connect to grpc on: %s", configuration.GRPCAddress)
		return errors.New(message)
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
		return errors.New(message)
	}
	circuitBreaker.F = func() (*facerecog.IdentifyResponse, error) {
		r, err := c.Identify(ctx, &facerecog.IdentifyRequest{
			ImagePath: path,
		})
		return r, err
	}
	circuitBreaker.Ping = func() bool {
		_, err := h.HealthCheck(ctx, &facerecog.Empty{})
		return err != nil
	}
	r, err := circuitBreaker.Call()
	if err != nil {
		dbErr := db.updateImageWithFailedStatus(i)
		if dbErr != nil {
			message := fmt.Sprintf("could not update image to failed status: %v. Reason for failed status: %v", dbErr, err)
			return errors.New(message)
		}
		return err
	}
	p, err := db.getPersonFromImage(r.GetImageName())
	if err != nil {
		message := fmt.Sprintf("warning: could not retrieve person: %v", err)
		return errors.New(message)
	}
	log.Println("got person: ", p.Name)
	log.Println("updating record with person id")
	err = db.updateImageWithPerson(p.ID, i)
	if err != nil {
		message := fmt.Sprintf("warning: could not update image record: %v", err)
		return errors.New(message)
	}
	log.Println("done")
	return nil
}
