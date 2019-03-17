package main

type GraphQLRequestBody struct {
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
	OperationName *string     `json:"operationName"`
}

type GraphQLResponseBody struct {
	Data   map[string]interface{}   `json:"data"`
	Errors []map[string]interface{} `json:"errors"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

// GitHub types

type Org struct {
	Name string `json:"name"`
}

type Status struct {
	CreatedAt                    string `json:"createdAt"`
	Emoji                        string `json:"emoji"`
	Id                           string `json:"id"`
	IndicatesLimitedAvailability bool   `json:"indicatesLimitedAvailability"`
	Message                      string `json:"message"`
	Organization                 Org    `json:"organization"`
	UpdatedAt                    string `json:"updatedAt"`
}

type RetrievalQueryResponseData struct {
	Viewer Viewer `json:"viewer,omitempty"`
}

type Viewer struct {
	Status Status `json:"status,omitempty"`
}

// Queries & Mutations

var retrievalQuery = `
query getStatus {
  viewer {
    status {
      createdAt
      message
      emoji
      indicatesLimitedAvailability
      organization {
        name
      }
    }
  }
}
`

var updateMutation = `
mutation ($newStatus: ChangeUserStatusInput!) {
  changeUserStatus(input: $newStatus) {
    status {
      id
      message
      emoji
    }
  }
}
`
