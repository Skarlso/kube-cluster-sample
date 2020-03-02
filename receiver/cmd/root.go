package main

import (
	"context"
	"flag"
	"log"

	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/images"
	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/sender"
	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/service"
)

var (
	rootArgs struct {
		service      service.Config
		imgConfig    images.Config
		senderConfig sender.Config
	}
)

func init() {
	flag.StringVar(&rootArgs.imgConfig.Hostname, "db-host", "localhost", "--db-host=localhost")
	flag.StringVar(&rootArgs.imgConfig.UsernamePassword, "db-username-password", "root:password123", "--db-username-password=root:password123")
	flag.StringVar(&rootArgs.imgConfig.Dbname, "db-dbname", "kube", "--db-dbname=kube")
	flag.StringVar(&rootArgs.imgConfig.Port, "db-port", "3306", "--db-port=3306")
	flag.StringVar(&rootArgs.senderConfig.Address, "producer-address", "127.0.0.1:4150", "--producer-address=127.0.0.1:4150")
	flag.Parse()
}

func main() {
	imgProvider := images.NewImageProvider(rootArgs.imgConfig)
	senderProvider := sender.NewNSQSender(rootArgs.senderConfig)

	srvc := service.New(rootArgs.service, service.Dependencies{
		ImageProvider: imgProvider,
		SendProvider:  senderProvider,
	})

	if err := srvc.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
