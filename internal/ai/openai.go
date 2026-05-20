package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

type openAIRequest struct {
	Model          string          `json:"model"`
	Messages       []openAIMessage `json:"messages"`
	Temperature    float64         `json:"temperature"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func openAIKey() (string, error) {
	for _, name := range []string{"OPENAI_API_KEY", "OpenAI_API_KEY"} {
		if k := os.Getenv(name); k != "" {
			return k, nil
		}
	}
	return "", fmt.Errorf("OPENAI_API_KEY is not set")
}

func openAIModel() string {
	if m := os.Getenv("OPENAI_MODEL"); m != "" {
		return m
	}
	return "gpt-4o-mini"
}

// CallLLM sends prompts to OpenAI Chat Completions and returns the assistant text.
func CallLLM(systemPrompt, userPrompt string) (string, error) {
	apiKey, err := openAIKey()
	if err != nil {
		return "", err
	}

	reqBody := openAIRequest{
		Model: openAIModel(),
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.2,
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{Type: "json_object"},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("openai marshal: %w", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}

	var raw []byte
	var statusCode int
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		req, err := http.NewRequest(http.MethodPost, openAIURL, bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("openai request: %w", err)
		}
		raw, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("openai read response: %w", err)
		}
		statusCode = resp.StatusCode
		if statusCode == http.StatusOK {
			break
		}
		if (statusCode == http.StatusTooManyRequests || statusCode >= 500) && attempt < 2 {
			continue
		}
		return "", fmt.Errorf("openai HTTP %d: %s", statusCode, string(raw))
	}

	var oaResp openAIResponse
	if err := json.Unmarshal(raw, &oaResp); err != nil {
		return "", fmt.Errorf("openai unmarshal: %w", err)
	}
	if oaResp.Error != nil {
		return "", fmt.Errorf("openai api error: %s", oaResp.Error.Message)
	}
	if len(oaResp.Choices) == 0 || oaResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("openai returned empty response")
	}

	return oaResp.Choices[0].Message.Content, nil
}
