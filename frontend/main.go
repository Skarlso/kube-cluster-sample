package main

import (
	"html/template"
	"log"
	"net/http"
)

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

// PageData returns the images that we would like to display.
type PageData struct {
	PageTitle string
	Images    []Image
}

// Image we'll let the DB assign an ID to an image.
type Image struct {
	ID     int
	Person Person
	Path   string
	Status Status
}

// Person is a person.
type Person struct {
	Name string
}

func init() {
	log.Println("Initiating environment...")
	initiateEnvironment()
	configuration = new(Configuration)
	configuration.loadConfiguration()
}

func view(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	db := new(DbConnection)
	images, err := db.loadImages()
	if err != nil {
		log.Fatal(err)
	}
	data := PageData{
		PageTitle: "Persons of Interest",
		Images:    images,
	}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", view)
	log.Printf("Started server under port: %s\n", configuration.Port)
	if err := http.ListenAndServe(":"+configuration.Port, nil); err != nil {
		log.Fatal(err)
	}
}
