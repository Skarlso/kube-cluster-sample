package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/models"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers"
)

// Config contains database access configuration.
type Config struct {
	Port             string
	Dbname           string
	UsernamePassword string
	Hostname         string
}

// MySQLStorage represents a storage implementation using MySQL.
type MySQLStorage struct {
	Config
}

var _ providers.ImageStorer = &MySQLStorage{}

// NewMySQLStorage creates a new Image storage.
func NewMySQLStorage(cfg Config) *MySQLStorage {
	return &MySQLStorage{
		Config: cfg,
	}
}

// GetImage returns an image.
func (m *MySQLStorage) GetImage(id int) (*models.Image, error) {
	var (
		path   string
		person int
		status int
	)
	f := func(tx *sql.Tx) error {
		return tx.QueryRow("select path, person, status from images where id = ?", id).Scan(&path, &person, &status)
	}
	if err := m.execInTx(context.Background(), f); err != nil {
		return nil, fmt.Errorf("failed to get path: %w", err)
	}
	return &models.Image{
		ID:     id,
		Path:   path,
		Person: person,
		Status: models.Status(status),
	}, nil
}

// UpdateImage updates an image.
// There is no check if the person exists or not, because that is happening in GetPersonFromImage.
func (m *MySQLStorage) UpdateImage(id int, person int, status models.Status) error {
	f := func(tx *sql.Tx) error {
		sets := []string{"status = ?"}
		args := []interface{}{status}
		if person != -1 {
			sets = append(sets, "person = ?")
			args = append(args, person)
		}
		args = append(args, id)
		res, err := tx.Exec(fmt.Sprintf("update images set %s where id = ?", strings.Join(sets, ", ")), args...)
		if err != nil {
			return fmt.Errorf("failed to update images: %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get affected rows: %w", err)
		}
		if rows == 0 {
			return fmt.Errorf("affect rows was zero")
		}
		return nil
	}
	if err := m.execInTx(context.Background(), f); err != nil {
		return fmt.Errorf("failed to update images with status in transaction: %w", err)
	}
	return nil
}

// GetPersonFromImage returns a person from an image url.
func (m *MySQLStorage) GetPersonFromImage(image string) (*models.Person, error) {
	var (
		name string
		id   int
	)
	f := func(tx *sql.Tx) error {
		return tx.QueryRow(`select person.name, person.id from person inner join person_images
					   as pi on person.id = pi.person_id where image_name = ?`, image).Scan(&name, &id)
	}
	if err := m.execInTx(context.Background(), f); err != nil {
		return nil, fmt.Errorf("failed to get person: %w", err)
	}
	return &models.Person{Name: name, ID: id}, nil
}

// execInTx executes in transaction. It will either commit, or rollback if there was an error.
func (m *MySQLStorage) execInTx(ctx context.Context, f func(tx *sql.Tx) error) error {
	db, err := m.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	// Defer a rollback in case anything fails.
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println("Failed to rollback: ", err)
		}
	}()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close database: ", err)
		}
	}()

	if err := f(tx); err != nil {
		return fmt.Errorf("failed to run function in transaction: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (m *MySQLStorage) createConnectionString() string {
	return fmt.Sprintf("%s@tcp(%s:%s)/%s",
		m.UsernamePassword,
		m.Hostname,
		m.Port,
		m.Dbname)
}

func (m *MySQLStorage) connect() (*sql.DB, error) {
	return sql.Open("mysql", m.createConnectionString())
}
