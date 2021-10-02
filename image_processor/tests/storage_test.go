package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/models"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/storage"
)

func TestGetImage(t *testing.T) {
	m := storage.NewMySQLStorage(storage.Config{
		Port:             dbPort,
		Dbname:           "kube",
		UsernamePassword: "root:password123",
		Hostname:         "localhost",
	})

	t.Run("image for image id", func(tt *testing.T) {
		image, err := m.GetImage(1)
		assert.NoError(tt, err, "Error received when trying to get path: ", err)
		assert.Equal(tt, "test/path", image.Path)
		assert.Equal(tt, 1, image.Person)
		assert.Equal(tt, models.PENDING, image.Status)
	})

	t.Run("not existing image", func(tt *testing.T) {
		image, err := m.GetImage(99)
		assert.Error(tt, err, "should have received error")
		assert.Nil(tt, image)
	})
}

func TestUpdateImage(t *testing.T) {
	m := storage.NewMySQLStorage(storage.Config{
		Port:             dbPort,
		Dbname:           "kube",
		UsernamePassword: "root:password123",
		Hostname:         "localhost",
	})

	t.Run("just update the image status without updating the person", func(tt *testing.T) {
		image, err := m.GetImage(4)
		assert.NoError(tt, err)
		assert.Equal(tt, models.PENDING, image.Status)
		err = m.UpdateImage(4, -1, models.FAILEDPROCESSING)
		assert.NoError(tt, err)
		image, err = m.GetImage(4)
		assert.NoError(tt, err)
		assert.Equal(tt, models.FAILEDPROCESSING, image.Status)
	})

	t.Run("update image person and status", func(tt *testing.T) {
		err := m.UpdateImage(4, 3, models.PROCESSED)
		assert.NoError(tt, err)
		image, err := m.GetImage(4)
		assert.NoError(tt, err)
		assert.Equal(tt, models.PROCESSED, image.Status)
		assert.Equal(tt, 3, image.Person)
	})

	t.Run("updating a non existing image", func(tt *testing.T) {
		err := m.UpdateImage(99, 3, models.PROCESSED)
		assert.Error(tt, err)
	})
}

func TestGetPersonFromImage(t *testing.T) {
	m := storage.NewMySQLStorage(storage.Config{
		Port:             dbPort,
		Dbname:           "kube",
		UsernamePassword: "root:password123",
		Hostname:         "localhost",
	})

	t.Run("can retrieve a person by image name", func(tt *testing.T) {
		person, err := m.GetPersonFromImage("hannibal_1.jpg")
		assert.NoError(tt, err)
		assert.Equal(tt, "Hannibal", person.Name)
		person, err = m.GetPersonFromImage("john_doe_1.jpg")
		assert.NoError(tt, err)
		assert.Equal(tt, "John Doe", person.Name)
	})

	t.Run("fails to get a person if there is no associated image name", func(tt *testing.T) {
		person, err := m.GetPersonFromImage("unknown_image_1.jpg")
		assert.Error(tt, err)
		assert.Nil(tt, person)
	})
}
