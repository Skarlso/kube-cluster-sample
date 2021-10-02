package models

// Image defines an image object.
type Image struct {
	ID     int
	Path   string
	Person int
	Status Status
}
