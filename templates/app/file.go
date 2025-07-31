package app

import (
	"errors"
	"fmt"
	"os"

	"{{CHIFRA}}/pkg/file"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/msgs"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/project"
)

// Generic errors
var ErrEmptyFilePath = errors.New("empty file path")
var ErrUnsavedChanges = errors.New("unsaved changes")
var ErrFileNotFound = errors.New("file not found")
var ErrOverwriteNotConfirmed = errors.New("file exists, overwrite not confirmed")

// File operation errors
var ErrReadFileFailed = errors.New("failed to read file")
var ErrWriteFileFailed = errors.New("failed to write file")
var ErrSerializeFailed = errors.New("failed to serialize data")
var ErrDeserializeFailed = errors.New("failed to deserialize data")

func (a *App) fileNew() error {
	a.Projects.New(a.uniqueProjectName("New Project"))

	// Ensure the newly created project is not marked as dirty
	activeProject := a.Projects.Active()
	if activeProject != nil {
		activeProject.SetDirty(false)
	}

	a.updateRecentProjects()
	return nil
}

func (a *App) fileSave() error {
	project := a.Projects.Active()
	if project == nil {
		return errors.New("no active project")
	}

	if project.GetPath() == "" {
		return ErrEmptyFilePath
	}

	if !project.IsDirty() {
		return nil
	}

	if err := a.Projects.SaveActive(); err != nil {
		return err
	}

	a.updateRecentProjects()
	return nil
}

func (a *App) fileSaveAs(newPath string, overwriteConfirmed bool) error {
	project := a.Projects.Active()
	if project == nil {
		return errors.New("no active project")
	}

	if newPath == "" {
		return ErrEmptyFilePath
	}

	if file.FileExists(newPath) && !overwriteConfirmed {
		return ErrOverwriteNotConfirmed
	}

	if err := a.Projects.SaveActiveAs(newPath); err != nil {
		return err
	}

	a.updateRecentProjects()
	return nil
}

func (a *App) fileOpen(path string) error {
	if path == "" {
		return ErrEmptyFilePath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrFileNotFound
	}

	_, err := a.Projects.Open(path)
	if err != nil {
		return err
	}

	a.updateRecentProjects()
	return nil
}

func (a *App) updateRecentProjects() {
	activeProject := a.Projects.Active()
	if activeProject == nil || activeProject.GetPath() == "" {
		return
	}

	path := activeProject.GetPath()

	if err := a.Preferences.AddRecentProject(path); err != nil {
		msgs.EmitError("add recent project failed", err)
		return
	}

	msgs.EmitManager("update_recent_projects")
}

func (a *App) GetFilename() *project.Project {
	return a.Projects.Active()
}

func (a *App) uniqueProjectName(baseName string) string {
	projectExists := func(name string) bool {
		for _, project := range a.Projects.GetOpenProjectIDs() {
			projectObj := a.Projects.GetProjectByID(project)
			if projectObj.Name == name {
				return true
			}
		}
		return false
	}

	count := 1
	uniqueName := baseName
	for {
		if !projectExists(uniqueName) {
			break
		}
		uniqueName = baseName + " " + fmt.Sprintf("%d", count)
		count++
	}

	return uniqueName
}
