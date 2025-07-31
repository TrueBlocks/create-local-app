package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"{{CHIFRA}}/pkg/file"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/preferences"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/project"
)

func TestFileNew(t *testing.T) {
	t.Run("CleanStateAllowsNewFile", func(t *testing.T) {
		app, dir, _ := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileNew()
		if err != nil {
			t.Fatalf("Expected no error for clean state, got: %v", err)
		}

		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected a new active project, but got nil")
		}

		if activeProject.GetPath() != "" {
			t.Fatalf("Expected empty path, got %s", activeProject.GetPath())
		}
	})

	t.Run("NewFileIsNotDirty", func(t *testing.T) {
		app, dir, _ := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileNew()
		if err != nil {
			t.Fatalf("Expected no error when creating new file, got: %v", err)
		}

		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected an active project after fileNew(), but got nil")
		}

		if activeProject.IsDirty() {
			t.Fatalf("Expected new project to not be dirty, but it was marked as dirty")
		}
	})
}

func TestFileSave(t *testing.T) {
	t.Run("BasicSaveCreatesFile", func(t *testing.T) {
		app, dir, filePath := getTestApp(t, true)
		defer preferences.SetConfigBaseForTest(t, dir)()

		if err := app.fileSave(); err != nil {
			t.Fatalf("fileSave() returned an error: %v", err)
		}

		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected an active project, but got nil")
		}

		if activeProject.IsDirty() {
			t.Fatalf("Expected dirty flag to be false, but it was true")
		}

		if !file.FileExists(filePath) {
			t.Fatalf("Expected file %s to be created, but it wasn't", filePath)
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to stat file: %v", err)
		}

		if fileInfo.Mode().Perm() != 0644 {
			t.Fatalf("Expected file permissions to be 0644, got: %v", fileInfo.Mode().Perm())
		}

		content, err := os.ReadFile(filePath)
		if err != nil || len(content) == 0 {
			t.Fatalf("Expected non-empty file content, got: %v", err)
		}
	})

	// New test for empty path
	t.Run("EmptyPathFails", func(t *testing.T) {
		app, dir, _ := getTestApp(t, true)
		defer preferences.SetConfigBaseForTest(t, dir)()

		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected an active project, but got nil")
		}

		activeProject.Path = ""

		err := app.fileSave()
		if !errors.Is(err, ErrEmptyFilePath) {
			t.Fatalf("Expected ErrEmptyFilePath, got: %v", err)
		}
	})

	// New test for non-dirty state
	t.Run("NotDirtyNoOp", func(t *testing.T) {
		app, dir, _ := getTestApp(t, false /* not dirty */)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileSave()
		if err != nil {
			t.Fatalf("Expected no error for non-dirty state, got: %v", err)
		}
	})
}

func TestFileSaveAs(t *testing.T) {
	t.Run("CreatesNewFile", func(t *testing.T) {
		app, dir, filePath := getTestApp(t, true)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileSaveAs(filePath, true)
		if err != nil {
			t.Fatalf("fileSaveAs() returned an error: %v", err)
		}

		if !file.FileExists(filePath) {
			t.Fatalf("Expected file %s to be created, but it wasn't", filePath)
		}

		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected an active project, but got nil")
		}

		if activeProject.GetPath() != filePath {
			t.Fatalf("Expected file path %s, got %s", filePath, activeProject.GetPath())
		}

		if activeProject.IsDirty() {
			t.Fatalf("Expected dirty flag to be false, but it was true")
		}
	})

	// New test for empty path
	t.Run("EmptyPathFails", func(t *testing.T) {
		app, dir, _ := getTestApp(t, true)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileSaveAs("", true)
		if !errors.Is(err, ErrEmptyFilePath) {
			t.Fatalf("Expected ErrEmptyFilePath, got: %v", err)
		}
	})

	// New test for overwrite protection
	t.Run("OverwriteNotConfirmedFails", func(t *testing.T) {
		app, dir, filePath := getTestApp(t, true)
		defer preferences.SetConfigBaseForTest(t, dir)()

		_ = os.WriteFile(filePath, []byte("{}"), 0644)

		err := app.fileSaveAs(filePath, false)
		if !errors.Is(err, ErrOverwriteNotConfirmed) {
			t.Fatalf("Expected ErrOverwriteNotConfirmed, got: %v", err)
		}
	})
}

