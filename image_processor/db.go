package main

// TODO: DB needs to be extracted

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
