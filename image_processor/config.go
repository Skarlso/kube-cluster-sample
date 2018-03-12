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
	MySQLHostname    string
	MySQLUserPass    string
	MySQLPort        int
	MySQLDBName      string
	NSQLookupAddress string
}

var configuration *Configuration

func (c *Configuration) loadConfiguration() {
	c.MySQLHostname = os.Getenv("MYSQL_CONNECTION")
	c.MySQLDBName = os.Getenv("MYSQL_DBNAME")
	c.MySQLPort, _ = strconv.Atoi(os.Getenv("MYSQL_PORT"))
	c.MySQLUserPass = os.Getenv("MYSQL_USERPASSWORD")
	c.NSQLookupAddress = os.Getenv("NSQ_LOOKUP_ADDRESS")
}

func initiateEnvironment() {
	ex, _ := os.Executable()
	path := filepath.Dir(ex)
	data, err := ioutil.ReadFile(filepath.Join(path, ".env"))
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
			os.Setenv(k, v)
		}
	}
}
