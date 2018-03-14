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
var sendTryCount = 0
var timeOut = time.Minute
var maxTryCount = 5
var t time.Time

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
	c := facerecog.NewIdentifyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if t.Add(timeOut).Before(time.Now()) {
		log.Printf("timeout over. allowing sending of images again.")
		sendTryCount = 0
	}
	if sendTryCount >= maxTryCount {
		log.Printf("max sending try count of %d reached. sending not allowed for %v time period.", maxTryCount, timeOut)
		return
	}

	path, _ := db.getPath(i)
	r, err := c.Identify(ctx, &facerecog.IdentifyRequest{
		ImagePath: path,
	})

	if err != nil {
		log.Printf("could not send image: %v", err)
		sendTryCount++
		if sendTryCount >= maxTryCount {
			log.Printf("maximum try of %d sends reached. disabling for %v time period.", maxTryCount, timeOut)
			t = time.Now()
		}
		return
	}
	p, err := db.getPersonFromImage(r.GetImageName())
	if err != nil {
		log.Fatalf("could not retrieve person: %v", err)
	}
	log.Println("got person: ", p.Name)
	log.Println("updating record with person id")
	err = db.updateImageWithPerson(p.ID, i)
	if err != nil {
		log.Fatalf("could not update image record: %v", err)
	}
	log.Println("done")
}
