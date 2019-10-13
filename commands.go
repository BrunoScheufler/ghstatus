package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
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

	if len(args) < 2 {
		return errors.New("Please supply at least two arguments for the status emoji and message")
	}

	updateStatusInput := UpdateStatusInput{
		Config:              config,
		Emoji:               args[0],
		Message:             strings.Join(args[1:], " "),
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

	fmt.Println(fmt.Sprintf(`✅ Successfully updated your status:
Status: %s %s
Busy: %s 
Organization: %s
Expires At: %s
`,
		updatedStatus.Emoji,
		aurora.Bold(updatedStatus.Message),
		aurora.Bold(updatedStatus.IndicatesLimitedAvailability),
		aurora.Bold(updatedStatus.Organization.Name),
		aurora.Bold(updatedStatus.ExpiresAt)),
	)

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

	// TODO Format and print status with Go template
	log.Println(fmt.Sprintf("Current status: %s", status.Message))

	return nil
}
