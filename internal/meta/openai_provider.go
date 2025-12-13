package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/biwakonbu/agent-runner/internal/agenttools"
	"github.com/biwakonbu/agent-runner/internal/logging"
	"gopkg.in/yaml.v3"
)

// OpenAIProvider uses compatible HTTP API.
type OpenAIProvider struct {
	apiKey       string
	model        string
	systemPrompt string
	client       *http.Client
	logger       *slog.Logger
}

// NewOpenAIProvider creates a new OpenAIProvider
func NewOpenAIProvider(apiKey, model, systemPrompt string) *OpenAIProvider {
	if model == "" {
		model = agenttools.DefaultMetaModel
	}
	return &OpenAIProvider{
		apiKey:       apiKey,
		model:        model,
		systemPrompt: systemPrompt,
		client:       &http.Client{Timeout: 60 * time.Second},
		logger:       logging.WithComponent(slog.Default(), "meta-openai"),
	}
}

// SetLogger sets a custom logger
func (p *OpenAIProvider) SetLogger(logger *slog.Logger) {
	p.logger = logging.WithComponent(logger, "meta-openai")
}

func (p *OpenAIProvider) Name() string {
	return "openai-chat"
}

// TestConnection verifies the API connection
func (p *OpenAIProvider) TestConnection(ctx context.Context) error {
	// API key is now optional (environment might use proxy without auth)
	// But if it's missing, we log warning (not error) or just proceed.
	// API key チェックは削除（プロキシ環境では不要なため）
	// 実際の OpenAI API では認証エラーになる
	_ = strings.TrimSpace(p.apiKey) // lint: 使用済みマーク

	// Simple check using Decompose with dummy input
	_, err := p.Decompose(ctx, &DecomposeRequest{
		UserInput: "Ping",
		Context:   DecomposeContext{},
	})
	if err != nil {
		return fmt.Errorf("connection check failed: %w", err)
	}

	return nil
}

// Implementation specific structs
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func isRetryableError(err error, resp *http.Response) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true
		}
		if err == context.Canceled {
			return false
		}
		if err == context.DeadlineExceeded {
			return true
		}
		return true
	}
	if resp != nil {
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			return true
		}
		if resp.StatusCode == 429 {
			return true
		}
	}
	return false
}

func (p *OpenAIProvider) callLLM(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	const maxRetries = 3
	const baseDelay = 1 * time.Second

	logger := logging.WithTraceID(p.logger, ctx)
	start := time.Now()

	// REMOVED: explicit empty check for apiKey.

	reqBody := chatRequest{
		Model: p.model,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	logger.Info("calling LLM",
		slog.String("model", p.model),
		slog.Int("request_size", len(jsonBody)),
	)
	logger.Debug("LLM request",
		slog.String("system_prompt", systemPrompt),
		slog.String("user_prompt", userPrompt),
	)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		if p.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+p.apiKey)
		}

		resp, err := p.client.Do(req)
		if err != nil {
			lastErr = err
			if !isRetryableError(err, nil) {
				return "", err
			}
			if attempt < maxRetries {
				delay := baseDelay * time.Duration(1<<uint(attempt))
				logger.Warn("LLM request failed, retrying",
					slog.Int("attempt", attempt+1),
					slog.Any("error", err),
				)
				select {
				case <-time.After(delay):
				case <-ctx.Done():
					return "", ctx.Err()
				}
			}
			continue
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("OpenAI API error: %s %s", resp.Status, string(body))

			if !isRetryableError(nil, resp) {
				return "", lastErr
			}
			if attempt < maxRetries {
				delay := baseDelay * time.Duration(1<<uint(attempt))
				logger.Warn("LLM request failed with retryable status, retrying",
					slog.Int("attempt", attempt+1),
					slog.Int("status_code", resp.StatusCode),
				)
				select {
				case <-time.After(delay):
				case <-ctx.Done():
					return "", ctx.Err()
				}
			}
			continue
		}

		var result chatResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}

		if len(result.Choices) == 0 {
			return "", fmt.Errorf("no choices returned from LLM")
		}

		responseContent := result.Choices[0].Message.Content
		logger.Info("LLM call completed",
			slog.Int("response_size", len(responseContent)),
			logging.LogDuration(start),
		)
		return responseContent, nil
	}

	if lastErr != nil {
		return "", fmt.Errorf("LLM request failed after %d retries: %w", maxRetries, lastErr)
	}
	return "", fmt.Errorf("LLM request failed after %d retries", maxRetries)
}

