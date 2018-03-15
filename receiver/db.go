package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// DbConnection mysql connection. I'm not going to do pooling here.
type DbConnection struct {
	*sql.DB
}

func (dc *DbConnection) setup() error {
	// Create the Image table here
	err := dc.open()
	if err != nil {
		return err
	}
	defer dc.close()

	_, err = dc.Exec(`
		create table images(
			id int not null auto_increment primary key,
			path varchar(255) not null,
			person int,
			status int
		)
	`)
	if e, ok := err.(*mysql.MySQLError); ok {
		if e.Number == 1050 {
			fmt.Println("table already exists.")
			return nil
		}
	}
	return err
}

func (dc *DbConnection) open() error {
	if dc.DB != nil {
		return nil
	}
	connectionString := fmt.Sprintf("%s@tcp(%s:%d)/%s",
		configuration.MySQLUserPass,
		configuration.MySQLHostname,
		configuration.MySQLPort,
		configuration.MySQLDBName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	dc.DB = db
	return db.Ping()
}

func (dc *DbConnection) close() error {
	if dc.DB != nil {
		dc.DB.Close()
	}
	return errors.New("trying to close a nil connection")
}

// saveImage saves an image with an incremental id, and sets that ID
// up for that image after save.
func (dc *DbConnection) saveImage(i *Image) error {
	err := dc.open()
	if err != nil {
		return err
	}
	result, err := dc.Exec("insert into images (path, person, status) values (?, ?, ?)", i.Path, i.PersonID, i.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	i.ID = int(id)
	defer dc.close()
	return nil
}

// loadImage gets an image from the database.
func (dc *DbConnection) loadImage(id int) (Image, error) {
	err := dc.open()
	if err != nil {
		return Image{}, err
	}
	defer dc.close()
	var (
		imageID int
		path    string
		person  int
		status  int
	)
	err = dc.QueryRow("select id, path, person, status from images where id = ?", id).Scan(&imageID, &path, &person, status)
	if err != nil {
		return Image{}, err
	}
	i := Image{
		ID:       imageID,
		Path:     []byte(path),
		PersonID: person,
		Status:   Status(status),
	}
	return i, nil
}

func (dc *DbConnection) searchPath(path string) (bool, error) {
	err := dc.open()
	if err != nil {
		return false, err
	}
	defer dc.close()
	row, err := dc.Query("select 1 from images where path = ?", path)
	if err != nil {
		return false, err
	}
	defer row.Close()
	return row.Next(), nil
}
