package sender

import (
	"encoding/binary"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
)

// Config is the necessary configuration for the nsq sender.
type Config struct {
	Address string
	Logger  zerolog.Logger
}

type nsqSender struct {
	config Config
}

// NewNSQSender defines a sender which uses NSQ as a queuing system.
func NewNSQSender(cfg Config) *nsqSender {
	return &nsqSender{config: cfg}
}

// SendImage takes the ID of an image and sends it along in the queue.
// I could argue here to send the path as is,
// but the ID is the key, so sending that makes sense.
// Also, I want to be able to change anything about
// the image later on, but the ID will still remain
// to be the ID.
func (s *nsqSender) SendImage(i uint64) error {
	s.config.Logger.Info().Uint64("id", i).Msg("sending image id to queue...")
	config := nsq.NewConfig()
	// The producer needs to be co-located with nsqd, so it can send messages to a local queue.
	// The consumers use lookupd to find a queue.
	// This means the image receiver needs to be on the same pod as the queue.
	// Or I have to split it up.
	w, _ := nsq.NewProducer(s.config.Address, config) // TODO: WRONG ADDRESS. This should be the producer
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, i)
	if err := w.Publish("images", buffer); err != nil {
		s.config.Logger.Debug().Err(err).Msg("Failed to publish image id...")
		return err
	}

	w.Stop()
	s.config.Logger.Debug().Msg("done")
	return nil
}
