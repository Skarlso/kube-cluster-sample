package consumer

import (
	"encoding/binary"
	"log"

	"github.com/rs/zerolog"

	"github.com/nsqio/go-nsq"

	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers"
)

// Config configuration options for the consumer service.
type Config struct {
	NsqAddress string
}

// Dependencies of the consumer provider.
type Dependencies struct {
	Logger zerolog.Logger
}

type consumer struct {
	cfg Config
}

// NewConsumer creates a new consumer provider.
func NewConsumer(cfg Config) providers.ConsumerProvider {
	return &consumer{cfg: cfg}
}

// Consume consums an entry from NSQ.
func (c *consumer) Consume(sendTo chan int) {
	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("images", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		data := binary.LittleEndian.Uint64(message.Body)
		sendTo <- int(data)
		return nil
	}))
	err := q.ConnectToNSQLookupd(c.cfg.NsqAddress)
	if err != nil {
		// TODO: Find a better way... :)
		log.Panic(err)
	}
}
