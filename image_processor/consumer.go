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
		log.Printf("Got a message: %v", message)
		// wg.Done()
		return nil
	}))
	err := q.ConnectToNSQLookupd(configuration.NSQLookupAddress)
	// err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}
	// wg.Wait()
}
