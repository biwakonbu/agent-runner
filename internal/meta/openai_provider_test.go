package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockRoundTripper allows controlling HTTP responses for testing
type mockRoundTripper struct {
	mu                sync.Mutex
	responses         []http.Response
	errors            []error
	callCount         int
	timeBetweenCalls  []time.Time
	returnedResponses int
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.timeBetweenCalls = append(m.timeBetweenCalls, time.Now())
	defer func() {
		m.callCount++
	}()

	if m.returnedResponses < len(m.errors) && m.errors[m.returnedResponses] != nil {
		m.returnedResponses++
		return nil, m.errors[m.returnedResponses-1]
	}

	if m.returnedResponses < len(m.responses) {
		resp := m.responses[m.returnedResponses]
		m.returnedResponses++
		return &resp, nil
	}

	// Default successful response
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"role":"assistant","content":"test response"}}]}`)),
		Header:     make(http.Header),
	}, nil
}

// mockRoundTripperFunc is a function-based RoundTripper for simple mock implementations
type mockRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f mockRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// timeoutError implements net.Error interface for testing
type timeoutError struct{}

func (t *timeoutError) Error() string   { return "timeout" }
func (t *timeoutError) Timeout() bool   { return true }
func (t *timeoutError) Temporary() bool { return false }

// TestIsRetryableError tests the isRetryableError helper function
func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		resp        *http.Response
		shouldRetry bool
	}{
		{
			name:        "No error, 200 response",
			err:         nil,
			resp:        &http.Response{StatusCode: 200},
			shouldRetry: false,
		},
		{
			name:        "No error, 4xx response",
			err:         nil,
			resp:        &http.Response{StatusCode: 400},
			shouldRetry: false,
		},
		{
			name:        "No error, 500 response",
			err:         nil,
			resp:        &http.Response{StatusCode: 500},
			shouldRetry: true,
		},
		{
			name:        "No error, 503 response",
			err:         nil,
			resp:        &http.Response{StatusCode: 503},
			shouldRetry: true,
		},
		{
			name:        "No error, 429 rate limit response",
			err:         nil,
			resp:        &http.Response{StatusCode: 429},
			shouldRetry: true,
		},
		{
			name:        "Timeout error",
			err:         &timeoutError{},
			resp:        nil,
			shouldRetry: true,
		},
		{
			name:        "Context canceled",
			err:         context.Canceled,
			resp:        nil,
			shouldRetry: false,
		},
		{
			name:        "Context deadline exceeded",
			err:         context.DeadlineExceeded,
			resp:        nil,
			shouldRetry: true,
		},
		{
			name:        "Generic error",
			err:         io.EOF,
			resp:        nil,
			shouldRetry: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isRetryableError(test.err, test.resp)
			if result != test.shouldRetry {
				t.Errorf("Expected shouldRetry=%v, got %v", test.shouldRetry, result)
			}
		})
	}
}

// TestCallLLM_SuccessFirstAttempt tests that successful first attempt returns immediately
func TestCallLLM_SuccessFirstAttempt(t *testing.T) {
	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: &mockRoundTripper{
				responses: []http.Response{
					{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"role":"assistant","content":"success"}}]}`)),
						Header:     make(http.Header),
					},
				},
			},
		},
	}

	ctx := context.Background()
	result, err := provider.callLLM(ctx, "system", "user")

	if err != nil {
		t.Fatalf("callLLM failed: %v", err)
	}

	if result != "success" {
		t.Errorf("Expected 'success', got %q", result)
	}
}

// TestCallLLM_RetryOn5xx tests retry behavior on 5xx errors
func TestCallLLM_RetryOn5xx(t *testing.T) {
	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: &mockRoundTripper{
				responses: []http.Response{
					{
						StatusCode: 503,
						Body:       io.NopCloser(bytes.NewBufferString("Service Unavailable")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"role":"assistant","content":"recovered"}}]}`)),
						Header:     make(http.Header),
					},
				},
			},
		},
	}

	ctx := context.Background()
	result, err := provider.callLLM(ctx, "system", "user")

	if err != nil {
		t.Fatalf("callLLM failed: %v", err)
	}

	if result != "recovered" {
		t.Errorf("Expected 'recovered', got %q", result)
	}
}

// TestCallLLM_RetryOnRateLimit tests retry behavior on 429 (rate limit) errors
func TestCallLLM_RetryOnRateLimit(t *testing.T) {
	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: &mockRoundTripper{
				responses: []http.Response{
					{
						StatusCode: 429,
						Body:       io.NopCloser(bytes.NewBufferString("Too Many Requests")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"role":"assistant","content":"rate limited"}}]}`)),
						Header:     make(http.Header),
					},
				},
			},
		},
	}

	ctx := context.Background()
	result, err := provider.callLLM(ctx, "system", "user")

	if err != nil {
		t.Fatalf("callLLM failed: %v", err)
	}

	if result != "rate limited" {
		t.Errorf("Expected 'rate limited', got %q", result)
	}
}

