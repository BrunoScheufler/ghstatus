package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		switch {
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

func LookupOrganization(Config *Config, Name string) (string, error) {
	// Construct and execute request
	variables := make(map[string]interface{})

	variables["input"] = Name

	rawResponse, err := SendApiRequest(Config.Data.Token, organizationLookupQuery, variables)
	if err != nil {
		return "", fmt.Errorf("could not send API request: %w", err)
	}

	responseData := OrganizationLookupQueryResponse{}
	err = json.Unmarshal(rawResponse, &responseData)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response data: %w", err)
	}

	return responseData.Organization.ID, nil
}
