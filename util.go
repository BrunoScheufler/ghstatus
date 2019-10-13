package main

import (
	"os"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		println("Could not retrieve home directory")
		os.Exit(1)
	}

	return homedir
}