func (p *OpenAIProvider) Decompose(ctx context.Context, req *DecomposeRequest) (*DecomposeResponse, error) {
	logger := logging.WithTraceID(p.logger, ctx)

	systemPrompt := decomposeSystemPrompt
	userPrompt := buildDecomposeUserPrompt(req)

	resp, err := p.callLLM(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	jsonStr := extractJSON(resp)
	yamlStr, err := jsonToYAML(jsonStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JSON to YAML: %w", err)
	}

	var msg MetaMessage
	if err := yaml.Unmarshal([]byte(yamlStr), &msg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	payloadBytes, err := yaml.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var decompose DecomposeResponse
	if err := yaml.Unmarshal(payloadBytes, &decompose); err != nil {
		return nil, fmt.Errorf("failed to parse decompose response: %w", err)
	}

	logger.Info("decompose completed", slog.Int("phases", len(decompose.Phases)))
	return &decompose, nil
}

func (p *OpenAIProvider) PlanPatch(ctx context.Context, req *PlanPatchRequest) (*PlanPatchResponse, error) {
	systemPrompt := planPatchSystemPrompt
	userPrompt := buildPlanPatchUserPrompt(req)

	resp, err := p.callLLM(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	jsonStr := extractJSON(resp)
	yamlStr, err := jsonToYAML(jsonStr)
	if err != nil {
		return nil, err
	}

	var msg MetaMessage
	if err := yaml.Unmarshal([]byte(yamlStr), &msg); err != nil {
		return nil, err
	}

	payloadBytes, err := yaml.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}
	var patch PlanPatchResponse
	if err := yaml.Unmarshal(payloadBytes, &patch); err != nil {
		return nil, err
	}

	return &patch, nil
}

func (p *OpenAIProvider) PlanTask(ctx context.Context, prdText string) (*PlanTaskResponse, error) {
	systemPrompt := p.systemPrompt
	if systemPrompt == "" {
		systemPrompt = `You are a Meta-agent that plans software development tasks.
Output MUST be a JSON block with the following structure:
{
  "type": "plan_task",
  "version": 1,
  "payload": { ... }
}`
	}
	userPrompt := fmt.Sprintf("PRD:\n%s\n\nGenerate the plan.", prdText)

	resp, err := p.callLLM(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	resp = extractYAML(resp)
	var msg MetaMessage
	if err := yaml.Unmarshal([]byte(resp), &msg); err != nil {
		if jsonErr := json.Unmarshal([]byte(resp), &msg); jsonErr != nil {
			return nil, fmt.Errorf("failed to parse response as YAML or JSON: %w", err)
		}
	}

	payloadBytes, err := yaml.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}
	var plan PlanTaskResponse
	if err := yaml.Unmarshal(payloadBytes, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (p *OpenAIProvider) NextAction(ctx context.Context, taskSummary *TaskSummary) (*NextActionResponse, error) {
	systemPrompt := p.systemPrompt
	if systemPrompt == "" {
		systemPrompt = `You are a Meta-agent that orchestrates a coding task.
Output MUST be a JSON block.`
	}
	// QH-005: Include WorkerRunsCount for mock detection
	contextSummary := fmt.Sprintf("Task: %s\nState: %s\nACs: %v\nWorkerRuns: %d",
		taskSummary.Title, taskSummary.State, len(taskSummary.AcceptanceCriteria), taskSummary.WorkerRunsCount)
	userPrompt := fmt.Sprintf("Context:\n%s\n\nDecide next action.", contextSummary)

	resp, err := p.callLLM(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	resp = extractYAML(resp)
	var msg MetaMessage
	if err := yaml.Unmarshal([]byte(resp), &msg); err != nil {
		return nil, err
	}

	payloadBytes, err := yaml.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}
	var action NextActionResponse
	if err := yaml.Unmarshal(payloadBytes, &action); err != nil {
		return nil, err
	}
	return &action, nil
}

func (p *OpenAIProvider) CompletionAssessment(ctx context.Context, taskSummary *TaskSummary) (*CompletionAssessmentResponse, error) {
	systemPrompt := p.systemPrompt
	if systemPrompt == "" {
		systemPrompt = `You are a Meta-agent evaluating task completion.`
	}
	userPrompt := fmt.Sprintf("Task: %s\nEvaluate completion.", taskSummary.Title)

	resp, err := p.callLLM(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	resp = extractYAML(resp)
	var msg MetaMessage
	if err := yaml.Unmarshal([]byte(resp), &msg); err != nil {
		return nil, err
	}

	payloadBytes, err := yaml.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}
	var assessment CompletionAssessmentResponse
	if err := yaml.Unmarshal(payloadBytes, &assessment); err != nil {
		return nil, err
	}
	return &assessment, nil
}
