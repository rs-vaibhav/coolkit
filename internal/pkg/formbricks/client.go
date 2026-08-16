package formbricks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const BaseURL = "https://app.formbricks.com/api/v1"

type Client struct {
	APIKey string
}

type Survey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // "link" or "app"
	Status      string `json:"status"` // "draft" or "published"
	EnvironmentID string `json:"environmentId"`
}

type CreateSurveyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	EnvironmentID string `json:"environmentId"`
}

type UpdateSurveyRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func NewClient() *Client {
	return &Client{
		APIKey: os.Getenv("FORMBRICKS_API_KEY"),
	}
}

// CreateSurvey creates a new survey in Formbricks
func (c *Client) CreateSurvey(name, description, envID string) (*Survey, error) {
	reqBody := CreateSurveyRequest{
		Name:        name,
		Description: description,
		Type:        "link",
		EnvironmentID: envID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/surveys", BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create survey: status %d", resp.StatusCode)
	}

	var survey Survey
	if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
		return nil, err
	}

	return &survey, nil
}

// GetSurvey retrieves a survey by ID
func (c *Client) GetSurvey(surveyID string) (*Survey, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/surveys/%s", BaseURL, surveyID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get survey: status %d", resp.StatusCode)
	}

	var survey Survey
	if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
		return nil, err
	}

	return &survey, nil
}

// UpdateSurvey updates an existing survey
func (c *Client) UpdateSurvey(surveyID string, updates UpdateSurveyRequest) (*Survey, error) {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/surveys/%s", BaseURL, surveyID), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update survey: status %d", resp.StatusCode)
	}

	var survey Survey
	if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
		return nil, err
	}

	return &survey, nil
}
