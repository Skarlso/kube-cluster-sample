package providers

// ConsumerProvider interface which defines the minimum needed functionality for a consumer
// of an image queue.
type ConsumerProvider interface {
	Consume(sendTo chan int)
}
