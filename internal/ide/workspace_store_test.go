package ide

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWorkspaceStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewWorkspaceStore(tmpDir)
	projectRoot := "/tmp/my-project"
	id := store.GetWorkspaceID(projectRoot)

	ws := &Workspace{
		Version:      "1.0",
		ProjectRoot:  projectRoot,
		DisplayName:  "My Project",
		CreatedAt:    time.Now(),
		LastOpenedAt: time.Now(),
	}

	if err := store.SaveWorkspace(ws); err != nil {
		t.Fatalf("SaveWorkspace failed: %v", err)
	}

	loadedWs, err := store.LoadWorkspace(id)
	if err != nil {
		t.Fatalf("LoadWorkspace failed: %v", err)
	}

	if loadedWs.ProjectRoot != ws.ProjectRoot {
		t.Errorf("expected ProjectRoot %s, got %s", ws.ProjectRoot, loadedWs.ProjectRoot)
	}
	if loadedWs.DisplayName != ws.DisplayName {
		t.Errorf("expected DisplayName %s, got %s", ws.DisplayName, loadedWs.DisplayName)
	}

	// Verify directory structure
	expectedDir := filepath.Join(tmpDir, id)
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("workspace directory not created: %s", expectedDir)
	}
}
