package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestContextWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-id-123"

	// Initially no trace ID
	if got := TraceIDFromContext(ctx); got != "" {
		t.Errorf("TraceIDFromContext(empty ctx) = %q, want empty string", got)
	}

	// Add trace ID
	ctx = ContextWithTraceID(ctx, traceID)
	if got := TraceIDFromContext(ctx); got != traceID {
		t.Errorf("TraceIDFromContext(ctx with trace) = %q, want %q", got, traceID)
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "default config",
			config: DefaultConfig(),
		},
		{
			name:   "production config",
			config: ProductionConfig(),
		},
		{
			name:   "debug config",
			config: DebugConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.config)
			if logger == nil {
				t.Error("NewLogger returned nil")
			}
		})
	}
}

func TestWithTraceID(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	traceID := "trace-abc-123"
	ctx := ContextWithTraceID(context.Background(), traceID)

	// Create logger with trace ID
	loggerWithTrace := WithTraceID(logger, ctx)
	loggerWithTrace.Info("test message")

	// Check that trace ID is in the log output
	output := buf.String()
	if !strings.Contains(output, traceID) {
		t.Errorf("Log output should contain trace ID %q, got: %s", traceID, output)
	}
}

func TestWithTraceID_NoTraceID(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	// Context without trace ID
	ctx := context.Background()
	loggerWithTrace := WithTraceID(logger, ctx)

	// Should return the same logger (no panic, no error)
	loggerWithTrace.Info("test message")

	output := buf.String()
	if strings.Contains(output, "trace_id") {
		t.Errorf("Log output should not contain trace_id when not set, got: %s", output)
	}
}

func TestWithComponent(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	component := "meta-client"
	loggerWithComp := WithComponent(logger, component)
	loggerWithComp.Info("test message")

	output := buf.String()
	if !strings.Contains(output, component) {
		t.Errorf("Log output should contain component %q, got: %s", component, output)
	}
}

func TestLogDuration(t *testing.T) {
	start := time.Now().Add(-100 * time.Millisecond) // 100ms ago
	attr := LogDuration(start)

	if attr.Key != AttrDuration {
		t.Errorf("LogDuration key = %q, want %q", attr.Key, AttrDuration)
	}

	durationMs := attr.Value.Float64()
	if durationMs < 100 {
		t.Errorf("LogDuration value = %f, want >= 100", durationMs)
	}
}

func TestLogRequest_ToAttrs(t *testing.T) {
	req := LogRequest{
		Method:       "POST",
		URL:          "https://api.openai.com/v1/chat/completions",
		StatusCode:   200,
		DurationMs:   1234.5,
		RequestSize:  500,
		ResponseSize: 2000,
	}

	attrs := req.ToAttrs()
	if len(attrs) != 6 {
		t.Errorf("ToAttrs() returned %d attributes, want 6", len(attrs))
	}

	// Test with error
	reqWithError := LogRequest{
		Method:       "POST",
		URL:          "https://api.openai.com/v1/chat/completions",
		StatusCode:   500,
		DurationMs:   100,
		RequestSize:  500,
		ResponseSize: 0,
		Error:        "internal server error",
	}

	attrsWithError := reqWithError.ToAttrs()
	if len(attrsWithError) != 7 {
		t.Errorf("ToAttrs() with error returned %d attributes, want 7", len(attrsWithError))
	}
}

func TestTaskLogContext_ToAttrs(t *testing.T) {
	taskCtx := TaskLogContext{
		TaskID:    "task-123",
		Title:     "Test Task",
		State:     "RUNNING",
		LoopCount: 2,
	}

	attrs := taskCtx.ToAttrs()
	if len(attrs) != 4 {
		t.Errorf("ToAttrs() returned %d attributes, want 4", len(attrs))
	}
}

func TestWorkerLogContext_ToAttrs(t *testing.T) {
	workerCtx := WorkerLogContext{
		ContainerID: "abc123",
		Image:       "agent-runner-codex:latest",
		Command:     "codex",
		ExitCode:    0,
		DurationMs:  5000,
	}

	attrs := workerCtx.ToAttrs()
	if len(attrs) != 5 {
		t.Errorf("ToAttrs() returned %d attributes, want 5", len(attrs))
	}
}

func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer

	// Create a custom logger that writes to our buffer
	cfg := ProductionConfig()
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	})
	logger := slog.New(handler)

	logger.Info("test message", slog.String("key", "value"))

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Log output is not valid JSON: %v\nOutput: %s", err, buf.String())
	}

	// Verify expected fields
	if _, ok := result["msg"]; !ok {
		t.Error("JSON log should have 'msg' field")
	}
	if _, ok := result["level"]; !ok {
		t.Error("JSON log should have 'level' field")
	}
}
