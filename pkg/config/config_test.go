package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantArgs *Args
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "no arguments",
			args:     []string{"program"},
			wantArgs: &Args{IsAuto: false, IsReverse: false, TemplateName: ""},
			wantErr:  false,
		},
		{
			name:     "auto mode",
			args:     []string{"program", "--auto"},
			wantArgs: &Args{IsAuto: true, IsReverse: false, TemplateName: ""},
			wantErr:  false,
		},
		{
			name:     "reverse mode with valid template name",
			args:     []string{"program", "--reverse", "my-template-123"},
			wantArgs: &Args{IsAuto: false, IsReverse: true, TemplateName: "my-template-123"},
			wantErr:  false,
		},
		{
			name:    "reverse mode missing template name",
			args:    []string{"program", "--reverse"},
			wantErr: true,
			errMsg:  "--reverse requires a template name parameter",
		},
		{
			name:    "reverse mode with invalid template name - spaces",
			args:    []string{"program", "--reverse", "invalid name"},
			wantErr: true,
			errMsg:  "invalid template name 'invalid name': must start with alphanumeric and contain only alphanumeric characters and dashes",
		},
		{
			name:    "reverse mode with invalid template name - starts with dash",
			args:    []string{"program", "--reverse", "-invalid"},
			wantErr: true,
			errMsg:  "invalid template name '-invalid': must start with alphanumeric and contain only alphanumeric characters and dashes",
		},
		{
			name:    "reverse mode with invalid template name - special chars",
			args:    []string{"program", "--reverse", "invalid@name"},
			wantErr: true,
			errMsg:  "invalid template name 'invalid@name': must start with alphanumeric and contain only alphanumeric characters and dashes",
		},
		{
			name:    "unknown argument",
			args:    []string{"program", "--unknown"},
			wantErr: true,
			errMsg:  "unknown argument: --unknown (valid options: --reverse <template-name>, --auto)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set test args
			os.Args = tt.args

			args, err := ParseArgs()

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseArgs() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ParseArgs() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseArgs() unexpected error = %v", err)
				return
			}

			if args.IsAuto != tt.wantArgs.IsAuto {
				t.Errorf("ParseArgs() IsAuto = %v, want %v", args.IsAuto, tt.wantArgs.IsAuto)
			}
			if args.IsReverse != tt.wantArgs.IsReverse {
				t.Errorf("ParseArgs() IsReverse = %v, want %v", args.IsReverse, tt.wantArgs.IsReverse)
			}
			if args.TemplateName != tt.wantArgs.TemplateName {
				t.Errorf("ParseArgs() TemplateName = %v, want %v", args.TemplateName, tt.wantArgs.TemplateName)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test configs
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name       string
		configData string
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "valid config file",
			configData: `{
  "Organization": "TrueBlocks, LLC",
  "ProjectName": "my-app",
  "Github": "https://github.com/TrueBlocks/my-app",
  "Domain": "trueblocks.io"
}`,
			wantConfig: &Config{
				Organization: "TrueBlocks, LLC",
				ProjectName:  "my-app",
				Github:       "https://github.com/TrueBlocks/my-app",
				Domain:       "trueblocks.io",
			},
			wantErr: false,
		},
		{
			name: "partial config file",
			configData: `{
  "Organization": "Test Org",
  "ProjectName": "test-project"
}`,
			wantConfig: &Config{
				Organization: "Test Org",
				ProjectName:  "test-project",
				Github:       "",
				Domain:       "",
			},
			wantErr: false,
		},
		{
			name:       "invalid JSON",
			configData: `{"Organization": "Test", invalid json}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, "test-config.json")

			// Write test config file
			if err := os.WriteFile(configPath, []byte(tt.configData), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			config, err := LoadConfig(configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfig() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadConfig() unexpected error = %v", err)
				return
			}

			if config.Organization != tt.wantConfig.Organization {
				t.Errorf("LoadConfig() Organization = %v, want %v", config.Organization, tt.wantConfig.Organization)
			}
			if config.ProjectName != tt.wantConfig.ProjectName {
				t.Errorf("LoadConfig() ProjectName = %v, want %v", config.ProjectName, tt.wantConfig.ProjectName)
			}
			if config.Github != tt.wantConfig.Github {
				t.Errorf("LoadConfig() Github = %v, want %v", config.Github, tt.wantConfig.Github)
			}
			if config.Domain != tt.wantConfig.Domain {
				t.Errorf("LoadConfig() Domain = %v, want %v", config.Domain, tt.wantConfig.Domain)
			}
		})
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// Test loading non-existent config file should return empty config
	nonExistentPath := "/path/that/does/not/exist/config.json"

	config, err := LoadConfig(nonExistentPath)
	if err != nil {
		t.Errorf("LoadConfig() with non-existent file should not error, got: %v", err)
	}

	// Should return empty config
	expected := &Config{}
	if config.Organization != expected.Organization ||
		config.ProjectName != expected.ProjectName ||
		config.Github != expected.Github ||
		config.Domain != expected.Domain {
		t.Errorf("LoadConfig() with non-existent file should return empty config, got: %+v", config)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for test configs
	tempDir, err := os.MkdirTemp("", "config-save-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &Config{
		Organization: "Test Organization",
		ProjectName:  "test-project",
		Github:       "https://github.com/test/test-project",
		Domain:       "test.io",
	}

	configPath := filepath.Join(tempDir, "test-save-config.json")

	// Save config
	err = SaveConfig(configPath, config)
	if err != nil {
		t.Errorf("SaveConfig() unexpected error = %v", err)
		return
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("SaveConfig() did not create config file")
		return
	}

	// Load it back and verify
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("Failed to load saved config: %v", err)
		return
	}

	if loadedConfig.Organization != config.Organization {
		t.Errorf("Saved/loaded Organization mismatch: got %v, want %v", loadedConfig.Organization, config.Organization)
	}
	if loadedConfig.ProjectName != config.ProjectName {
		t.Errorf("Saved/loaded ProjectName mismatch: got %v, want %v", loadedConfig.ProjectName, config.ProjectName)
	}
	if loadedConfig.Github != config.Github {
		t.Errorf("Saved/loaded Github mismatch: got %v, want %v", loadedConfig.Github, config.Github)
	}
	if loadedConfig.Domain != config.Domain {
		t.Errorf("Saved/loaded Domain mismatch: got %v, want %v", loadedConfig.Domain, config.Domain)
	}
}

func TestGetConfigPath(t *testing.T) {
	configPath, err := GetConfigPath()
	if err != nil {
		t.Errorf("GetConfigPath() unexpected error = %v", err)
		return
	}

	if configPath == "" {
		t.Errorf("GetConfigPath() returned empty path")
		return
	}

	// Should end with the expected config filename
	expectedFilename := "config.json"
	if filepath.Base(configPath) != expectedFilename {
		t.Errorf("GetConfigPath() should end with %s, got %s", expectedFilename, filepath.Base(configPath))
	}

	// Should be an absolute path
	if !filepath.IsAbs(configPath) {
		t.Errorf("GetConfigPath() should return absolute path, got %s", configPath)
	}

	// Should be in the user's home directory under .create-local-app
	expectedDir := ".create-local-app"
	if !strings.Contains(configPath, expectedDir) {
		t.Errorf("GetConfigPath() should contain %s, got %s", expectedDir, configPath)
	}
}

// TestIsValidTemplateName tests the template name validation logic
// This is a table-driven test for the internal validation function
func TestTemplateNameValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"valid alphanumeric", []string{"program", "--reverse", "template123"}, false},
		{"valid with dashes", []string{"program", "--reverse", "my-template-123"}, false},
		{"valid starting with letter", []string{"program", "--reverse", "a-template"}, false},
		{"valid starting with number", []string{"program", "--reverse", "1template"}, false},
		{"invalid with spaces", []string{"program", "--reverse", "my template"}, true},
		{"invalid starting with dash", []string{"program", "--reverse", "-template"}, true},
		{"invalid with special chars", []string{"program", "--reverse", "template@123"}, true},
		{"invalid with underscore", []string{"program", "--reverse", "template_123"}, true},
		{"invalid empty", []string{"program", "--reverse", ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = tt.args
			_, err := ParseArgs()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error for template name validation but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for valid template name: %v", err)
			}
		})
	}
}
