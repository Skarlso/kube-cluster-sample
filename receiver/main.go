package main

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.Println("Initiating configuration...")
	configuration = new(Configuration)
	ex, _ := os.Executable()
	configuration.loadConfiguration(filepath.Dir(ex))
}

func main() {
	// log.Println("Starting up microservice one.")
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// })

	// log.Fatal(http.ListenAndServe(":8080", nil))
	db := new(DbConnection)
	err := db.setup()
	if err != nil {
		log.Fatal(err)
	}
	image := Image{
		ID:       -1,
		Path:     []byte("testpath"),
		PersonID: -1,
	}
	err = image.saveImage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("saved image with ID: ", image.ID)
	i := new(Image)
	i.loadImage(3)
	log.Println("loaded image: ", string(i.Path))
}