func TestFileOpen(t *testing.T) {
	t.Run("ValidProjectFile", func(t *testing.T) {
		app, dir, filePath := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		proj := project.Project{
			Version: "1.0",
			Name:    strings.ToLower(preferences.GetAppId().OrgName) + "-project",
		}
		bytes, _ := json.Marshal(proj)
		_ = os.WriteFile(filePath, bytes, 0644)

		err := app.fileOpen(filePath)
		if err != nil {
			t.Fatalf("fileOpen() returned an error: %v", err)
		}

		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected an active project, but got nil")
		}

		if activeProject.GetPath() != filePath {
			t.Fatalf("Expected file path %s, got %s", filePath, activeProject.GetPath())
		}

		if activeProject.Version != "1.0" {
			t.Fatalf("Expected version '1.0', got %s", activeProject.Version)
		}
	})

	t.Run("NonexistentFile", func(t *testing.T) {
		app, dir, _ := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileOpen("./nonexistent_project.json")
		if !errors.Is(err, ErrFileNotFound) {
			t.Fatalf("Expected ErrFileNotFound, got: %v", err)
		}
	})

	// New test to verify file open doesn't mark file as dirty
	t.Run("FileOpenDoesNotMakeTheFileDirty", func(t *testing.T) {
		app, dir, filePath := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		// Set up a valid project file
		proj := project.Project{
			Version: "1.0",
			Name:    strings.ToLower(preferences.GetAppId().OrgName) + "-project",
		}
		bytes, _ := json.Marshal(proj)
		_ = os.WriteFile(filePath, bytes, 0644)

		// Open the file
		err := app.fileOpen(filePath)
		if err != nil {
			t.Fatalf("fileOpen() returned an error: %v", err)
		}

		// Check that the project is not dirty after opening
		activeProject := app.Projects.Active()
		if activeProject == nil {
			t.Fatalf("Expected an active project, but got nil")
		}

		if activeProject.IsDirty() {
			t.Fatalf("Expected project to not be dirty after opening, but it was marked as dirty")
		}
	})

	// New test for empty path
	t.Run("EmptyPathFails", func(t *testing.T) {
		app, dir, _ := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		err := app.fileOpen("")
		if !errors.Is(err, ErrEmptyFilePath) {
			t.Fatalf("Expected ErrEmptyFilePath, got: %v", err)
		}
	})

	t.Run("InvalidJSONDoesNotFail", func(t *testing.T) {
		app, dir, _ := getTestApp(t, false)
		defer preferences.SetConfigBaseForTest(t, dir)()

		// Use a different path for the invalid JSON file
		invalidFilePath := filepath.Join(dir, "invalid_project.json")

		// Write invalid JSON content to the file
		_ = os.WriteFile(invalidFilePath, []byte("{invalid json"), 0644)

		// Attempt to open the invalid file
		err := app.fileOpen(invalidFilePath)
		if err != nil {
			t.Fatalf("Expected no error opening invalid JSON, got: %v", err)
		}
	})
}

func TestUpdateRecentProjects(t *testing.T) {
	app, dir, _ := getTestApp(t, true)
	defer preferences.SetConfigBaseForTest(t, dir)()

	activeProject := app.Projects.Active()
	if activeProject == nil {
		t.Fatalf("Expected an active project, but got nil")
	}

	_ = app.fileSave()

	if len(app.Preferences.App.RecentProjects) == 0 {
		t.Fatalf("Expected recently used files list to contain entries, but it was empty")
	}

	ruf := strings.Replace(app.Preferences.App.RecentProjects[0], dir, ".", -1)
	if ruf != "./test_project.json" {
		t.Fatalf("Expected recently used file at position 0 to be './test_project.json', got %s", ruf)
	}
}

// Helper functions remain unchanged
func getTestApp(t *testing.T, dirty bool) (*App, string, string) {
	t.Helper()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "test_project.json")

	preferences := &preferences.Preferences{
		App: preferences.AppPreferences{
			RecentProjects: []string{},
		},
	}

	manager := project.NewManager()

	proj := manager.New("Test Project")
	proj.SetDirty(dirty)
	_ = proj.SaveAs(filePath)
	proj.SetDirty(dirty)

	app := &App{
		Preferences: preferences,
		Projects:    manager,
	}

	return app, dir, filePath
}
