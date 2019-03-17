package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	configPath := flag.String("config", "~/.config/ghstatus/config.json", "configuration file path")
	organization := flag.String("org", "", "organization name")
	busy := flag.Bool("busy", false, "limited availability")
	flag.Parse()

	// Check if config is in home directory
	if strings.HasPrefix(*configPath, "~") {
		*configPath = strings.Replace(*configPath, "~", getHomeDir(), 1)
	}

	// Convert config file path to absolute path
	if !filepath.IsAbs(*configPath) {
		absPath, err := filepath.Abs(*configPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		*configPath = absPath
	}

	config := &Config{path: *configPath}

	// Check if config exists, otherwise use default config
	if !configExists(*configPath) {
		config.data = DefaultConfig
		err := config.write(false)
		if err != nil {
			println("Failed to write to config:", err.Error())
			os.Exit(1)
		}
	} else {
		err := config.load()
		if err != nil {
			println("Failed to load config:", err.Error())
			os.Exit(1)
		}
	}

	args := flag.Args()

	// Handle case if no args were supplied
	if len(args) == 0 {
		getCommand(config)
		os.Exit(0)
	}

	// Match first arg
	switch args[0] {
	case "set":
		setCommand(config, organization, busy, args[1:])
	case "get":
		getCommand(config)
	case "config":
		configCommand(config)
		break
	case "help":
		helpCommand()
		break
	default:
		helpCommand()
		break
	}
}

func configExists(path string) bool {
	exists, err := fileExists(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return exists
}
