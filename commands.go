package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func configCommand(config *Config) {
	fmt.Println(fmt.Sprintf("The current configuration is located at %v.", aurora.Cyan(config.path).String()))
}

func helpCommand() {
	helpMenu := `
%v

%v
get - Retrieve current status
set [emoji] [status] - Set new status

%v
--config [path to config file] (default: ~/.config/ghstatus/config.json) - Set config path
--org | -org [name] - Set organization to limit visibility of status
--limited | -l [true/false] - Set limited availability / busy
`
	fmt.Println(fmt.Sprintf(helpMenu, aurora.Bold("ghstatus"), aurora.Bold("Available commands:"), aurora.Bold("Available arguments:")))
}

func setCommand(config *Config) {
	if !validateTokenSet(config) {
		fmt.Println(aurora.Red("Please set your auth token for GitHub in the configuration file first!"))
	}
}

func getCommand(config *Config) {
	if !validateTokenSet(config) {
		fmt.Println(aurora.Red("Please set your auth token for GitHub in the configuration file first!"))
		return
	}

	err := get(config)
	if err != nil {
		fmt.Println(aurora.Red(fmt.Sprintf("Failed to send status request: %v.", err.Error())))
	}
}
