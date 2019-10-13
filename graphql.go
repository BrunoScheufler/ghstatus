package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func SendApiRequest(token, query string, variables map[string]interface{}) (json.RawMessage, error) {
	requestHeaders := make(map[string]string)
	requestHeaders["Authorization"] = fmt.Sprintf("Bearer %s", token)
	requestHeaders["Content-Type"] = "application/json"

	requestBody := GraphQLRequestBody{
		Query:     query,
		Variables: variables,
	}

	requestInput := SendRequestInput{
		Endpoint: "https://api.github.com/graphql",
		Headers:  requestHeaders,
		Body:     &requestBody,
	}

	response, err := SendGraphQLRequest(&requestInput)
	if err != nil {
		return nil, fmt.Errorf("could not send GraphQL request: %w", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.New(fmt.Sprintf("response contained errors: %s", response.Errors[0].Message))
	}

	return response.Data, nil
}

type GraphQLRequestBody struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName,omitempty"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponseBody struct {
	Errors []GraphQLError  `json:"errors"`
	Data   json.RawMessage `json:"data"`
}

type SendRequestInput struct {
	Endpoint string
	Headers  map[string]string
	Body     *GraphQLRequestBody
}

func SendGraphQLRequest(input *SendRequestInput) (*GraphQLResponseBody, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	body, err := json.Marshal(input.Body)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %w", err)
	}

	request, err := http.NewRequest("POST", input.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	for k, v := range input.Headers {
		request.Header.Set(k, v)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request returned non-200 response status")
	}

	rawResponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	responseBody := GraphQLResponseBody{}

	err = json.Unmarshal(rawResponseBody, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &responseBody, nil
}

// GitHub types

// Use omitempty for optional values

type OrganizationLookupQueryResponse struct {
	Organization struct {
		ID string `json:"id"`
	} `json:"organization"`
}

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
	IndicatesLimitedAvailability bool   `json:"indicatesLimitedAvailability"`
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

var organizationLookupQuery = `
query OrganizationLookupQuery ($input: String!){
  organization(login: $input) {
    id
  }
}
`
