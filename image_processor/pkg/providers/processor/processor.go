package processor

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/models"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/circuitbreaker"
)

// Config needed for the processor.
type Config struct {
	GrpcAddress string
}

// Dependencies of the processor provider.
type Dependencies struct {
	Logger         zerolog.Logger
	CircuitBreaker circuitbreaker.CircuitBreaker
	Storer         providers.ImageStorer
}

// processor defines a processor which uses a real database to store and process data.
type processor struct {
	deps    Dependencies
	conf    Config
	client  facerecog.IdentifyClient
	hclient facerecog.HealthCheckClient
}

// NewProcessorProvider creates a new processor provider with an active grpc connection.
func NewProcessorProvider(cfg Config, deps Dependencies) (providers.ProcessorProvider, error) {
	conn, err := grpc.Dial(cfg.GrpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc on: %s", cfg.GrpcAddress)
	}
	c := facerecog.NewIdentifyClient(conn)
	h := facerecog.NewHealthCheckClient(conn)
	return &processor{
		deps:    deps,
		conf:    cfg,
		client:  c,
		hclient: h,
	}, nil
}

// updateImageWithFailedStatus updates a given image ID with failed status.
func (p *processor) updateImageWithFailedStatus(imageID int) error {
	return p.deps.Storer.UpdateImageStatus(imageID, models.FAILEDPROCESSING)
}

// updateImageWithPerson updates a record with the person's ID to which it belongs to.
func (p *processor) updateImageWithPerson(personID, imageID int) error {
	return p.deps.Storer.UpdateImageWithPerson(imageID, personID, models.PROCESSED)
}

// ProcessImages takes a channel for input and waits on that channel for processable items.
// This channel must never be closed.
func (p *processor) ProcessImages(in chan int) {
	// continuously get ids for image processing, block until something is received.
	for {
		i := <-in
		p.deps.Logger.Info().Int("image-id", i).Msg("Processing image...")

		path, err := p.deps.Storer.GetPath(i)
		if err != nil {
			p.deps.Logger.Error().Err(err).Int("image-id", i).Msg("error while getting path for image")
			// log the error then continue
			continue
		}
		p.deps.CircuitBreaker.SetCallF(func() (*facerecog.IdentifyResponse, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			r, err := p.client.Identify(ctx, &facerecog.IdentifyRequest{
				ImagePath: path,
			})
			return r, err
		})

		p.deps.CircuitBreaker.SetPingF(func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			_, err := p.hclient.HealthCheck(ctx, &facerecog.Empty{})
			return err != nil
		})
		r, err := p.deps.CircuitBreaker.Call()
		if err != nil {
			if err := p.updateImageWithFailedStatus(i); err != nil {
				p.deps.Logger.Error().Err(err).Msg("could not update image to failed status")
				continue
			}
			p.deps.Logger.Error().Err(err).Msg("image processing failed, updated image to failed status.")
			continue
		}
		name := r.GetImageName()
		if name == "not_found" {
			if err := p.updateImageWithFailedStatus(i); err != nil {
				p.deps.Logger.Error().Err(err).Msg("could not update image to failed status")
				continue
			}
			p.deps.Logger.Error().Msg("the person could not be identified")
			continue
		}
		p.deps.Logger.Info().Str("name", name).Msg("got name from face recog processor")
		person, err := p.deps.Storer.GetPersonFromImage(name)
		if err != nil {
			p.deps.Logger.Error().Err(err).Msg("could not retrieve person")
			continue
		}
		p.deps.Logger.Info().Str("person-name", person.Name).Msg("got person... updating record with person id")
		err = p.updateImageWithPerson(person.ID, i)
		if err != nil {
			p.deps.Logger.Error().Err(err).Msg("warning: could not update image record")
			continue
		}
		p.deps.Logger.Info().Str("name", name).Msg("done")
	}
}
