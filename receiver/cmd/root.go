package cmd

import (
	"context"
	"log"

	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/sender"

	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/images"
	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/service"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the root (and only) command of this service
	RootCmd = &cobra.Command{
		Use:   "receiver",
		Short: "Receiver service",
		Run:   runRootCmd,
	}

	rootArgs struct {
		service      service.Config
		imgConfig    images.Config
		senderConfig sender.Config
	}
)

func init() {
	f := RootCmd.Flags()
	f.StringVar(&rootArgs.imgConfig.Hostname, "db-host", "localhost", "--db-host=localhost")
	f.StringVar(&rootArgs.imgConfig.Password, "db-password", "password123", "--db-password=password123")
	f.StringVar(&rootArgs.imgConfig.Username, "db-username", "root", "--db-username=root")
	f.StringVar(&rootArgs.imgConfig.Dbname, "db-dbname", "kube", "--db-dbname=kube")
	f.StringVar(&rootArgs.imgConfig.Port, "db-port", "3306", "--db-port=3306")
	f.StringVar(&rootArgs.senderConfig.Address, "nsq-address", "127.0.0.1", "--nsq-address=127.0.0.1")
	f.StringVar(&rootArgs.service.Producer.Address, "producer-address", "127.0.0.1:4150", "--producer-address=127.0.0.1:4150")
}

func runRootCmd(cmd *cobra.Command, args []string) {
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
