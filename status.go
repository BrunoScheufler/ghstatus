package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/logrusorgru/aurora"
	"strings"
	"text/template"
	"time"
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
	if input.Emoji != "" && (!strings.HasPrefix(input.Emoji, ":") || !strings.HasSuffix(input.Emoji, ":")) {
		return nil, errors.New("invalid emoji format, please supply a valid emoji")
	}

	// Construct and send query
	variables := make(map[string]interface{})

	updateInput := UpdateStatusMutationInput{
		Message: input.Message,
		Emoji:   input.Emoji,
	}

	// Handle status expiry, add to variables
	if input.ExpiresAt != nil && *input.ExpiresAt != "" {
		duration, err := time.ParseDuration(*input.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("could not parse duration for expiresAt value: %w", err)
		}

		now := time.Now()

		updateInput.ExpiresAt = now.Add(duration).UTC().Format(time.RFC3339Nano)
	}

	// Handle organization, add to variables
	if input.Organization != nil {
		if *input.Organization != "" {
			// Look up organization by name (login)
			orgId, err := LookupOrganization(input.Config, *input.Organization)
			if err != nil {
				return nil, fmt.Errorf("could not lookup organization by name: %w", err)
			}
			updateInput.OrganizationId = orgId
		} else {
			updateInput.OrganizationId = *input.Organization
		}
	}

	// Handle limitedAvailability (busy), add to variables
	if input.LimitedAvailability != nil {
		updateInput.LimitedAvailability = *input.LimitedAvailability
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

func FormatStatus(status *UserStatus) (string, error) {
	if status.Message == "" {
		return "👉 Status not set.", nil
	}

	templateFuncs := template.FuncMap{
		"printEmoji": func(text string) string {
			if text == "" {
				return ""
			}
			return emoji.Sprint(text)
		},
		"isEmptyString": func(str string) bool {
			return str == ""
		},
		"formatBold": func(content interface{}) string {
			return aurora.Bold(content).String()
		},
	}

	tpl, err := template.New("status-template").Funcs(templateFuncs).Parse(`
{{ formatBold "Status" }}: {{ printEmoji .Emoji }}{{ .Message }}
🚫 Busy: {{ formatBold .IndicatesLimitedAvailability }}
⏱  {{ if isEmptyString .ExpiresAt }}Status does not expire. {{else}} Expires at {{ .ExpiresAt }} {{end}}
🏢 {{ if isEmptyString .Organization.Name }}Visible for everyone {{ else }} Visible for {{ .Organization.Name}} {{ end }}
`)
	if err != nil {
		return "", fmt.Errorf("could not create status template: %w", err)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, status)
	if err != nil {
		return "", fmt.Errorf("could not execute template: %w", err)
	}

	return buf.String(), nil
}
