package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Configuration represent a db configuration.
type Configuration struct {
	MySQLHostname   string
	MySQLUserPass   string
	MySQLPort       int
	MySQLDBName     string
	NSQAddress      string
	ProducerAddress string
}

var configuration *Configuration

func (c *Configuration) loadConfiguration() {
	c.MySQLHostname = os.Getenv("MYSQL_CONNECTION")
	c.MySQLDBName = os.Getenv("MYSQL_DBNAME")
	c.MySQLPort, _ = strconv.Atoi(os.Getenv("MYSQL_PORT"))
	c.NSQAddress = os.Getenv("NSQ_ADDRESS")
	c.ProducerAddress = os.Getenv("PRODUCER_ADDRESS")
}

func initiateEnvironment() {
	ex, _ := os.Executable()
	path := filepath.Dir(ex)
	envPath := filepath.Join(path, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Println(".env doesn't exists. moving on assuming env is setup.")
		return
	}
	data, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.Contains(line, "=") {
			continue
		}
		split := strings.Split(line, "=")
		k, v := split[0], split[1]
		if _, ok := os.LookupEnv(k); !ok {
			_ = os.Setenv(k, v)
		}
	}
}
