package orchestrator

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/biwakonbu/agent-runner/internal/orchestrator/persistence"
)

type SchedulerV2 struct {
	repo     persistence.WorkspaceRepository
	executor ExecutorV2
	logger   *slog.Logger
}

func NewSchedulerV2(repo persistence.WorkspaceRepository, executor ExecutorV2, logger *slog.Logger) *SchedulerV2 {
	return &SchedulerV2{
		repo:     repo,
		executor: executor,
		logger:   logger,
	}
}

// CheckAndSchedule is the main entry point to be called periodically or on event.
func (s *SchedulerV2) CheckAndSchedule(ctx context.Context) error {
	// 1. Load pending tasks
	tasksState, err := s.repo.State().LoadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	nodesRuntime, err := s.repo.State().LoadNodesRuntime()
	if err != nil {
		return fmt.Errorf("failed to load nodes runtime: %w", err)
	}

	agentsState, err := s.repo.State().LoadAgents()
	if err != nil {
		return fmt.Errorf("failed to load agents: %w", err)
	}

	// 2. Filter schedule-able tasks
	var candidates []persistence.TaskState
	for _, t := range tasksState.Tasks {
		if t.Status == "pending" {
			// Check dependencies
			if s.allDependenciesSatisfied(t.NodeID, nodesRuntime) {
				candidates = append(candidates, t)
			}
		}
	}

	if len(candidates) == 0 {
		return nil // Nothing to schedule
	}

	// 3. Dispatch to Agents
	// Simple FIFO dispatch for MVP
	actionsToAppend := []persistence.Action{}
	modifiedTasks := false
	modifiedAgents := false

	for _, task := range candidates {
		agentID, ok := s.findAvailableAgent(agentsState, task)
		if !ok {
			continue // No agent available
		}

		// Dispatch!
		now := time.Now()

		// Create Action
		action := persistence.Action{
			ID:          newID("act"),
			At:          now,
			Kind:        "task.started",
			WorkspaceID: "TODO-ws-id", // Needs to be plumbed
			Payload: map[string]interface{}{
				"task_id":  task.TaskID,
				"agent_id": agentID,
			},
		}
		actionsToAppend = append(actionsToAppend, action)

		// Update Task State
		for i, t := range tasksState.Tasks {
			if t.TaskID == task.TaskID {
				tasksState.Tasks[i].Status = "running"
				tasksState.Tasks[i].AssignedAgent = agentID
				tasksState.Tasks[i].UpdatedAt = now
				modifiedTasks = true
				break
			}
		}

		// Update Agent State
		for i, a := range agentsState.Agents {
			if a.AgentID == agentID {
				agentsState.Agents[i].RunningTasks = append(agentsState.Agents[i].RunningTasks, task.TaskID)
				modifiedAgents = true
				break
			}
		}

		// Trigger Executor (Async)
		// Note: Executor needs to be updated to accept TaskStateV2
		go func(t persistence.TaskState) {
			s.logger.Info("Executing task", "task_id", t.TaskID)
			_ = s.executor.Execute(ctx, t)
		}(task)
	}

	// 4. Persist Changes
	if len(actionsToAppend) > 0 {
		for _, a := range actionsToAppend {
			if err := s.repo.History().AppendAction(&a); err != nil {
				return fmt.Errorf("failed to append action: %w", err)
			}
		}
	}
	if modifiedTasks {
		if err := s.repo.State().SaveTasks(tasksState); err != nil {
			return fmt.Errorf("failed to save tasks: %w", err)
		}
	}
	if modifiedAgents {
		if err := s.repo.State().SaveAgents(agentsState); err != nil {
			return fmt.Errorf("failed to save agents: %w", err)
		}
	}

	return nil
}

func (s *SchedulerV2) allDependenciesSatisfied(nodeID string, runtime *persistence.NodesRuntime) bool {
	nodeDesign, err := s.repo.Design().GetNode(nodeID)
	if err != nil {
		s.logger.Error("failed to get node design", "node_id", nodeID, "err", err)
		return false // Fail safe
	}

	for _, depID := range nodeDesign.Dependencies {
		satisfied := false
		for _, n := range runtime.Nodes {
			if n.NodeID == depID {
				// Dependency must be implemented or verified to proceed?
				// MVP: "implemented" is enough to start next task?
				if n.Status == "implemented" || n.Status == "verified" {
					satisfied = true
				}
				break
			}
		}
		if !satisfied {
			return false
		}
	}
	return true
}

func (s *SchedulerV2) findAvailableAgent(agents *persistence.AgentsState, _ persistence.TaskState) (string, bool) {
	for _, a := range agents.Agents {
		if len(a.RunningTasks) < a.MaxParallel {
			// Check capabilities (naive check)
			// Real logic: check if agent kind matches task requirements
			return a.AgentID, true
		}
	}
	return "", false
}

// Helper to generate ID (placeholders)
func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
