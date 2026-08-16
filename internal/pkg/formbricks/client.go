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

type Environment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Survey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // "link" or "app"
	Status      string `json:"status"` // "draft" or "published"
	EnvironmentID string `json:"environmentId"`
}

type Question struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Headline map[string]string `json:"headline"`
	Required bool              `json:"required"`
}

type CreateEnvironmentRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateEnvironmentResponse struct {
	Data Environment `json:"data"`
}

type CreateSurveyRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	EnvironmentID string   `json:"environmentId"`
	Questions   []Question `json:"questions,omitempty"`
}

type UpdateSurveyRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Questions   *[]Question `json:"questions,omitempty"`
}

type CreateSurveyResponse struct {
	Data Survey `json:"data"`
}

func NewClient() *Client {
	return &Client{
		APIKey: os.Getenv("FORMBRICKS_API_KEY"),
	}
}

// CreateEnvironment creates a new environment in Formbricks
func (c *Client) CreateEnvironment(eventName string) (string, error) {
	reqBody := CreateEnvironmentRequest{
		Name: fmt.Sprintf("Event: %s", eventName),
		Type: "production",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/environments", BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create environment: status %d", resp.StatusCode)
	}

	var result CreateEnvironmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.ID, nil
}

// CreateSurveyWithDefaults creates a new survey with default registration questions
func (c *Client) CreateSurveyWithDefaults(name, envID string) (*Survey, error) {
	// Default questions for event registration
	questions := []Question{
		{
			ID:       "q1",
			Type:     "shortText",
			Headline: map[string]string{"default": "Full Name"},
			Required: true,
		},
		{
			ID:       "q2",
			Type:     "shortText",
			Headline: map[string]string{"default": "Email Address"},
			Required: true,
		},
		{
			ID:       "q3",
			Type:     "shortText",
			Headline: map[string]string{"default": "Phone Number"},
			Required: false,
		},
		{
			ID:       "q4",
			Type:     "multipleChoiceSingle",
			Headline: map[string]string{"default": "Are you a member of this club?"},
			Required: true,
		},
	}

	reqBody := CreateSurveyRequest{
		Name:        fmt.Sprintf("%s Registration", name),
		Description: fmt.Sprintf("Registration form for %s", name),
		Type:        "link",
		EnvironmentID: envID,
		Questions:   questions,
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

	var result CreateSurveyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
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

	var result CreateSurveyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
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

	var result CreateSurveyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
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

	var result CreateSurveyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
