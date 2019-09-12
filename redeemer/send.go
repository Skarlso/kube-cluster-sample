package main

import (
	"encoding/binary"

	nsq "github.com/bitly/go-nsq"
)

// sendImage sends an image ID to the queue.
func sendImage(i int) error {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(configuration.ProducerAddress, config)
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(buffer, uint32(i))
	err := w.Publish("images", buffer)
	if err != nil {
		return err
	}

	w.Stop()
	return nil
}
