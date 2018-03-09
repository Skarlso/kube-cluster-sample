package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alecthomas/template"
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
}

// Person is a person.
type Person struct {
	Name string
}

func init() {
	log.Println("Initiating configuration...")
	configuration = new(Configuration)
	ex, _ := os.Executable()
	configuration.loadConfiguration(filepath.Dir(ex))
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

	http.ListenAndServe(":8081", nil)
}
