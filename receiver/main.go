package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func init() {
	log.Println("Initiating configuration...")
	configuration = new(Configuration)
	ex, _ := os.Executable()
	configuration.loadConfiguration(filepath.Dir(ex))
}

// Path is a path to an image
type Path struct {
	Path string `json:"path"`
}

// PostImage handles a post of an image. Saves it to the database
// and sends it to NSQ for further processing.
func PostImage(w http.ResponseWriter, r *http.Request) {
	var p Path
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprintf(w, "got error: %s", err)
		return
	}
	fmt.Fprintf(w, "got path: %+v\n", p)
	image := Image{
		ID:       -1,
		PersonID: -1,
		Path:     []byte(p.Path),
	}
	err = image.saveImage()
	fmt.Fprintf(w, "image saved with id: %d\n", image.ID)
	if err != nil {
		fmt.Fprintf(w, "got error while saving image: %s", err)
		return
	}
	nsq := new(NSQ)
	err = nsq.sendImage(image)
	if err != nil {
		fmt.Fprintf(w, "error while sending image to queue: %s", err)
		return
	}
	fmt.Fprintln(w, "image sent to nsq")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/image/post", PostImage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
