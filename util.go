package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		switch true {
		case os.IsNotExist(err):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		println("Could not retrieve home directory")
		os.Exit(1)
	}

	return homedir
}

func sendAPIRequest(token, query string, variables map[string]interface{}) (*GraphQLResponseBody, error) {
	// Create graphql-compatible body
	gqlBody := &GraphQLRequestBody{
		Query:         query,
		Variables:     variables,
		OperationName: nil,
	}

	// Serialize request body
	jsonBody, err := json.Marshal(gqlBody)
	if err != nil {
		return &GraphQLResponseBody{}, err
	}

	body := bytes.NewBuffer(jsonBody)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", body)
	if err != nil {
		return &GraphQLResponseBody{}, err
	}

	// Add request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return &GraphQLResponseBody{}, err
	}

	// Check response status code
	if res.StatusCode == http.StatusUnauthorized {
		return &GraphQLResponseBody{}, errors.New("received 'Unauthorized' response, please add a valid token")
	}

	if res.StatusCode != http.StatusOK {
		return &GraphQLResponseBody{}, errors.New(fmt.Sprintf("request failed with status %v: %v", res.StatusCode, res.Status))
	}

	// Read all incoming bytes of response body
	rawResponseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &GraphQLResponseBody{}, err
	}

	responseBody := &GraphQLResponseBody{}

	// Infuse response body into responseBody
	err = json.Unmarshal(rawResponseBody, responseBody)
	if err != nil {
		return &GraphQLResponseBody{}, err
	}

	return responseBody, nil
}

func validateTokenSet(c *Config) bool {
	return c.data.Token != ""
}
