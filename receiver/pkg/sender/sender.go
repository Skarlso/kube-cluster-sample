package sender

import (
	"encoding/binary"

	"github.com/nsqio/go-nsq"
)

// Config is the necessary configuration for the nsq sender.
type Config struct {
	Address string
}

type nsqSender struct {
	config Config
}

// NewNSQSender defines a sender which uses NSQ as a queuing system.
func NewNSQSender(cfg Config) *nsqSender {
	return &nsqSender{config: cfg}
}

// SendImage takes the ID of an image and sends it along in the queue.
// I could agrue here to send the path as is,
// but the ID is the key, so sending that makes sense.
// Also, I want to be able to change anything about
// the image later on, but the ID will still remain
// to be the ID.
func (s *nsqSender) SendImage(i uint64) error {
	config := nsq.NewConfig()
	// The procuder needs to be co-located with nsqd so it can send messages to a local queue.
	// The consumers use lookupd to find a queue.
	// This means the image receiver needs to be on the same pod as the queue.
	// Or I have to split it up.
	w, _ := nsq.NewProducer(s.config.Address, config) // TODO: WRONG ADDRESS. This should be the producer
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, i)
	err := w.Publish("images", buffer)
	if err != nil {
		return err
	}

	w.Stop()
	return nil
}
