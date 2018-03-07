package main

import (
	"log"

	nsq "github.com/bitly/go-nsq"
)

func consume() {
	// wg := &sync.WaitGroup{}
	// wg.Add(1)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("images", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", string(message.Body))
		imageQueue = append(imageQueue, 1, 2, 3, 4, 5)
		log.Println(imageQueue)
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
