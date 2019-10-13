package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
)

func configCommand(config *Config) {
	fmt.Println(fmt.Sprintf("The current configuration is located at %v.", aurora.Cyan(config.Path)))
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

func setCommand(config *Config, organization *string, limited *bool, args []string) error {
	if config.Data.Token == "" {
		return errors.New("Please set your auth token for GitHub in the configuration file first!")
	}

	message := ""
	emoji := ""
	if len(args) > 1 {
		emoji = args[0]
		message = strings.Join(args[1:], " ")
	}

	updateStatusInput := UpdateStatusInput{
		Config:              config,
		Emoji:               emoji,
		Message:             message,
		Organization:        organization,
		LimitedAvailability: limited,
	}
	updatedStatus, err := UpdateStatus(&updateStatusInput)
	if err != nil {
		return fmt.Errorf("could not send status update: %w", err)
	}

	if updatedStatus == nil {
		return errors.New("could not retrieve updated status")
	}

	formattedStatus, err := FormatStatus(updatedStatus)
	if err != nil {
		return fmt.Errorf("could not format status: %w", err)
	}

	fmt.Println(formattedStatus)

	return nil
}

func getCommand(config *Config) error {
	if config.Data.Token == "" {
		return errors.New("Please set your auth token for GitHub in the configuration file first!")
	}

	fmt.Println(aurora.Gray("Retrieving your current status..."))

	status, err := GetCurrentStatus(config)
	if err != nil {
		return fmt.Errorf("could not send request for current status: %w", err)
	}

	if status == nil {
		return errors.New("could noot retrieve updated status from GitHub API")
	}

	formattedStatus, err := FormatStatus(status)
	if err != nil {
		return fmt.Errorf("could not format status: %w", err)
	}

	fmt.Println(formattedStatus)

	return nil
}
