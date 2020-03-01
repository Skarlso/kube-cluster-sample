package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/rs/zerolog"

	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/circuitbreaker"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/consumer"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/processor"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/service"
)

var (
	rootArgs struct {
		service         service.Config
		consumerConfig  consumer.Config
		processorConfig processor.Config
	}
)

func init() {
	flag.StringVar(&rootArgs.processorConfig.Hostname, "db-host", "localhost", "--db-host=localhost")
	flag.StringVar(&rootArgs.processorConfig.UsernamePassword, "db-username-password", "root:password123", "--db-username-password=root:password123")
	flag.StringVar(&rootArgs.processorConfig.Dbname, "db-dbname", "kube", "--db-dbname=kube")
	flag.StringVar(&rootArgs.processorConfig.Port, "db-port", "3306", "--db-port=3306")
	flag.StringVar(&rootArgs.processorConfig.GrpcAddress, "grpc-address", "localhost:50051", "--grpc-address=localhost:50051")
	flag.StringVar(&rootArgs.consumerConfig.NsqAddress, "nsq-lookup-address", "127.0.0.1:4161", "--nsq-lookup-address=127.0.0.1:4161")
	flag.Parse()
}

func main() {
	// Wire up the service and its dependencies.
	cb := circuitbreaker.NewcircuitBreaker()
	logger := zerolog.New(os.Stdout)
	proc, err := processor.NewProcessorProvider(rootArgs.processorConfig, processor.Dependencies{
		CircuitBreaker: cb,
		Logger:         logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	cons := consumer.NewConsumer(rootArgs.consumerConfig)

	srvc := service.New(rootArgs.service, service.Dependencies{
		Processor: proc,
		Consumer:  cons,
	})

	srvc.Run(context.Background())
}
