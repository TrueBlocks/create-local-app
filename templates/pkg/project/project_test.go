package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/project"
)

// TestProject tests the basic functionality of the Project type
func TestProject(t *testing.T) {
	// Create a new project
	p := project.New("test-project")

	// Set some data
	p.SetData(map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	})

	// Verify initial state
	if !p.IsDirty() {
		t.Error("New project should be marked as dirty")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "project-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save the project
	tempPath := filepath.Join(tempDir, "test-project.json")
	if err := p.SaveAs(tempPath); err != nil {
		t.Fatalf("Failed to save project: %v", err)
	}

	// Verify the project is no longer dirty
	if p.IsDirty() {
		t.Error("Project should not be dirty after saving")
	}

	// Modify the project
	p.SetName("renamed-project")

	// Verify the project is dirty again
	if !p.IsDirty() {
		t.Error("Project should be dirty after modification")
	}

	// Save again
	if err := p.Save(); err != nil {
		t.Fatalf("Failed to save project after renaming: %v", err)
	}

	// Load the project
	loadedProject, err := project.Load(tempPath)
	if err != nil {
		t.Fatalf("Failed to load project: %v", err)
	}

	// Verify loaded data
	if loadedProject.GetName() != "renamed-project" {
		t.Errorf("Expected project name '%s', got '%s'", "renamed-project", loadedProject.GetName())
	}
}
