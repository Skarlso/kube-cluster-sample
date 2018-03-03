package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
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
			person int
		)
	`)
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

func (dc *DbConnection) saveImage(i *Image) error {
	err := dc.open()
	if err != nil {
		return err
	}
	defer dc.close()
	return nil
}

func (dc *DbConnection) loadImage(id int) (*Image, error) {
	err := dc.open()
	if err != nil {
		return nil, err
	}
	defer dc.close()
	return &Image{}, nil
}
