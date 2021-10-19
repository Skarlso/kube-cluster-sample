package images

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"

	"github.com/Skarlso/kube-cluster-sample/receiver/models"
)

// Config db configs here
type Config struct {
	Port             string
	Dbname           string
	UsernamePassword string
	Hostname         string
	Logger           zerolog.Logger
}

type imageProvider struct {
	config Config
}

// NewImageProvider creates a new image provider using the db.
func NewImageProvider(cfg Config) *imageProvider {
	return &imageProvider{config: cfg}
}

func (i *imageProvider) openConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s@tcp(%s:%s)/%s",
		i.config.UsernamePassword,
		i.config.Hostname,
		i.config.Port,
		i.config.Dbname)
	return sql.Open("mysql", connectionString)
}

// SaveImage takes an image model and saves it into the database.
func (i *imageProvider) SaveImage(image *models.Image) (*models.Image, error) {
	i.config.Logger.Debug().Str("path", string(image.Path)).Msg("Saving image path...")
	if ok, _ := i.SearchPath(string(image.Path)); ok {
		return nil, fmt.Errorf("image with path '%s' already exists", string(image.Path))
	}

	conn, err := i.openConnection()
	if err != nil {
		i.config.Logger.Debug().Err(err).Msg("failed to open connection")
		return nil, err
	}
	defer conn.Close()
	result, err := conn.Exec("insert into images (path, person, status) values (?, ?, ?)", image.Path, image.PersonID, image.Status)
	if err != nil {
		i.config.Logger.Debug().Err(err).Msg("failed to run insert")
		return nil, err
	}
	id, _ := result.LastInsertId()
	image.ID = id

	return image, nil
}

// LoadImage takes an id and looks for that id in the database.
func (i *imageProvider) LoadImage(id int64) (*models.Image, error) {
	log.Println("Loading image with ID: ", id)
	conn, err := i.openConnection()
	if err != nil {
		i.config.Logger.Debug().Err(err).Msg("failed to open connection")
		return &models.Image{}, err
	}
	defer conn.Close()
	var (
		imageID int
		path    string
		person  int
		status  int
	)
	err = conn.QueryRow("select id, path, person, status from images where id = ?", id).Scan(&imageID, &path, &person, status)
	if err != nil {
		i.config.Logger.Debug().Err(err).Msg("failed to query images")
		return &models.Image{}, err
	}
	ret := &models.Image{
		ID:       int64(imageID),
		Path:     []byte(path),
		PersonID: person,
		Status:   models.Status(status),
	}
	return ret, nil
}

// SearchPath takes an image path and checks if there is an image in the database with that path.
func (i *imageProvider) SearchPath(path string) (bool, error) {
	conn, err := i.openConnection()
	if err != nil {
		i.config.Logger.Debug().Err(err).Msg("failed to open connection")
		return false, err
	}
	defer conn.Close()
	row, err := conn.Query("select 1 from images where path = ?", path)
	if err != nil {
		i.config.Logger.Debug().Err(err).Msg("failed to select from images")
		return false, err
	}
	defer row.Close()
	return row.Next(), nil
}
