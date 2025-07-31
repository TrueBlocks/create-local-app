package app

import (
	"context"
	"testing"

	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/preferences"
)

func TestIsReady(t *testing.T) {
	app := &App{}

	if app.IsReady() {
		t.Fatal("Expected IsReady() to be false when ctx is nil")
	}

	app.ctx = context.Background()
	if !app.IsReady() {
		t.Fatal("Expected IsReady() to be true when ctx is set")
	}
}

func TestSaveBounds_SavesToPreferences(t *testing.T) {
	tmp := t.TempDir()
	defer preferences.SetConfigBaseForTest(t, tmp)()

	app := &App{
		ctx: context.Background(),
		Preferences: &preferences.Preferences{
			App: preferences.AppPreferences{},
		},
	}

	app.SaveBounds(111, 222, 333, 444)

	loaded, err := preferences.GetAppPreferences()
	if err != nil {
		t.Fatalf("Failed to load saved preferences: %v", err)
	}

	if loaded.Bounds.X != 111 || loaded.Bounds.Y != 222 || loaded.Bounds.Width != 333 || loaded.Bounds.Height != 444 {
		t.Fatalf("Unexpected saved bounds:\n  x: %d\n  y: %d\n  width: %d\n  height: %d",
			loaded.Bounds.X, loaded.Bounds.Y, loaded.Bounds.Width, loaded.Bounds.Height)
	}
}
