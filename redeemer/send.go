package main

import (
	"encoding/binary"

	nsq "github.com/bitly/go-nsq"
)

// NSQ handles sending topics to the queue
type NSQ struct {
}

// Image we'll let the DB assign an ID to an image
type Image struct {
	ID       int
	PersonID int
	Path     []byte
	Status   Status
}

// sendImage sends an image ID to the queue.
func (n *NSQ) sendImage(i Image) error {
	config := nsq.NewConfig()
	// The procuder needs to be co-located with nsqd so it can send messages to a local queue.
	// The consumers use lookupd to find a queue.
	// This means the image receiver needs to be on the same pod as the queue.
	w, _ := nsq.NewProducer(configuration.ProducerAddress, config)
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(buffer, uint32(i.ID))
	err := w.Publish("images", buffer)
	if err != nil {
		return err
	}

	w.Stop()
	return nil
}
