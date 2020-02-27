package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Skarlso/kube-cluster-sample/receiver/models"
	"github.com/Skarlso/kube-cluster-sample/receiver/pkg/providers"
)

// Config is everything that this service needs to work.
type Config struct {
	Nsq struct {
		Address string
	}
	Producer struct {
		Address string
	}
}

// Dependencies are providers which this service operates with.
type Dependencies struct {
	ImageProvider providers.ImageProvider
	SendProvider  providers.SendProvider
}

// Service interface defines a service which can Run something.
type Service interface {
	Run(ctx context.Context) error
}

// Service represents the service object of the receiver.
type receiver struct {
	config Config
	deps   Dependencies
}

// Path is a single path of an image to process.
type Path struct {
	Path string `json:"path"`
}

// Paths is a batch of paths to process.
type Paths struct {
	Paths []Path `json:"paths"`
}

// PostImage handles a post of an image. Saves it to the database
// and sends it to NSQ for further processing.
func (s *receiver) postImage(w http.ResponseWriter, r *http.Request) {
	var p Path
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprintf(w, "got error while decoding body: %s", err)
		return
	}
	fmt.Fprintf(w, "got path: %+v\n", p)
	var ps Paths
	paths := make([]Path, 0)
	paths = append(paths, p)
	ps.Paths = paths
	var pathsJSON bytes.Buffer
	err = json.NewEncoder(&pathsJSON).Encode(ps)
	if err != nil {
		fmt.Fprintf(w, "failed to encode paths: %s", err)
		return
	}
	r.Body = ioutil.NopCloser(&pathsJSON)
	r.ContentLength = int64(pathsJSON.Len())
	s.postImages(w, r)
}

// PostImages handles a post of an image. Saves it to the database
// and sends it to NSQ for further processing.
func (s *receiver) postImages(w http.ResponseWriter, r *http.Request) {
	var p Paths
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprintf(w, "got error while decoding request body: %s", err)
		return
	}
	fmt.Fprintf(w, "got paths: %+v\n", p)
	for _, path := range p.Paths {
		image := models.Image{
			ID:       -1,
			PersonID: -1,
			Path:     []byte(path.Path),
			Status:   models.PENDING,
		}
		savedImage, err := s.deps.ImageProvider.SaveImage(&image)
		if err != nil {
			fmt.Fprintf(w, "got error while saving image: %s; moving on to next...", err)
			continue
		}
		fmt.Fprintf(w, "image saved with id: %d\n", savedImage.ID)
		err = s.deps.SendProvider.SendImage(uint64(savedImage.ID))
		if err != nil {
			fmt.Fprintf(w, "error while sending image to queue: %s", err)
			continue
		}
		fmt.Fprintln(w, "image sent to nsq")
	}
}

// New creates a new service will all its needed configuration.
func New(cfg Config, deps Dependencies) Service {
	s := &receiver{
		config: cfg,
		deps:   deps,
	}
	return s
}

// Run starts this service.
func (s *receiver) Run(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/image/post", s.postImage).Methods("POST")
	router.HandleFunc("/images/post", s.postImages).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
	return nil
}
