package ide

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Workspace represents an IDE workspace configuration.
type Workspace struct {
	Version      string    `json:"version"`
	ProjectRoot  string    `json:"projectRoot"`
	DisplayName  string    `json:"displayName"`
	CreatedAt    time.Time `json:"createdAt"`
	LastOpenedAt time.Time `json:"lastOpenedAt"`
}

// WorkspaceStore handles workspace persistence.
type WorkspaceStore struct {
	BaseDir string
}

// NewWorkspaceStore creates a new WorkspaceStore.
// baseDir should be typically $HOME/.multiverse/workspaces
func NewWorkspaceStore(baseDir string) *WorkspaceStore {
	return &WorkspaceStore{BaseDir: baseDir}
}

// GetWorkspaceID generates a deterministic ID for a project root.
// ID = sha1(projectRoot)[:12]
func (s *WorkspaceStore) GetWorkspaceID(projectRoot string) string {
	h := sha1.New()
	h.Write([]byte(projectRoot))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)[:12]
}

// GetWorkspaceDir returns the directory path for a given workspace ID.
func (s *WorkspaceStore) GetWorkspaceDir(id string) string {
	return filepath.Join(s.BaseDir, id)
}

// LoadWorkspace loads a workspace by ID.
func (s *WorkspaceStore) LoadWorkspace(id string) (*Workspace, error) {
	path := filepath.Join(s.GetWorkspaceDir(id), "workspace.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workspace: %w", err)
	}
	return &ws, nil
}

// SaveWorkspace saves a workspace configuration.
// It automatically creates the workspace directory if it doesn't exist.
func (s *WorkspaceStore) SaveWorkspace(ws *Workspace) error {
	if ws.ProjectRoot == "" {
		return fmt.Errorf("projectRoot is required")
	}

	id := s.GetWorkspaceID(ws.ProjectRoot)
	dir := s.GetWorkspaceDir(id)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %w", err)
	}

	path := filepath.Join(dir, "workspace.json")
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workspace: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write workspace file: %w", err)
	}

	return nil
}
