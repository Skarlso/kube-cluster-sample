package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DbConnection mysql connection. I'm not going to do pooling here.
type DbConnection struct {
	*sql.DB
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

// loadImage gets an image from the database.
func (dc *DbConnection) loadImages() ([]Image, error) {
	images := make([]Image, 0)
	err := dc.open()
	if err != nil {
		return images, err
	}
	rows, err := dc.Query("select id, path, person, status from images")
	if err != nil {
		return images, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			imageID int
			path    string
			person  int
			status  int
		)
		if err := rows.Scan(&imageID, &path, &person, &status); err != nil {
			log.Println("fatal while trying to scan rows")
			return images, err
		}
		p := Person{Name: "Pending..."}
		if person != -1 {
			p, err = dc.getPerson(person)
			if err != nil {
				return images, err
			}
		}
		i := Image{
			ID:     imageID,
			Path:   path,
			Person: p,
			Status: Status(status),
		}
		images = append(images, i)
	}
	return images, nil
}

func (dc *DbConnection) getPerson(id int) (Person, error) {
	err := dc.open()
	if err != nil {
		return Person{}, err
	}
	var name string
	err = dc.QueryRow("select name from person where id = ?", id).Scan(&name)
	if err != nil {
		return Person{}, err
	}
	p := Person{
		Name: name,
	}
	return p, nil
}
