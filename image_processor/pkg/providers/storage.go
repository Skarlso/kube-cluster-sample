package providers

import "github.com/Skarlso/kube-cluster-sample/image_processor/pkg/models"

// ImageStorer handles storing and updating images for the image processor.
type ImageStorer interface {
	GetPath(id int) (string, error)
	UpdateImageStatus(id int, status models.Status) error
	UpdateImageWithPerson(id int, person int, status models.Status) error
	GetPersonFromImage(image string) (*models.Person, error)
}
