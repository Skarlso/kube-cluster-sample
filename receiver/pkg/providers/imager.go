package providers

import "github.com/Skarlso/kube-cluster-sample/receiver/models"

// ImageProvider defines functions which are used to handle images.
type ImageProvider interface {
	SaveImage(image *models.Image) (*models.Image, error)
	LoadImage(i int64) (*models.Image, error)
	SearchPath(path string) (bool, error)
}
