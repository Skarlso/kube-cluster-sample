package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alecthomas/template"
)

type PageData struct {
	PageTitle string
	Images    []Image
}

// Image we'll let the DB assign an ID to an image
type Image struct {
	ID       int
	PersonID int
	Path     string
}

func init() {
	log.Println("Initiating configuration...")
	configuration = new(Configuration)
	ex, _ := os.Executable()
	configuration.loadConfiguration(filepath.Dir(ex))
}

func main() {
	db := new(DbConnection)
	tmpl := template.Must(template.ParseFiles("index.html"))
	images, _ := db.loadImages()
	log.Println(images)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			PageTitle: "Persons of Interest",
			Images:    images,
		}
		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":8081", nil)
}
