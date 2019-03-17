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

	// Handle GraphQL errors
	err = handleGraphQLErrors(response)
	if err != nil {
		return err
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

func set(c *Config, emoji, message string, organization *string, limitedAvailability *bool) error {
	// Construct and send query
	variables := map[string]interface{}{"emoji": emoji, "message": message}

	if organization != nil && *organization != "" {
		// TODO add org support (requires another query to fetch organizationId by name)
		fmt.Println("Note: Supplying an organization is currently not supported")
	}

	if limitedAvailability != nil {
		variables["limitedAvailability"] = *limitedAvailability
	}

	response, err := sendAPIRequest(c.data.Token, updateMutation, map[string]interface{}{"newStatus": variables})
	if err != nil {
		return err
	}

	err = handleGraphQLErrors(response)
	if err != nil {
		return err
	}

	responseData := UpdateMutationResponseData{}

	// Try to decode body
	err = mapstructure.Decode(response.Data, &responseData)
	if err != nil {
		return err
	}

	status := responseData.ChangeUserStatus.Status

	if status.Message != message ||
		status.Emoji != emoji ||
		(organization != nil && status.Organization.Name != *organization) ||
		(limitedAvailability != nil && status.IndicatesLimitedAvailability != *limitedAvailability) {
		return errors.New("some fields were not updated accordingly, please try again")
	}

	fmt.Println(aurora.Green("🎉 Updated your status!"))

	return nil
}
