package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
	"strings"
)

func GetCurrentStatus(c *Config) (string, error) {
	// Construct and execute request
	variables := make(map[string]interface{})
	response, err := sendAPIRequest(c.Data.Token, retrievalQuery, variables)
	if err != nil {
		return "", fmt.Errorf("could not send API request: %w", err)
	}

	log.Println(response)

	// Handle GraphQL errors
	//err = handleGraphQLErrors(response)
	//if err != nil {
	//	return err
	//}
	//
	//responseData := RetrievalQueryResponseData{}
	//
	//// Try to decode body
	//err = mapstructure.Decode(response.Data, &responseData)
	//if err != nil {
	//	return err
	//}
	//
	//status := responseData.Viewer.Status
	//
	//// Print status
	//fmt.Println(aurora.Bold("Current status:"), emoji.Sprint(status.Emoji), status.Message)
	//
	//// Check if visibility is limited to organization
	//if status.Organization.Name != "" {
	//	fmt.Println(fmt.Sprintf("Status visible to %v.", aurora.Bold(status.Organization.Name)))
	//}
	//
	//// Print availability
	//if status.IndicatesLimitedAvailability == true {
	//	fmt.Println(aurora.Bold("Your availability is marked as limited."))
	//}
	//
	//return nil
	return "test", nil
}

type UpdateStatusInput struct {
	Config *Config

	Emoji   string
	Message string

	ExpiresAt           *string
	Organization        *string
	LimitedAvailability *bool
}

func UpdateStatus(input *UpdateStatusInput) error {
	// Validate emoji
	if !strings.HasPrefix(input.Emoji, ":") || !strings.HasSuffix(input.Emoji, ":") {
		return errors.New("invalid emoji format, please supply a valid emoji")
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
		return fmt.Errorf("could not send API request: %w", err)
	}

	responseData, ok := rawResponse.(UpdateUserStatusMutationResponse)
	if !ok {
		return fmt.Errorf("could not cast response into expected type: %w", err)
	}

	updatedStatus := responseData.ChangeUserStatus.Status

	// TODO Use template here
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
