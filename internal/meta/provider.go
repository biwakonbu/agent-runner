package meta

import "context"

// Provider は LLM プロバイダを抽象化するインターフェース
// HTTP ベース (OpenAI) や CLI ベース (Codex, Claude) の実装を統一的に扱う
type Provider interface {
	// Name returns the provider name (e.g. "openai-chat", "codex-cli")
	Name() string

	// TestConnection verifies if the provider is reachable and configured correctly
	TestConnection(ctx context.Context) error

	// Decompose breaks down a user request into tasks
	Decompose(ctx context.Context, req *DecomposeRequest) (*DecomposeResponse, error)

	// PlanPatch generates operations to modify the project plan
	PlanPatch(ctx context.Context, req *PlanPatchRequest) (*PlanPatchResponse, error)

	// PlanTask generates acceptance criteria from a PRD
	PlanTask(ctx context.Context, prdText string) (*PlanTaskResponse, error)

	// NextAction decides the next step in task execution
	NextAction(ctx context.Context, taskSummary *TaskSummary) (*NextActionResponse, error)

	// CompletionAssessment evaluates if a task matches its criteria
	CompletionAssessment(ctx context.Context, taskSummary *TaskSummary) (*CompletionAssessmentResponse, error)
}
