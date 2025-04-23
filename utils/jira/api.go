package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/config"
)

// JiraIssue represents the data needed to create a JIRA ticket
type JiraIssue struct {
	Fields Fields `json:"fields"`
}

// Fields represents the fields of a JIRA ticket
type Fields struct {
	Project     Project   `json:"project"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	IssueType   IssueType `json:"issuetype"`
	Assignee    *Assignee `json:"assignee,omitempty"`
	Parent      *Parent   `json:"parent,omitempty"`
	Labels      []string  `json:"labels,omitempty"`
}

// Project represents the JIRA project
type Project struct {
	Key string `json:"key"`
}

// IssueType represents the type of a JIRA ticket
type IssueType struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Assignee represents the assignee of the ticket
type Assignee struct {
	ID string `json:"accountId"`
}

// Parent represents the parent ticket
type Parent struct {
	Key string `json:"key"`
}

// CreateJiraTicket creates a JIRA ticket via the API
func CreateJiraTicket(title, description, assignee, parentTicket string) (string, error) {
	cfg := config.Get()

	// Check that the necessary configurations are defined
	if cfg.App.JiraAPIToken == "" || cfg.App.JiraUsername == "" {
		return "", fmt.Errorf("JIRA API credentials not set")
	}

	// Use the projectKey from the configuration
	projectKey := cfg.App.JiraProjectKey
	if projectKey == "" {
		return "", fmt.Errorf("JIRA project key not set")
	}

	// Use the ticket type ID from the configuration or a default value
	issueType := cfg.App.JiraIssueType
	if issueType == "" {
		issueType = "10080" // Default ID for tasks (adjust according to your Jira instance)
	}

	// Prepare the JSON payload
	issue := JiraIssue{
		Fields: Fields{
			Project: Project{
				Key: projectKey,
			},
			Summary:     title,
			Description: description,
			IssueType: IssueType{
				Name: issueType,
			},
		},
	}

	if strings.ToLower(issueType) != "sous-tâche" &&
		strings.ToLower(issueType) != "sous-tache" &&
		strings.ToLower(issueType) != "sub-task" {
		issue.Fields.Labels = []string{"CI/CD", "automation", "git-report"}
	}

	if assignee != "" {
		accountID, err := FindAssignableUser(assignee)
		if err != nil {
			fmt.Printf("Warning: Unable to find user %s: %v\n", assignee, err)
		} else {
			issue.Fields.Assignee = &Assignee{
				ID: accountID,
			}
		}
	}

	// Add the parent ticket if specified
	if parentTicket != "" {
		issue.Fields.Parent = &Parent{
			Key: parentTicket,
		}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(issue)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Display the JSON for debugging
	fmt.Printf("JIRA Request JSON: %s\n", string(jsonData))

	// Create the request
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	apiURL := fmt.Sprintf("%s/rest/api/2/issue", cfg.App.JiraBaseURL)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(cfg.App.JiraUsername, cfg.App.JiraAPIToken)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body for debugging
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("JIRA API Response [%d]: %s\n", resp.StatusCode, string(respBody))

	// Process the response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("error response from JIRA: %s - %s", resp.Status, string(respBody))
	}

	// Extract the ID of the created ticket
	var result struct {
		ID   string `json:"id"`
		Key  string `json:"key"`
		Self string `json:"self"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Return the key of the created ticket
	ticketURL := fmt.Sprintf("%s/browse/%s", cfg.App.JiraBaseURL, result.Key)
	return ticketURL, nil
}

// FindAssignableUser searches for a user that can be assigned to a ticket in a project
func FindAssignableUser(displayName string) (string, error) {
	cfg := config.Get()

	// Build the API URL
	apiURL := fmt.Sprintf("%s/rest/api/2/user/assignable/search?project=%s&maxResults=100",
		cfg.App.JiraBaseURL,
		cfg.App.JiraProjectKey)

	// Create the request
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	// Add authentication
	req.SetBasicAuth(cfg.App.JiraUsername, cfg.App.JiraAPIToken)
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check the response status
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error while searching for users: %s", string(body))
	}

	// Parse the JSON
	var users []struct {
		AccountID   string `json:"accountId"`
		DisplayName string `json:"displayName"`
	}

	if err := json.Unmarshal(body, &users); err != nil {
		return "", err
	}

	// Search for the user by display name (case insensitive)
	lowerDisplayName := strings.ToLower(displayName)
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.DisplayName), lowerDisplayName) {
			return user.AccountID, nil
		}
	}

	// If the user is not found
	return "", fmt.Errorf("no assignable user found with the name: %s", displayName)
}
