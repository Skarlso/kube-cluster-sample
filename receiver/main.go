package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	log.Println("Initiating environment...")
	initiateEnvironment()
	configuration = new(Configuration)
	configuration.loadConfiguration()
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
func PostImage(w http.ResponseWriter, r *http.Request) {
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
	PostImages(w, r)
}

// PostImages handles a post of an image. Saves it to the database
// and sends it to NSQ for further processing.
func PostImages(w http.ResponseWriter, r *http.Request) {
	var p Paths
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprintf(w, "got error while decoding request body: %s", err)
		return
	}
	fmt.Fprintf(w, "got paths: %+v\n", p)
	nsq := new(NSQ)
	for _, path := range p.Paths {
		image := Image{
			ID:       -1,
			PersonID: -1,
			Path:     []byte(path.Path),
			Status:   PENDING,
		}
		err = image.saveImage()
		if err != nil {
			fmt.Fprintf(w, "got error while saving image: %s; moving on to next...", err)
			continue
		}
		fmt.Fprintf(w, "image saved with id: %d\n", image.ID)
		err = nsq.sendImage(image)
		if err != nil {
			fmt.Fprintf(w, "error while sending image to queue: %s", err)
			continue
		}
		fmt.Fprintln(w, "image sent to nsq")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/image/post", PostImage).Methods("POST")
	router.HandleFunc("/images/post", PostImages).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
