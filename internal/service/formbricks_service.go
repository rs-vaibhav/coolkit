package service

import (
"bytes"
"context"
"encoding/json"
"fmt"
"net/http"
"os"
)

type FormbricksService struct {
apiKey string
apiURL string
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
Name          string `json:"name"`
Description   string `json:"description"`
Type          string `json:"type"`
Status        string `json:"status"`
EnvironmentID string `json:"environmentId"`
}

func NewFormbricksService() *FormbricksService {
return &FormbricksService{
apiKey: os.Getenv("FORMBRICKS_API_KEY"),
apiURL: "https://app.formbricks.com/api/v1",
}
}

func (s *FormbricksService) CreateSurvey(ctx context.Context, envID, eventName string) (*FormbricksSurvey, error) {
reqBody := CreateSurveyRequest{
Name:          fmt.Sprintf("%s Registration", eventName),
Description:   fmt.Sprintf("Registration form for %s", eventName),
Type:          "link",
Status:        "production",
EnvironmentID: envID,
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
return nil, fmt.Errorf("failed to create survey: status %d", resp.StatusCode)
}

var survey FormbricksSurvey
if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
return nil, err
}

return &survey, nil
}

func (s *FormbricksService) GetSurvey(ctx context.Context, surveyID string) (*FormbricksSurvey, error) {
req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/surveys/%s", s.apiURL, surveyID), nil)
if err != nil {
return nil, err
}

req.Header.Set("x-api-key", s.apiKey)

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
return nil, err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
return nil, fmt.Errorf("failed to get survey: status %d", resp.StatusCode)
}

var survey FormbricksSurvey
if err := json.NewDecoder(resp.Body).Decode(&survey); err != nil {
return nil, err
}

return &survey, nil
}
