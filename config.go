package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var DefaultConfig = ConfigContents{Token: ""}

type ConfigContents struct {
	Token string `json:"token"`
}

type Config struct {
	data ConfigContents
	path string
}

func (c *Config) serialize() ([]byte, error) {
	return json.Marshal(&c.data)
}

func (c *Config) write(exists bool) error {
	serialized, err := c.serialize()
	if err != nil {
		return err
	}

	if !exists {
		err := os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(c.path, serialized, 0644)
}

func (c *Config) load() error {
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c.data)
}

func configExists(path string) bool {
	exists, err := fileExists(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return exists
}
