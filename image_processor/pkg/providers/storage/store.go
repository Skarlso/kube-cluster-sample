package storage

import (
	"context"
	"database/sql"
	"fmt"

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

type MySQLStorage struct {
	Config
}

func (m *MySQLStorage) GetPath(id int) (string, error) {
	var path string
	f := func(tx *sql.Tx) error {
		return tx.QueryRow("select path from images where id = ? and status = ?", id, models.PENDING).Scan(&path)
	}
	if err := m.execInTx(context.Background(), f); err != nil {
		return "", fmt.Errorf("failed to get path: %w", err)
	}
	return path, nil
}

func (m *MySQLStorage) UpdateImageStatus(id int, status models.Status) error {
	f := func(tx *sql.Tx) error {
		res, err := tx.Exec("update images set status = ? where id = ?", status, id)
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

func (m *MySQLStorage) UpdateImageWithPerson(id int, person int, status models.Status) error {
	f := func(tx *sql.Tx) error {
		// TODO: this is duplicating the image status code. Should be handled better.
		res, err := tx.Exec("update images set person = ?, status = ? where id = ?", person, status, id)
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
		return nil, fmt.Errorf("failed to get path: %w", err)
	}
	return &models.Person{Name: name, ID: id}, nil
}

var _ providers.ImageStorer = &MySQLStorage{}

// NewMySQLStorage creates a new Image storage.
func NewMySQLStorage(cfg Config) *MySQLStorage {
	return &MySQLStorage{
		Config: cfg,
	}
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
	defer tx.Rollback()
	defer db.Close()

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
