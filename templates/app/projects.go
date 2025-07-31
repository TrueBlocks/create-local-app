package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SwitchToProject(id string) error {
	if a.Projects.GetProjectByID(id) == nil {
		return fmt.Errorf("no project with ID %s exists", id)
	}
	return a.Projects.SetActive(id)
}

func (a *App) CloseProject(id string) error {
	project := a.Projects.GetProjectByID(id)
	if project == nil {
		return fmt.Errorf("no project with ID %s exists", id)
	}

	// Check if project has unsaved changes
	if project.IsDirty() {
		response, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Unsaved Changes",
			Message: fmt.Sprintf("Do you want to save changes to project '%s' before closing?", project.GetName()),
			Buttons: []string{"Yes", "No", "Cancel"},
		})

		if err != nil {
			return err
		}

		switch response {
		case "Yes":
			// Save the project before closing
			if project.GetPath() == "" {
				// Project hasn't been saved before, need to use SaveAs
				path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
					Title: "Save Project Before Closing",
				})
				if err != nil || path == "" {
					return fmt.Errorf("save canceled")
				}

				// Use project's SaveAs method instead of SaveProjectAs
				if err := project.SaveAs(path); err != nil {
					return err
				}
			} else {
				// Project has a path, use normal save
				// Use project's Save method instead of SaveProject
				if err := project.Save(); err != nil {
					return err
				}
			}
		case "Cancel":
			return fmt.Errorf("close canceled")
		}
	}

	return a.Projects.Close(id)
}
