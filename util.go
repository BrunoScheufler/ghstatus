package main

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
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

func sendAPIRequest(token, query string, variables map[string]interface{}) (interface{}, error) {
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
