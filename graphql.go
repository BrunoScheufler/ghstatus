package main

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
)

func SendApiRequest(token, query string, variables map[string]interface{}) (interface{}, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	request := graphql.NewRequest(query)

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	for k, v := range variables {
		request.Var(k, v)
	}

	// Create new context
	ctx := context.Background()

	var responseData interface{}
	err := client.Run(ctx, request, &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not send GitHub GraphQL API request: %w", err)
	}

	return responseData, nil
}

// GitHub types

// Use omitempty for optional values
type UpdateStatusMutationInput struct {
	Emoji               string `json:"emoji,omitempty"`
	Message             string `json:"message,omitempty"`
	OrganizationId      string `json:"organizationId,omitempty"`
	LimitedAvailability bool   `json:"limitedAvailability,omitempty"`
	ExpiresAt           string `json:"expiresAt,omitempty"`
}

type UpdateUserStatusMutationResponse struct {
	ChangeUserStatus struct {
		Status UserStatus `json:"status"`
	} `json:"changeUserStatus"`
}

type RetrieveUserStatusQueryResponse struct {
	Viewer struct {
		Status UserStatus `json:"status"`
	} `json:"viewer"`
}

type UserStatus struct {
	CreatedAt                    string `json:"createdAt"`
	Emoji                        string `json:"emoji"`
	ExpiresAt                    string `json:"expiresAt"`
	ID                           string `json:"id"`
	IndicatesLimitedAvailability string `json:"indicatesLimitedAvailability"`
	Message                      string
	Organization                 struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	UpdatedAt string `json:"updatedAt"`
}

// Queries & Mutations

var statusFragment = `
fragment StatusFragment on UserStatus {
	createdAt
	updatedAt
	expiresAt
	message
	emoji
	indicatesLimitedAvailability
	organization {
		name
	}
}
`

var retrievalQuery = `
query StatusRetrievalQuery {
  viewer {
    status {
      ...StatusFragment
    }
  }
}
` + statusFragment

var updateMutation = `
mutation UpdateUserStatusMutation ($input: ChangeUserStatusInput!) {
  changeUserStatus(input: $input) {
    status {
      ...StatusFragment
    }
  }
}
` + statusFragment
