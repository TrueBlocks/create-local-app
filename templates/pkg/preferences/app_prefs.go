package preferences

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"{{CHIFRA}}/pkg/file"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/logging"
	"github.com/kbinani/screenshot"
)

type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewBounds() Bounds {
	bounds := screenshot.GetDisplayBounds(0)
	screenW := bounds.Dx()
	screenH := bounds.Dy()

	ret := Bounds{}
	ret.Width = screenW * 3 / 4
	ret.Height = screenH * 3 / 4
	ret.X = (screenW - ret.Width) / 2
	ret.Y = (screenH - ret.Height) / 2
	return ret
}

func (b *Bounds) IsValid() bool {
	return b.X >= 0 && b.Y >= 0 && b.Width > 100 && b.Height > 100
}

type AppPreferences struct {
	Bounds           Bounds            `json:"bounds,omitempty"`
	DebugMode        bool              `json:"debugMode,omitempty"`
	HelpCollapsed    bool              `json:"helpCollapsed,omitempty"`
	LastAddress      string            `json:"lastAddress,omitempty"`
	LastChain        string            `json:"lastChain,omitempty"`
	LastContract     string            `json:"lastContract,omitempty"`
	LastLanguage     string            `json:"lastLanguage,omitempty"`
	LastProject      string            `json:"lastProject,omitempty"`
	LastTab          map[string]string `json:"lastTab"`
	LastTheme        string            `json:"lastTheme,omitempty"`
	LastView         string            `json:"lastView,omitempty"`
	LastViewNoWizard string            `json:"lastViewNoWizard,omitempty"`
	MenuCollapsed    bool              `json:"menuCollapsed,omitempty"`
	Name             string            `json:"name,omitempty"`
	RecentProjects   []string          `json:"recentProjects"`
	SilencedDialogs  map[string]bool   `json:"silencedDialogs"`
	Version          string            `json:"version"`
}

func (p *AppPreferences) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

// NewAppPreferences creates a new AppPreferences instance with default values
func NewAppPreferences() *AppPreferences {
	return &AppPreferences{
		Bounds:           NewBounds(),
		DebugMode:        false,
		HelpCollapsed:    false,
		LastAddress:      "0xf503017d7baf7fbc0fff7492b751025c6a78179b",
		LastContract:     "0x52df6e4d9989e7cf4739d687c765e75323a1b14c",
		LastLanguage:     "en",
		LastTab:          make(map[string]string),
		LastTheme:        "dark",
		LastView:         "/",
		LastViewNoWizard: "/",
		MenuCollapsed:    false,
		RecentProjects:   []string{},
		SilencedDialogs:  make(map[string]bool),
		Version:          "1.0",
	}
}

func GetAppPreferences() (AppPreferences, error) {
	path := getAppPrefsPath()

	if !file.FileExists(path) {
		defaults := *NewAppPreferences()
		if err := SetAppPreferences(&defaults); err != nil {
			return AppPreferences{}, err
		}
		return defaults, nil
	}

	var appPrefs AppPreferences
	contents := file.AsciiFileToString(path)
	if err := json.Unmarshal([]byte(contents), &appPrefs); err != nil {
		// Log the corruption issue for debugging
		logging.LogBackend(fmt.Sprintf("Warning: App preferences file corrupted (%v), creating new defaults", err))
		logging.LogBackend(fmt.Sprintf("Corrupted content: %s", contents))
		backupPath := path + ".corrupted"
		if backupErr := os.WriteFile(backupPath, []byte(contents), 0644); backupErr == nil {
			logging.LogBackend(fmt.Sprintf("Corrupted file backed up to: %s", backupPath))
		}

		appPrefs = *NewAppPreferences()
		if err = SetAppPreferences(&appPrefs); err != nil {
			return AppPreferences{}, fmt.Errorf("failed to save repaired preferences: %w", err)
		}
		logging.LogBackend("App preferences reset to defaults and saved")
	}

	var needsSave bool
	if appPrefs.LastTab == nil {
		appPrefs.LastTab = make(map[string]string)
		needsSave = true
	}
	if appPrefs.RecentProjects == nil {
		appPrefs.RecentProjects = []string{}
		needsSave = true
	}
	if appPrefs.SilencedDialogs == nil {
		appPrefs.SilencedDialogs = make(map[string]bool)
		needsSave = true
	}
	if appPrefs.Version == "" {
		appPrefs.Version = "1.0"
		needsSave = true
	}
	if appPrefs.LastView == "" {
		appPrefs.LastView = "/"
		needsSave = true
	}
	if appPrefs.LastViewNoWizard == "" {
		appPrefs.LastViewNoWizard = "/"
		needsSave = true
	}

	if needsSave {
		logging.LogBackend("App preferences had missing fields, saving corrected version")
		if err := SetAppPreferences(&appPrefs); err != nil {
			logging.LogBackend(fmt.Sprintf("Warning: Could not save corrected app preferences: %v", err))
		}
	}

	if tbDebug := os.Getenv("TB_DEBUG"); tbDebug != "" {
		envDebugMode := tbDebug == "true" || tbDebug == "1" || tbDebug == "on"
		if appPrefs.DebugMode != envDebugMode {
			logging.LogBackend(fmt.Sprintf("TB_DEBUG environment variable (%s) overriding stored debug mode (%t) with %t", tbDebug, appPrefs.DebugMode, envDebugMode))
			appPrefs.DebugMode = envDebugMode
			if err := SetAppPreferences(&appPrefs); err != nil {
				logging.LogBackend(fmt.Sprintf("Warning: Could not save debug mode from TB_DEBUG: %v", err))
			}
		}
	}

	return appPrefs, nil
}

func SetAppPreferences(appPrefs *AppPreferences) error {
	path := getAppPrefsPath()

	data, err := json.MarshalIndent(appPrefs, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getAppPrefsPath() string {
	return filepath.Join(getConfigBase(), ToCamel(configBaseApp), "app_prefs.json")
}
