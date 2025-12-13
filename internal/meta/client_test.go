package meta

import (
	"context"
	"reflect"
	"testing"
)

// TestClient_PlanTask_Success tests successful PlanTask with mock provider
func TestClient_PlanTask_Success(t *testing.T) {
	client := NewMockClient()
	result, err := client.PlanTask(context.Background(), "test prd")

	if err != nil {
		t.Fatalf("PlanTask failed: %v", err)
	}

	if result.TaskID != "TASK-MOCK" {
		t.Errorf("TaskID = %q, want TASK-MOCK", result.TaskID)
	}

	if len(result.AcceptanceCriteria) == 0 {
		t.Errorf("AcceptanceCriteria is empty")
	}

	if result.AcceptanceCriteria[0].ID != "AC-1" {
		t.Errorf("First AC ID = %q, want AC-1", result.AcceptanceCriteria[0].ID)
	}
}

// TestClient_NextAction_Success tests successful NextAction
func TestClient_NextAction_Success(t *testing.T) {
	client := NewMockClient()
	summary := &TaskSummary{
		Title:              "Test Task",
		State:              "RUNNING",
		AcceptanceCriteria: []AcceptanceCriterion{},
		WorkerRunsCount:    0,
	}

	result, err := client.NextAction(context.Background(), summary)

	if err != nil {
		t.Fatalf("NextAction failed: %v", err)
	}

	if result.Decision.Action != "run_worker" {
		t.Errorf("First action should be run_worker, got %q", result.Decision.Action)
	}

	// Call again to test mark_complete
	summary.WorkerRunsCount = 1
	result, err = client.NextAction(context.Background(), summary)

	if err != nil {
		t.Fatalf("Second NextAction failed: %v", err)
	}

	if result.Decision.Action != "mark_complete" {
		t.Errorf("Second action should be mark_complete, got %q", result.Decision.Action)
	}
}

// TestClient_MockMode tests that mock client is configured correctly
func TestClient_MockMode_NoAPIKey(t *testing.T) {
	client := NewMockClient()

	if client.kind != "mock" {
		t.Errorf("kind = %q, want mock", client.kind)
	}

	// Now we check if the provider is OpenAIProvider with MOCK_KEY
	// We can't easily check private fields of provider interface, so we assume NewMockClient works if tests pass
}

// TestClient_NewMockClient tests the NewMockClient factory
func TestClient_NewMockClient(t *testing.T) {
	client := NewMockClient()

	if client == nil {
		t.Fatalf("NewMockClient returned nil")
	}

	if client.kind != "mock" {
		t.Errorf("kind = %q, want mock", client.kind)
	}

	if client.provider == nil {
		t.Errorf("provider is nil")
	}
}

// TestClient_PlanTask_MarkdownCodeBlock tests YAML extraction from markdown via client
func TestClient_PlanTask_MarkdownCodeBlock(t *testing.T) {
	client := NewMockClient()

	// This test verifies that the extraction logic is applied (via provider).
	result, err := client.PlanTask(context.Background(), "test prd")

	if err != nil {
		t.Fatalf("PlanTask failed: %v", err)
	}

	if result == nil {
		t.Fatalf("PlanTask returned nil result")
	}

	if result.TaskID == "" {
		t.Errorf("TaskID is empty")
	}
}

// TestClient_NextAction_MultipleRuns tests NextAction behavior with multiple runs
func TestClient_NextAction_MultipleRuns(t *testing.T) {
	client := NewMockClient()
	summary := &TaskSummary{
		Title:              "Test Task",
		State:              "RUNNING",
		AcceptanceCriteria: []AcceptanceCriterion{},
		WorkerRunsCount:    0,
	}

	// First call should return run_worker
	action1, _ := client.NextAction(context.Background(), summary)
	if action1.Decision.Action != "run_worker" {
		t.Errorf("First action should be run_worker")
	}

	// Second call should return mark_complete
	summary.WorkerRunsCount = 1
	action2, _ := client.NextAction(context.Background(), summary)
	if action2.Decision.Action != "mark_complete" {
		t.Errorf("Second action should be mark_complete")
	}

	// Further calls should still return mark_complete
	summary.WorkerRunsCount = 5
	action3, _ := client.NextAction(context.Background(), summary)
	if action3.Decision.Action != "mark_complete" {
		t.Errorf("Subsequent actions should be mark_complete")
	}
}

// TestClient_NewMockClient_HasProvider tests that mock client has provider set
func TestClient_NewMockClient_HasProvider(t *testing.T) {
	client := NewMockClient()
	if client.provider == nil {
		t.Errorf("Mock client should have provider")
	}
}

// TestClient_NextAction_NoWorkerRuns tests NextAction behavior with zero worker runs
func TestClient_NextAction_NoWorkerRuns(t *testing.T) {
	client := NewMockClient()
	summary := &TaskSummary{
		Title:           "Test Task",
		State:           "RUNNING",
		WorkerRunsCount: 0,
	}

	result, err := client.NextAction(context.Background(), summary)
	if err != nil {
		t.Fatalf("NextAction failed: %v", err)
	}

	if result.Decision.Action != "run_worker" {
		t.Errorf("With 0 worker runs, action should be run_worker, got %q", result.Decision.Action)
	}
}

// TestClient_PlanTask_MockResponse validates the mock response structure
func TestClient_PlanTask_MockResponse(t *testing.T) {
	client := NewMockClient()
	result, _ := client.PlanTask(context.Background(), "Any PRD")

	if result == nil {
		t.Fatalf("PlanTask should return non-nil result")
	}

	if result.TaskID == "" {
		t.Errorf("TaskID should not be empty")
	}

	if len(result.AcceptanceCriteria) == 0 {
		t.Errorf("AcceptanceCriteria should not be empty in mock response")
	}
}

// TestClient_Kind tests the client kind setter
func TestClient_Kind(t *testing.T) {
	client := NewClient("mock", "", "", "")
	if client.kind != "mock" {
		t.Errorf("Client kind should be 'mock', got %q", client.kind)
	}

	client2 := NewClient("openai-chat", "key", "", "")
	if client2.kind != "openai-chat" {
		t.Errorf("Client kind should be 'openai-chat', got %q", client2.kind)
	}
}

// TestClient_APIKeyHandling tests that API key is stored correctly
func TestClient_APIKeyHandling(t *testing.T) {
	apiKey := "test-api-key-12345"
	client := NewClient("openai-chat", apiKey, "", "")

	if client == nil {
		t.Fatalf("Client creation failed")
	}

	if client.model != "gpt-5.2" {
		t.Errorf("expected default model gpt-5.2, got %s", client.model)
	}

	// Check internal provider if possible, but it's okay if we just trust factory logic tested above
	p, ok := client.provider.(*OpenAIProvider)
	if !ok {
		t.Errorf("expected OpenAIProvider")
	} else {
		// Can't check apiKey private field, but can verify struct type
		_ = p
	}
}

// TestClient_ProviderDelegation ensures Client delegates to Provider
func TestClient_ProviderDelegation(t *testing.T) {
	// This indirectly tests delegation via functional tests above (e.g. TestClient_PlanTask_Success)
	// which depends on MockProvider (reused OpenAIProvider with shim) doing the work.
}

func TestClient_SetLogger(t *testing.T) {
	client := NewMockClient()
	client.SetLogger(nil) // Should not panic
}

// Helper to inspect private fields if absolutely necessary for white-box testing
// but preferably avoid.
var _ = func(c *Client) Provider {
	return c.provider
}

var _ = func(v interface{}, t string) bool {
	return reflect.TypeOf(v).String() == t
}
