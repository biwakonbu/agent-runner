// Package logging provides unified structured logging for the Multiverse system.
// It supports trace ID propagation, multiple log levels, and JSON/Text formatting.
package logging

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// TraceIDKey is the context key for trace ID.
	TraceIDKey contextKey = "trace_id"

	// AttrTraceID is the slog attribute name for trace ID.
	AttrTraceID = "trace_id"

	// AttrComponent is the slog attribute name for component.
	AttrComponent = "component"

	// AttrDuration is the slog attribute name for duration.
	AttrDuration = "duration_ms"
)

// Config holds logger configuration.
type Config struct {
	// Level is the minimum log level to output.
	Level slog.Level

	// JSONFormat enables JSON output format. If false, text format is used.
	JSONFormat bool

	// AddSource adds source file information to logs.
	AddSource bool
}

// DefaultConfig returns the default logger configuration.
func DefaultConfig() Config {
	return Config{
		Level:      slog.LevelInfo,
		JSONFormat: false,
		AddSource:  false,
	}
}

// ProductionConfig returns a production-ready logger configuration.
func ProductionConfig() Config {
	return Config{
		Level:      slog.LevelInfo,
		JSONFormat: true,
		AddSource:  true,
	}
}

// DebugConfig returns a debug-friendly logger configuration.
func DebugConfig() Config {
	return Config{
		Level:      slog.LevelDebug,
		JSONFormat: false,
		AddSource:  true,
	}
}

// ContextWithTraceID adds a trace ID to the context.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// TraceIDFromContext extracts the trace ID from context.
// Returns empty string if no trace ID is set.
func TraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// NewLogger creates a new structured logger with the given configuration.
func NewLogger(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

// WithTraceID returns a new logger with the trace ID from context as an attribute.
// If no trace ID is in context or logger is nil, returns the original logger (or default logger if nil).
func WithTraceID(logger *slog.Logger, ctx context.Context) *slog.Logger {
	if logger == nil {
		logger = slog.Default()
	}
	traceID := TraceIDFromContext(ctx)
	if traceID == "" {
		return logger
	}
	return logger.With(slog.String(AttrTraceID, traceID))
}

// WithComponent returns a new logger with the component name as an attribute.
// If logger is nil, uses the default logger.
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return logger.With(slog.String(AttrComponent, component))
}

// LogDuration logs a duration in milliseconds.
func LogDuration(start time.Time) slog.Attr {
	return slog.Float64(AttrDuration, float64(time.Since(start).Milliseconds()))
}

// LogRequest logs request/response metadata.
// Used for logging API calls (e.g., LLM requests) at INFO level.
type LogRequest struct {
	Method       string
	URL          string
	StatusCode   int
	DurationMs   float64
	RequestSize  int
	ResponseSize int
	Error        string
}

// ToAttrs converts LogRequest to slog attributes.
func (r LogRequest) ToAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.String("method", r.Method),
		slog.String("url", r.URL),
		slog.Int("status_code", r.StatusCode),
		slog.Float64("duration_ms", r.DurationMs),
		slog.Int("request_size", r.RequestSize),
		slog.Int("response_size", r.ResponseSize),
	}
	if r.Error != "" {
		attrs = append(attrs, slog.String("error", r.Error))
	}
	return attrs
}

// TaskLogContext provides context for task-related logging.
type TaskLogContext struct {
	TaskID    string
	Title     string
	State     string
	LoopCount int
}

// ToAttrs converts TaskLogContext to slog attributes.
func (t TaskLogContext) ToAttrs() []slog.Attr {
	return []slog.Attr{
		slog.String("task_id", t.TaskID),
		slog.String("task_title", t.Title),
		slog.String("task_state", t.State),
		slog.Int("loop_count", t.LoopCount),
	}
}

// WorkerLogContext provides context for worker-related logging.
type WorkerLogContext struct {
	ContainerID string
	Image       string
	Command     string
	ExitCode    int
	DurationMs  float64
}

// ToAttrs converts WorkerLogContext to slog attributes.
func (w WorkerLogContext) ToAttrs() []slog.Attr {
	return []slog.Attr{
		slog.String("container_id", w.ContainerID),
		slog.String("image", w.Image),
		slog.String("command", w.Command),
		slog.Int("exit_code", w.ExitCode),
		slog.Float64("duration_ms", w.DurationMs),
	}
}
