package main

import (
	"encoding/binary"
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

// sendImage sends an image ID to the queue.
func sendImage(i int) error {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(configuration.ProducerAddress, config)
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(buffer, uint32(i))
	err := w.Publish("images", buffer)
	if err != nil {
		return fmt.Errorf("failed to publish image: %w", err)
	}

	w.Stop()
	return nil
}
