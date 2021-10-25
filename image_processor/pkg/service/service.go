package service

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers"
)

// Dependencies are providers which this service operates with.
type Dependencies struct {
	Consumer  providers.ConsumerProvider
	Processor providers.ProcessorProvider
	Logger    zerolog.Logger
}

// Service interface defines a service which can Run something.
type Service interface {
	Run(ctx context.Context)
}

// New creates a new service with all of its dependencies and configurations.
func New(deps Dependencies) Service {
	return &imageProcessor{
		deps: deps,
	}
}

// Service represents the service object of the receiver.
type imageProcessor struct {
	deps Dependencies
}

// Run starts the this service.
// TODO: Pass the context?
func (s *imageProcessor) Run(ctx context.Context) {
	s.deps.Logger.Info().Msg("Starting service...")
	done := make(chan struct{})

	// Create the channel on which the consumer and the processor can communicate.
	// This should be buffered.
	mediator := make(chan int, 1)

	// Start the consumer routine.
	go s.deps.Consumer.Consume(mediator)

	// Start the processor routine.
	go s.deps.Processor.ProcessImages(context.Background(), mediator)

	// Block forever.
	<-done
}
