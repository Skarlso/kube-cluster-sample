package providers

import "github.com/Skarlso/kube-cluster-sample/image_processor/pkg/models"

// ImageStorer handles storing and updating images for the image processor.
type ImageStorer interface {
	GetImage(id int) (*models.Image, error)
	UpdateImage(id int, person int, status models.Status) error
	GetPersonFromImage(image string) (*models.Person, error)
}
