package recaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var goolgleApisURL = func(projectID string, authKey string) string {
	return fmt.Sprintf(
		"https://recaptchaenterprise.googleapis.com/v1/projects/%s/assessments?key=%s",
		projectID,
		authKey)
}

type event struct {
	Token          string `json:"token"`
	SiteKey        string `json:"siteKey"`
	ExpectedAction string `json:"expectedAction"`
}

type riskAnalyhsis struct {
	Score   float64 `json:"score,omitempty"`
	Reasons string  `json:"reasosn,omitempty"`
}

type tokenProperties struct {
	Valid         bool      `json:"valid,omitempty"`
	InvalidReason string    `json:"invalidReason,omitempty"`
	Hostname      string    `json:"hostname,omitempty"`
	Action        string    `json:"action,omitempty"`
	CreateTime    time.Time `json:"createTime,omitempty"`
}

type assessmentData struct {
	event           `json:"event"`
	riskAnalyhsis   `json:"riskAnalysis,omitempty"`
	tokenProperties `json:"tokenProperties,omitempty"`
}

type AssessmentParams struct {
	ProjectID        string
	AuthKey          string
	RecapthcaSiteKey string
}

type Assessment struct {
	data assessmentData
}

func CreateAssessment(token string, action string, params AssessmentParams) (Assessment, error) {
	request, err := json.Marshal(
		assessmentData{
			event: event{
				Token:          token,
				SiteKey:        params.RecapthcaSiteKey,
				ExpectedAction: strings.ToLower(action),
			},
		})
	if err != nil {
		return Assessment{}, err
	}

	req, err := http.NewRequest(
		"POST",
		goolgleApisURL(params.ProjectID, params.AuthKey),
		bytes.NewBuffer(request),
	)
	if err != nil {
		return Assessment{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Assessment{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Assessment{}, fmt.Errorf("failed to retrieve assessment: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Assessment{}, err
	}

	result := assessmentData{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return Assessment{}, err
	}

	return Assessment{
		data: result,
	}, nil
}

func (a *Assessment) IsTokenValid() bool {
	return a.data.Valid
}

func (a *Assessment) GetInvalidReason() string {
	return a.data.InvalidReason
}

func (a *Assessment) IsActionValid() bool {
	return strings.ToLower(a.data.ExpectedAction) == strings.ToLower(a.data.Action)
}

func (a *Assessment) GetAction() string {
	return a.data.Action
}

func (a *Assessment) GetExpectedAction() string {
	return a.data.ExpectedAction
}

func (a *Assessment) IsScoreValid(threshold float64) bool {
	return threshold <= a.data.Score
}

func (a *Assessment) GetScore() float64 {
	return a.data.Score
}

func (a *Assessment) GetScoreReasons() string {
	return a.data.Reasons
}
