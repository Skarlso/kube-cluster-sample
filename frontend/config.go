package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Configuration represent a db configuration.
type Configuration struct {
	MySQLHostname   string `yaml:"mysql_hostname"`
	MySQLUserPass   string `yaml:"mysql_userpassword"`
	MySQLPort       int    `yaml:"mysql_port"`
	MySQLDBName     string `yaml:"mysql_dbname"`
	NSQAddress      string `yaml:"nsq_address"`
	ProducerAddress string `yaml:"producer_address"`
}

var configuration *Configuration

func (c *Configuration) loadConfiguration(path string) {
	data, err := ioutil.ReadFile(filepath.Join(path, ".config.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}
}
