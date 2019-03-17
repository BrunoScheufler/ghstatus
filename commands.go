package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
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

--org [name] - Set organization to limit visibility of status
--busy - Set limited availability
`
	fmt.Println(fmt.Sprintf(helpMenu, aurora.Bold("ghstatus"), aurora.Bold("Available commands:"), aurora.Bold("Available arguments:")))
}

func setCommand(config *Config, organization *string, limited *bool, args []string) {
	if !validateTokenSet(config) {
		fmt.Println(aurora.Red("Please set your auth token for GitHub in the configuration file first!"))
		return
	}

	if len(args) < 2 {
		fmt.Println(aurora.Red("Please supply two arguments for the status emoji and message"))
		return
	}

	err := set(config, args[0], strings.Join(args[1:], " "), organization, limited)
	if err != nil {
		fmt.Println(aurora.Red(fmt.Sprintf("Failed to send status update: %v.", err.Error())))
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
