package main

import (
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

	responseData, ok := rawResponse.(RetrieveUserStatusQueryResponse)
	if !ok {
		return nil, errors.New("could not cast response data into expected type")
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

	responseData, ok := rawResponse.(UpdateUserStatusMutationResponse)
	if !ok {
		return nil, errors.New("could not cast response into expected type")
	}

	updatedStatus := responseData.ChangeUserStatus.Status

	return &updatedStatus, nil
}
