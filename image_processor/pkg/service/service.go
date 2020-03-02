package service

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers"
)

// Config is everything that this service needs to work.
type Config struct {
}

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
func New(cfg Config, deps Dependencies) Service {
	return &imageProcessor{
		deps:   deps,
		config: cfg,
	}
}

// Service represents the service object of the receiver.
type imageProcessor struct {
	config Config
	deps   Dependencies
}

// Run starts the this service.
// TODO: Pass the context?
func (s *imageProcessor) Run(ctx context.Context) {
	s.deps.Logger.Info().Msg("Starting service...")
	done := make(chan struct{})

	// Create the channel on which the consumer and the processor can comunicate.
	// This should be buffered.
	mediator := make(chan int, 1)

	// Start the consumer routine.
	go s.deps.Consumer.Consume(mediator)

	// Start the processor routine.
	go s.deps.Processor.ProcessImages(mediator)

	// Block forever.
	<-done
}
