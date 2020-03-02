package processor

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/circuitbreaker"
)

// Person is a person
type Person struct {
	ID   int
	Name string
}

// Status is an Image status representation
type Status int

const (
	// PENDING -- not yet send to face recognition service
	PENDING Status = iota
	// PROCESSED -- processed by face recognition service; even if no person was found for the image
	PROCESSED
	// FAILEDPROCESSING -- for whatever reason the processing failed and this image is flagged for a retry
	FAILEDPROCESSING
)

// Config needed for the processor.
type Config struct {
	Port             string
	Dbname           string
	UsernamePassword string
	Hostname         string
	GrpcAddress      string
}

// Dependencies of the processor provider.
type Dependencies struct {
	Logger         zerolog.Logger
	CircuitBreaker circuitbreaker.CircuitBreaker
}

// processor defines a processor which uses a real database to store and process data.
type processor struct {
	deps    Dependencies
	conf    Config
	db      *sql.DB
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

// open opens a connection to the database.
func (p *processor) open() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s@tcp(%s:%s)/%s",
		p.conf.UsernamePassword,
		p.conf.Hostname,
		p.conf.Port,
		p.conf.Dbname)
	return sql.Open("mysql", connectionString)
}

// getPath returns the path information of the image.
func (p *processor) getPath(id int) (string, error) {
	var path string
	err := p.db.QueryRow("select path from images where id = ? and status = ?", id, PENDING).Scan(&path)
	return path, err
}

// updateImageWithFailedStatus updates a given image ID with failed status.
func (p *processor) updateImageWithFailedStatus(imageID int) error {
	return p.updateImage("update images set status = ? where id = ?", FAILEDPROCESSING, imageID)
}

// updateImageWithPerson updates a record with the person's ID to which it belongs to.
func (p *processor) updateImageWithPerson(personID, imageID int) error {
	return p.updateImage("update images set person = ?, status = ? where id = ?", personID, PROCESSED, imageID)
}

// updateImage takes a sql and arguments to it and performs an update of the image record.
// includes a row count. if no records were affected, it will return an error.
func (p *processor) updateImage(sql string, args ...interface{}) error {
	res, err := p.db.Exec(sql, args...)
	if err != nil {
		return err
	}
	rowCount, _ := res.RowsAffected()
	if rowCount == 0 {
		return fmt.Errorf("no rows were affected")
	}
	return nil
}

// getPersonFromImage returns the person matching this image.
func (p *processor) getPersonFromImage(img string) (Person, error) {
	var (
		name string
		id   int
	)
	err := p.db.QueryRow(`select person.name, person.id from person inner join person_images
					   as pi on person.id = pi.person_id where image_name = ?`, img).Scan(&name, &id)
	if err != nil {
		return Person{}, err
	}
	return Person{Name: name, ID: id}, nil
}

// ProcessImages takes a channel for input and waits on that channel for processable items.
// This channel must never be closed.
func (p *processor) ProcessImages(in chan int) {
	// continously get ids for image processing, block until something is received.
	for {
		i := <-in
		p.deps.Logger.Info().Int("image-id", i).Msg("Processing image...")

		// open db connection
		db, err := p.open()
		if err != nil {
			p.deps.Logger.Error().Err(err).Msg("failed to open db connection")
			continue
		}
		p.db = db
		path, err := p.getPath(i)
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
		person, err := p.getPersonFromImage(name)
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
