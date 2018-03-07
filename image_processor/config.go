package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Configuration represent a db configuration.
type Configuration struct {
	NSQLookupAddress string `yaml:"nsq_lookup_address"`
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