// TestCallLLM_NoRetryOn4xx tests that 4xx errors are not retried
func TestCallLLM_NoRetryOn4xx(t *testing.T) {
	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: &mockRoundTripper{
				responses: []http.Response{
					{
						StatusCode: 401,
						Status:     "401 Unauthorized",
						Body:       io.NopCloser(bytes.NewBufferString("Unauthorized")),
						Header:     make(http.Header),
					},
				},
			},
		},
	}

	ctx := context.Background()
	_, err := provider.callLLM(ctx, "system", "user")

	if err == nil {
		t.Fatal("callLLM should have failed with 401 error")
	}

	if !strings.Contains(err.Error(), "Unauthorized") {
		t.Errorf("Expected Unauthorized error, got: %v", err)
	}
}

// TestCallLLM_MaxRetriesExceeded tests failure after max retries
func TestCallLLM_MaxRetriesExceeded(t *testing.T) {
	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: &mockRoundTripper{
				responses: []http.Response{
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 1")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 2")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 3")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 4")),
						Header:     make(http.Header),
					},
				},
			},
		},
	}

	ctx := context.Background()
	_, err := provider.callLLM(ctx, "system", "user")

	if err == nil {
		t.Fatal("callLLM should have failed after max retries")
	}

	if !strings.Contains(err.Error(), "failed after 3 retries") {
		t.Errorf("Expected 'failed after 3 retries', got: %v", err)
	}
}

// TestCallLLM_RetryOnTimeout tests retry behavior on timeout errors
func TestCallLLM_RetryOnTimeout(t *testing.T) {
	callCount := 0

	customTransport := mockRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if callCount == 0 {
			callCount++
			return nil, &timeoutError{}
		}
		// Second attempt succeeds
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"role":"assistant","content":"timeout recovered"}}]}`)),
			Header:     make(http.Header),
		}, nil
	})

	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: customTransport,
		},
	}

	ctx := context.Background()
	result, err := provider.callLLM(ctx, "system", "user")

	if err != nil {
		t.Fatalf("callLLM failed: %v", err)
	}

	if result != "timeout recovered" {
		t.Errorf("Expected 'timeout recovered', got %q", result)
	}
}

// TestCallLLM_ExponentialBackoff tests that backoff delays are exponential
func TestCallLLM_ExponentialBackoff(t *testing.T) {
	startTime := time.Now()

	provider := &OpenAIProvider{
		apiKey: "test-api-key",
		model:  "gpt-4-turbo",
		client: &http.Client{
			Transport: &mockRoundTripper{
				responses: []http.Response{
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 1")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 2")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 3")),
						Header:     make(http.Header),
					},
					{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString("Error 4")),
						Header:     make(http.Header),
					},
				},
			},
		},
	}

	ctx := context.Background()
	_, _ = provider.callLLM(ctx, "system", "user")

	elapsed := time.Since(startTime)

	// Should have delays: 1s (after attempt 0) + 2s (after attempt 1) + 4s (after attempt 2) = 7 seconds minimum
	expectedMinimum := 7 * time.Second
	if elapsed < expectedMinimum {
		t.Errorf("Expected minimum elapsed time of %v, but got %v", expectedMinimum, elapsed)
	}

	// Allow some overhead (max 9 seconds to account for execution time)
	expectedMaximum := 9 * time.Second
	if elapsed > expectedMaximum {
		t.Errorf("Elapsed time %v exceeds expected maximum of %v", elapsed, expectedMaximum)
	}
}

// TestOpenAIProvider_SystemPromptUsage checks if system prompt is sent correctly
func TestOpenAIProvider_SystemPromptUsage(t *testing.T) {
	override := "You are a custom agent."

	type testChatRequest struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}

	var capturedRequest testChatRequest

	mockTransport := mockRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &capturedRequest); err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"role":"assistant","content":"ok"}}]}`)),
			Header:     make(http.Header),
		}, nil
	})

	provider := NewOpenAIProvider("key", "gpt-4-turbo", override)
	// We need to inject the mock transport. Since client is private, we recreate it here manually
	// similar to how we did in other tests.
	provider.client = &http.Client{Transport: mockTransport}

	ctx := context.Background()
	_, err := provider.callLLM(ctx, override, "user") // callLLM takes systemPrompt as arg, but PlanTask uses p.systemPrompt
	// Wait, callLLM takes systemPrompt. So checking if callLLM uses argument is trivial.
	// We want to check if PlanTask uses p.systemPrompt.

	if err != nil {
		t.Fatalf("callLLM failed: %v", err)
	}

	// Verify system prompt in request
	foundSystem := false
	for _, msg := range capturedRequest.Messages {
		if msg.Role == "system" {
			foundSystem = true
			if msg.Content != override {
				t.Errorf("Expected system prompt %q, got %q", override, msg.Content)
			}
		}
	}

	if !foundSystem {
		t.Error("System prompt not found in request")
	}

	// Verify PlanTask uses stored system prompt
	capturedRequest = testChatRequest{} // Request capture reset
	_, err = provider.PlanTask(ctx, "PRD")
	// Ignore parse error as response is "ok" not YAML
	_ = err

	foundSystem = false
	for _, msg := range capturedRequest.Messages {
		if msg.Role == "system" {
			foundSystem = true
			// The request should contain p.systemPrompt if we called PlanTask
			if msg.Content != override {
				t.Errorf("Expected system prompt %q in PlanTask, got %q", override, msg.Content)
			}
		}
	}
	if !foundSystem {
		t.Error("System prompt not found in PlanTask request")
	}
}
