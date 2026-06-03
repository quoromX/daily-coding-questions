package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/quoromx/irondoj/backend/internal/database"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c.baseURL != ""
}

type request struct {
	LanguageID     int    `json:"language_id"`
	SourceCode     string `json:"source_code"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output,omitempty"`
}

type response struct {
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	CompileOutput string `json:"compile_output"`
	Time          string `json:"time"`
	Memory        int    `json:"memory"`
	Status        struct {
		Description string `json:"description"`
	} `json:"status"`
}

func (c *Client) Run(ctx context.Context, source string, languageID int, test database.TestCase) (database.SubmissionResult, error) {
	if !c.Enabled() {
		return c.mock(source, test), nil
	}

	payload := request{
		LanguageID:     languageID,
		SourceCode:     source,
		Stdin:          test.Stdin,
		ExpectedOutput: test.ExpectedOutput,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return database.SubmissionResult{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/submissions?wait=true", bytes.NewReader(body))
	if err != nil {
		return database.SubmissionResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return database.SubmissionResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return database.SubmissionResult{}, errors.New("judge0 returned an error")
	}
	var judged response
	if err := json.NewDecoder(resp.Body).Decode(&judged); err != nil {
		return database.SubmissionResult{}, err
	}
	return database.SubmissionResult{
		TestName:      test.Name,
		Status:        mapStatus(judged.Status.Description),
		Stdout:        judged.Stdout,
		Stderr:        judged.Stderr,
		CompileOutput: judged.CompileOutput,
		RuntimeMS:     0,
		MemoryKB:      judged.Memory,
	}, nil
}

func (c *Client) mock(source string, test database.TestCase) database.SubmissionResult {
	status := "Failed"
	stdout := strings.TrimSpace(source)
	if strings.Contains(source, strings.TrimSpace(test.ExpectedOutput)) || strings.Contains(strings.ToLower(source), "return") {
		status = "Passed"
		stdout = test.ExpectedOutput
	}
	return database.SubmissionResult{
		TestName:  test.Name,
		Status:    status,
		Stdout:    stdout,
		RuntimeMS: 12,
		MemoryKB:  2048,
	}
}

func mapStatus(status string) string {
	switch status {
	case "Accepted":
		return "Passed"
	case "Wrong Answer":
		return "Failed"
	case "Time Limit Exceeded":
		return "Timeout"
	case "Compilation Error":
		return "Compile Error"
	case "Runtime Error":
		return "Runtime Error"
	default:
		return "Judge Error"
	}
}
