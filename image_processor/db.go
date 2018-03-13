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

func (dc *DbConnection) getPath(id int) (string, error) {
	err := dc.open()
	if err != nil {
		return "", err
	}
	var path string
	err = dc.QueryRow("select path from images where id = ?", id).Scan(&path)
	return path, err
}

func (dc *DbConnection) getPersonFromImage(img string) (Person, error) {
	err := dc.open()
	if err != nil {
		return Person{}, err
	}
	var (
		name string
		id   int
	)
	err = dc.QueryRow(`select person.name, person.id from person inner join person_images
					   as pi on person.id = pi.person_id where image_name = ?`, img).Scan(&name, &id)
	if err != nil {
		return Person{}, err
	}
	return Person{Name: name, ID: id}, nil
}

func (dc *DbConnection) updateImageWithPerson(personID, imageID int) error {
	err := dc.open()
	if err != nil {
		return err
	}
	res, err := dc.Exec("update images set person = ? where id = ?", personID, imageID)
	rowCount, _ := res.RowsAffected()
	if rowCount == 0 {
		log.Println("warning: no rows were affected")
	}
	return err
}
