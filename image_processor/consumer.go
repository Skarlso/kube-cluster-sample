package main

import (
	"encoding/binary"
	"log"

	nsq "github.com/bitly/go-nsq"
)

func consume() {
	// wg := &sync.WaitGroup{}
	// wg.Add(1)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("images", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		data := binary.LittleEndian.Uint64(message.Body)
		log.Printf("Got a message: %d\n", data)
		imageQueue.add(int(data))
		c.Signal()
		// wg.Done()
		return nil
	}))
	err := q.ConnectToNSQLookupd(configuration.NSQLookupAddress)
	if err != nil {
		log.Panic("Could not connect")
	}
	// wg.Wait()
}
