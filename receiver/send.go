package main

import (
	"encoding/binary"
	"net"

	nsq "github.com/bitly/go-nsq"
)

// NSQ handles sending topics to the queue
type NSQ struct {
}

// sendImage sends an image ID to the queue.
// I could agrue here to send the path as is,
// but the ID is the key, so sending that makes sense.
// Also, I want to be able to change anything about
// the image later on, but the ID will still remain
// to be the ID.
func (n *NSQ) sendImage(i Image) error {
	config := nsq.NewConfig()
	laddr := "127.0.0.1"

	config.LocalAddr, _ = net.ResolveTCPAddr("tcp", laddr+":0")

	w, _ := nsq.NewProducer("127.0.0.1:4150", config)
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(buffer, uint32(i.ID))
	err := w.Publish("images", buffer)
	if err != nil {
		return err
	}

	w.Stop()
	return nil
}
