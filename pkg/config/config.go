package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Config represents the application configuration
type Config struct {
	Organization string `json:"Organization"`
	ProjectName  string `json:"ProjectName"`
	Github       string `json:"Github"`
	Domain       string `json:"Domain"`
}

// Args represents parsed command line arguments
type Args struct {
	IsAuto       bool
	IsReverse    bool
	TemplateName string
}

// ParseArgs parses command line arguments
func ParseArgs() (*Args, error) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--reverse":
			if len(os.Args) < 3 {
				return nil, fmt.Errorf("--reverse requires a template name parameter")
			}
			templateName := os.Args[2]
			if !isValidTemplateName(templateName) {
				return nil, fmt.Errorf("invalid template name '%s': must start with alphanumeric and contain only alphanumeric characters and dashes", templateName)
			}
			return &Args{IsReverse: true, TemplateName: templateName}, nil
		case "--auto":
			return &Args{IsAuto: true}, nil
		default:
			return nil, fmt.Errorf("unknown argument: %s (valid options: --reverse <template-name>, --auto)", os.Args[1])
		}
	}
	return &Args{}, nil
}

// isValidTemplateName validates that a template name starts with alphanumeric
// and contains only alphanumeric characters and dashes
func isValidTemplateName(name string) bool {
	if name == "" {
		return false
	}

	// Template name must start with alphanumeric and contain only alphanumeric and dashes
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9-]*$`, name)
	return matched
}

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil // Return empty config if file doesn't exist
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(configPath string, config *Config) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetUserConfigDir returns the path to the user's configuration directory
func GetUserConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".create-local-app"), nil
}

// InitializeUserConfig creates the user configuration directory structure if it doesn't exist
func InitializeUserConfig() error {
	configDir, err := GetUserConfigDir()
	if err != nil {
		return err
	}

	// Create main config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Create templates directory structure
	systemTemplatesDir := filepath.Join(configDir, "templates", "system")
	contributedTemplatesDir := filepath.Join(configDir, "templates", "contributed")

	if err := os.MkdirAll(systemTemplatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create system templates directory %s: %w", systemTemplatesDir, err)
	}

	if err := os.MkdirAll(contributedTemplatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create contributed templates directory %s: %w", contributedTemplatesDir, err)
	}

	return nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	configDir, err := GetUserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}
