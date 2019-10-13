package main

import (
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Declare and parse flags to be used
	configPath := flag.String("config", "~/.config/ghstatus/config.json", "configuration file path")

	organization := flag.String("org", "", "organization name")
	expiresIn := flag.String("expire", "", "status should expire after duration")
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

	config := &Config{Path: *configPath}

	// Check if config exists, otherwise use default config
	if !ConfigExists(*configPath) {
		config.Data = DefaultConfig
		err := config.Write(false)
		if err != nil {
			println("Failed to write to config:", err.Error())
			os.Exit(1)
		}
	} else {
		err := config.Load()
		if err != nil {
			println("Failed to load config:", err.Error())
			os.Exit(1)
		}
	}

	args := flag.Args()

	// Handle case if no args were supplied
	if len(args) == 0 {
		err := getCommand(config)
		if err != nil {
			fmt.Println(aurora.Red(err.Error()))
			os.Exit(1)
			return
		}
		os.Exit(0)
	}

	// Match first arg
	switch flag.Arg(0) {
	case "set":
		err := setCommand(config, organization, expiresIn, busy, args[1:])
		if err != nil {
			fmt.Println(aurora.Red(err.Error()))
			os.Exit(1)
			return
		}
	case "get":
		err := getCommand(config)
		if err != nil {
			fmt.Println(aurora.Red(err.Error()))
			os.Exit(1)
			return
		}
	case "config":
		configCommand(config)
		break
	case "help":
		fallthrough
	default:
		helpCommand()
		break
	}
}
