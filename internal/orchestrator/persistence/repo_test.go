package persistence

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWorkspaceRepository(t *testing.T) {
	// Setup temp dir
	tmpDir, err := os.MkdirTemp("", "multiverse-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	repo := NewWorkspaceRepository(tmpDir)

	// Test Init
	if err := repo.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Verify dirs created
	dirs := []string{"design", "design/nodes", "state", "history", "snapshots"}
	for _, d := range dirs {
		if _, err := os.Stat(filepath.Join(tmpDir, d)); os.IsNotExist(err) {
			t.Errorf("Directory %s not created", d)
		}
	}

	// --- Design Test ---
	t.Run("DesignRepository", func(t *testing.T) {
		wbs := &WBS{
			WBSID:       "wbs-1",
			ProjectRoot: "/tmp/foo",
			NodeIndex: []NodeIndex{
				{NodeID: "n1", Children: []string{}},
			},
		}
		if err := repo.Design().SaveWBS(wbs); err != nil {
			t.Fatalf("SaveWBS failed: %v", err)
		}

		loadedWBS, err := repo.Design().LoadWBS()
		if err != nil {
			t.Fatalf("LoadWBS failed: %v", err)
		}
		if loadedWBS.WBSID != wbs.WBSID {
			t.Errorf("Expected WBSID %s, got %s", wbs.WBSID, loadedWBS.WBSID)
		}

		node := &NodeDesign{
			NodeID: "n1",
			Name:   "Test Node",
		}
		if err := repo.Design().SaveNode(node); err != nil {
			t.Fatalf("SaveNode failed: %v", err)
		}

		loadedNode, err := repo.Design().GetNode("n1")
		if err != nil {
			t.Fatalf("GetNode failed: %v", err)
		}
		if loadedNode.Name != node.Name {
			t.Errorf("Expected Node Name %s, got %s", node.Name, loadedNode.Name)
		}
	})

	// --- State Test ---
	t.Run("StateRepository", func(t *testing.T) {
		tasks := &TasksState{
			Tasks: []TaskState{
				{TaskID: "t1", Status: "pending"},
			},
		}
		if err := repo.State().SaveTasks(tasks); err != nil {
			t.Fatalf("SaveTasks failed: %v", err)
		}

		loadedTasks, err := repo.State().LoadTasks()
		if err != nil {
			t.Fatalf("LoadTasks failed: %v", err)
		}
		if len(loadedTasks.Tasks) != 1 || loadedTasks.Tasks[0].TaskID != "t1" {
			t.Errorf("Unexpected loaded tasks: %+v", loadedTasks)
		}
	})

	// --- History Test ---
	t.Run("HistoryRepository", func(t *testing.T) {
		now := time.Now()
		action := &Action{
			ID:   "act-1",
			At:   now,
			Kind: "test",
		}
		if err := repo.History().AppendAction(action); err != nil {
			t.Fatalf("AppendAction failed: %v", err)
		}

		actions, err := repo.History().ListActions(now.Add(-1*time.Hour), now.Add(1*time.Hour))
		if err != nil {
			t.Fatalf("ListActions failed: %v", err)
		}
		if len(actions) != 1 {
			t.Fatalf("Expected 1 action, got %d", len(actions))
		}
		if actions[0].ID != "act-1" {
			t.Errorf("Expected action ID act-1, got %s", actions[0].ID)
		}
	})
}
