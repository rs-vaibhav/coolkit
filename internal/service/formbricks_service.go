package service

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"os"
)

type FormbricksService struct {
apiKey   string
apiURL   string
apiURLV2 string
}

type FormbricksSurvey struct {
ID          string `json:"id"`
Name        string `json:"name"`
Description string `json:"description"`
Status      string `json:"status"`
Type        string `json:"type"`
Link        string `json:"link"`
}

type CreateSurveyRequest struct {
Name          string           `json:"name"`
Description   string           `json:"description"`
Type          string           `json:"type"`
Status        string           `json:"status"`
EnvironmentID string           `json:"environmentId"`
Questions     []SurveyQuestion `json:"questions,omitempty"`
}

type SurveyQuestion struct {
ID       string            `json:"id"`
Type     string            `json:"type"`
Headline map[string]string `json:"headline"`
Required bool              `json:"required"`
}

type UserInfoResponse struct {
Data UserInfo `json:"data"`
}

type UserInfo struct {
OrganizationID string `json:"organizationId"`
}

type CreateEnvironmentRequest struct {
Name string `json:"name"`
Type string `json:"type"`
}

type CreateEnvironmentResponse struct {
Data Environment `json:"data"`
}

type Environment struct {
ID   string `json:"id"`
Name string `json:"name"`
Type string `json:"type"`
}

func NewFormbricksService() *FormbricksService {
return &FormbricksService{
apiKey:   os.Getenv("FORMBRICKS_API_KEY"),
apiURL:   "https://app.formbricks.com/api/v1",
apiURLV2: "https://app.formbricks.com/api/v2",
}
}

// GetOrganizationID fetches the organization ID from Formbricks API
func (s *FormbricksService) GetOrganizationID(ctx context.Context) (string, error) {
req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/me", s.apiURLV2), nil)
if err != nil {
return "", err
}

req.Header.Set("x-api-key", s.apiKey)

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
body, _ := io.ReadAll(resp.Body)
return "", fmt.Errorf("failed to get user info: %s (Status: %d)", string(body), resp.StatusCode)
}

var result UserInfoResponse
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
return "", err
}

if result.Data.OrganizationID == "" {
return "", fmt.Errorf("no organization found for this API key")
}

return result.Data.OrganizationID, nil
}

// CreateEnvironment creates a new environment in Formbricks
func (s *FormbricksService) CreateEnvironment(ctx context.Context, orgID, eventName string) (string, error) {
reqBody := CreateEnvironmentRequest{
Name: fmt.Sprintf("Event: %s", eventName),
Type: "production",
}

jsonData, err := json.Marshal(reqBody)
if err != nil {
return "", err
}

req, err := http.NewRequestWithContext(ctx, "POST",
fmt.Sprintf("%s/organizations/%s/environments", s.apiURL, orgID),
bytes.NewBuffer(jsonData))
if err != nil {
return "", err
}

req.Header.Set("Content-Type", "application/json")
req.Header.Set("x-api-key", s.apiKey)

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
body, _ := io.ReadAll(resp.Body)
return "", fmt.Errorf("failed to create environment: %s (Status: %d)", string(body), resp.StatusCode)
}

var result CreateEnvironmentResponse
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
return "", err
}

return result.Data.ID, nil
}

func (s *FormbricksService) CreateSurvey(ctx context.Context, envID, eventName string) (*FormbricksSurvey, error) {
// Default questions for event registration
questions := []SurveyQuestion{
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
Name:          fmt.Sprintf("%s Registration", eventName),
Description:   fmt.Sprintf("Registration form for %s", eventName),
Type:          "link",
Status:        "production",
EnvironmentID: envID,
Questions:     questions,
}

jsonData, err := json.Marshal(reqBody)
if err != nil {
return nil, err
}

req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/surveys", s.apiURL), bytes.NewBuffer(jsonData))
if err != nil {
return nil, err
}

req.Header.Set("Content-Type", "application/json")
req.Header.Set("x-api-key", s.apiKey)

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return nil, err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
body, _ := io.ReadAll(resp.Body)
return nil, fmt.Errorf("failed to create survey: %s (Status: %d)", string(body), resp.StatusCode)
}

var survey FormbricksSurvey
if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
return nil, err
}

return &survey, nil
}
