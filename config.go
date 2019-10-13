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
	Data ConfigContents
	Path string
}

func (c *Config) Serialize() ([]byte, error) {
	return json.Marshal(&c.Data)
}

func (c *Config) Write(exists bool) error {
	serialized, err := c.Serialize()
	if err != nil {
		return err
	}

	if !exists {
		err := os.MkdirAll(filepath.Dir(c.Path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(c.Path, serialized, 0644)
}

func (c *Config) Load() error {
	data, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c.Data)
}

func ConfigExists(path string) bool {
	exists, err := fileExists(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return exists
}
