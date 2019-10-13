package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func GetCurrentStatus(c *Config) (*UserStatus, error) {
	// Construct and execute request
	variables := make(map[string]interface{})
	rawResponse, err := SendApiRequest(c.Data.Token, retrievalQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("could not send API request: %w", err)
	}

	responseData := RetrieveUserStatusQueryResponse{}
	err = json.Unmarshal(rawResponse, &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response data: %w", err)
	}

	return &responseData.Viewer.Status, nil
}

type UpdateStatusInput struct {
	Config *Config

	Emoji   string
	Message string

	ExpiresAt           *string
	Organization        *string
	LimitedAvailability *bool
}

func UpdateStatus(input *UpdateStatusInput) (*UserStatus, error) {
	// Validate emoji
	if !strings.HasPrefix(input.Emoji, ":") || !strings.HasSuffix(input.Emoji, ":") {
		return nil, errors.New("invalid emoji format, please supply a valid emoji")
	}

	// Construct and send query
	variables := make(map[string]interface{})

	updateInput := UpdateStatusMutationInput{
		Message: input.Message,
		Emoji:   input.Emoji,
	}

	// Add organization to variables
	if input.Organization != nil {
		// TODO add org support (requires another query to fetch organizationId by name)
		fmt.Println("Note: Supplying an organization is currently not supported")
	}

	// Add limitedAvailability to variables
	if input.LimitedAvailability != nil {
		variables["limitedAvailability"] = *input.LimitedAvailability
	}

	variables["input"] = updateInput

	rawResponse, err := SendApiRequest(input.Config.Data.Token, updateMutation, variables)
	if err != nil {
		return nil, fmt.Errorf("could not send API request: %w", err)
	}

	responseData := UpdateUserStatusMutationResponse{}
	err = json.Unmarshal(rawResponse, &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response data: %w", err)
	}

	return &responseData.ChangeUserStatus.Status, nil
}
