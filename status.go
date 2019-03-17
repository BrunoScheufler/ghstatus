package main

import (
	"errors"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/mapstructure"
)

func get(c *Config) error {
	// Construct and execute request
	variables := map[string]interface{}{}
	response, err := sendAPIRequest(c.data.Token, retrievalQuery, variables)
	if err != nil {
		return err
	}

	// Handle potential errors
	if len(response.Errors) > 0 {
		for _, e := range response.Errors {
			error := GraphQLError{}

			err = mapstructure.Decode(e, &error)
			if err != nil {
				continue
			}

			println(error.Message)
		}

		return errors.New("request failed")
	}

	responseData := RetrievalQueryResponseData{}

	// Try to decode body
	err = mapstructure.Decode(response.Data, &responseData)
	if err != nil {
		return err
	}

	status := responseData.Viewer.Status

	// Print status
	fmt.Println(aurora.Bold("Current status:"), emoji.Sprint(status.Emoji), status.Message)

	// Check if visibility is limited to organization
	if status.Organization.Name != "" {
		fmt.Println(fmt.Sprintf("Status visible to %v.", aurora.Bold(status.Organization.Name)))
	}

	// Print availability
	if status.IndicatesLimitedAvailability == true {
		fmt.Println(aurora.Bold("Your availability is marked as limited."))
	}

	return nil
}
