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
	IsForce      bool
	TemplateName string
}

// ParseArgs parses command line arguments and returns Args struct or handles special commands
func ParseArgs(version string) (*Args, error) {
	args := &Args{}

	if len(os.Args) > 1 {
		i := 1
		for i < len(os.Args) {
			switch os.Args[i] {
			case "--version", "version":
				fmt.Printf("create-local-app version %s\n", version)
				os.Exit(0)
			case "--help", "help":
				printHelp()
				os.Exit(0)
			case "--reverse":
				if i+1 >= len(os.Args) {
					return nil, fmt.Errorf("--reverse requires a template name parameter")
				}
				templateName := os.Args[i+1]
				if !isValidTemplateName(templateName) {
					return nil, fmt.Errorf("invalid template name '%s': must start with alphanumeric and contain only alphanumeric characters and dashes", templateName)
				}
				args.IsReverse = true
				args.TemplateName = templateName
				i += 2 // Skip the template name argument
			case "--auto":
				args.IsAuto = true
				i++
			case "--force":
				args.IsForce = true
				i++
			default:
				return nil, fmt.Errorf("unknown argument: %s (valid options: --reverse <template-name>, --auto, --force, --version, --help)", os.Args[i])
			}
		}
	}
	return args, nil
}

// printHelp displays usage information
func printHelp() {
	fmt.Println("create-local-app - A powerful Go-based scaffolding tool for TrueBlocks/Wails desktop applications")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  create-local-app [options]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --auto                    Use saved configuration without prompts")
	fmt.Println("  --reverse <template-name> Create a template from the current directory")
	fmt.Println("  --force                   Force operation without confirmation (overwrite existing files)")
	fmt.Println("  --version                 Show version information")
	fmt.Println("  --help                    Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  create-local-app                    # Interactive mode - prompts for project details")
	fmt.Println("  create-local-app --auto             # Use previously saved configuration")
	fmt.Println("  create-local-app --force            # Interactive mode - overwrite existing files without confirmation")
	fmt.Println("  create-local-app --reverse my-app   # Create template from current directory")
	fmt.Println("  create-local-app --reverse my-app --force  # Create template, overwrite if exists")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/TrueBlocks/create-local-app")
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
