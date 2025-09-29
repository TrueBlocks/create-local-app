package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ViewConfigEntry represents configuration for a single view in .create-local-app.json
type ViewConfigEntry struct {
	MenuOrder      int             `json:"menuOrder,omitempty"`
	Disabled       bool            `json:"disabled,omitempty"`
	DisabledFacets map[string]bool `json:"disabledFacets,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Organization  string                     `json:"Organization"`
	ProjectName   string                     `json:"ProjectName"`
	Github        string                     `json:"Github"`
	Domain        string                     `json:"Domain"`
	Template      string                     `json:"Template"`
	PreserveFiles []string                   `json:"PreserveFiles,omitempty"`
	ViewConfig    map[string]ViewConfigEntry `json:"ViewConfig,omitempty"`
}

// Args represents parsed command line arguments
type Args struct {
	IsAuto       bool
	IsCreate     bool
	IsRemove     bool
	IsForce      bool
	IsList       bool
	TemplateName string
	UseTemplate  string
}

// ParseArgs parses command line arguments and returns Args struct or handles special commands
func ParseArgs(version, buildTime string) (*Args, error) {
	args := &Args{}

	if len(os.Args) > 1 {
		i := 1
		for i < len(os.Args) {
			switch os.Args[i] {
			case "--version", "version":
				fmt.Printf("create-local-app version %s\n", version)
				fmt.Printf("built: %s\n", buildTime)
				os.Exit(0)
			case "--help", "help":
				printHelp()
				os.Exit(0)
			case "--create":
				if i+1 >= len(os.Args) {
					return nil, fmt.Errorf("--create requires a template name parameter")
				}
				templateName := os.Args[i+1]
				if !isValidTemplateName(templateName) {
					return nil, fmt.Errorf("invalid template name '%s': must start with alphanumeric and contain only alphanumeric characters and dashes", templateName)
				}
				args.IsCreate = true
				args.TemplateName = templateName
				i += 2 // Skip the template name argument
			case "--remove":
				if i+1 >= len(os.Args) {
					return nil, fmt.Errorf("--remove requires a template name parameter")
				}
				templateName := os.Args[i+1]
				if !isValidTemplateName(templateName) {
					return nil, fmt.Errorf("invalid template name '%s': must start with alphanumeric and contain only alphanumeric characters and dashes", templateName)
				}
				args.IsRemove = true
				args.TemplateName = templateName
				i += 2 // Skip the template name argument
			case "--auto":
				args.IsAuto = true
				i++
			case "--force":
				args.IsForce = true
				i++
			case "--list":
				args.IsList = true
				i++
			case "--template":
				if i+1 >= len(os.Args) {
					return nil, fmt.Errorf("--template requires a template name parameter")
				}
				templateName := os.Args[i+1]
				if !isValidTemplateName(templateName) {
					return nil, fmt.Errorf("invalid template name '%s': must start with alphanumeric and contain only alphanumeric characters and dashes", templateName)
				}
				args.UseTemplate = templateName
				i += 2 // Skip the template name argument
			default:
				return nil, fmt.Errorf("unknown argument: %s (valid options: --create <template-name>, --remove <template-name>, --template <template-name>, --auto, --force, --list, --version, --help)", os.Args[i])
			}
		}
	}

	// Validate incompatible flag combinations
	if args.IsCreate && args.IsAuto {
		return nil, fmt.Errorf("--create and --auto flags are incompatible (auto mode is for project creation, not template creation)")
	}
	if args.IsCreate && args.IsForce {
		return nil, fmt.Errorf("--create and --force flags are incompatible (force mode is for project creation, not template creation)")
	}
	if args.IsRemove && args.IsAuto {
		return nil, fmt.Errorf("--remove and --auto flags are incompatible (auto mode is for project creation, not template removal)")
	}
	if args.IsRemove && args.IsForce {
		return nil, fmt.Errorf("--remove and --force flags are incompatible (force mode is for project creation, not template removal)")
	}
	if args.IsCreate && args.IsRemove {
		return nil, fmt.Errorf("--create and --remove flags are incompatible (cannot create and remove template simultaneously)")
	}
	if args.UseTemplate != "" && args.IsCreate {
		return nil, fmt.Errorf("--template and --create flags are incompatible (cannot specify template when creating one)")
	}
	if args.UseTemplate != "" && args.IsRemove {
		return nil, fmt.Errorf("--template and --remove flags are incompatible (cannot specify template when removing one)")
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
	fmt.Println("  --auto                           Use saved configuration without prompts")
	fmt.Println("  --list                           List available templates")
	fmt.Println("  --create <template-name>         Create a template from the current directory")
	fmt.Println("  --remove <template-name>         Remove a contributed template")
	fmt.Println("  --template <template-name>       Optionally, use a specific template")
	fmt.Println("  --force                          Force operation without confirmation (overwrite existing files)")
	fmt.Println("  --version                        Show version information")
	fmt.Println("  --help                           Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  create-local-app                           # Interactive mode - prompts for project details")
	fmt.Println("  create-local-app --auto                    # Use previously saved configuration")
	fmt.Println("  create-local-app --list                    # List available templates")
	fmt.Println("  create-local-app --create my-template      # Create template from current directory")
	fmt.Println("  create-local-app --remove my-template      # Remove contributed template")
	fmt.Println("  create-local-app --template my-template    # Use a specific template")
	fmt.Println("  create-local-app --force                   # Overwrite existing files without confirmation")
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

// GetProjectConfigPath returns the path to the project-local configuration file
func GetProjectConfigPath() string {
	return "./.create-local-app.json"
}

// LoadProjectConfig loads configuration from project-local file, falling back to global config
func LoadProjectConfig() (*Config, string, error) {
	// First try project-local config
	projectConfigPath := GetProjectConfigPath()
	if _, err := os.Stat(projectConfigPath); err == nil {
		config, err := LoadConfig(projectConfigPath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load project config: %w", err)
		}
		return config, projectConfigPath, nil
	}

	// Fall back to global config
	globalConfigPath, err := GetConfigPath()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get global config path: %w", err)
	}

	config, err := LoadConfig(globalConfigPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load global config: %w", err)
	}

	return config, globalConfigPath, nil
}

// SaveProjectConfig saves configuration to project-local file
func SaveProjectConfig(config *Config) error {
	return SaveConfig(GetProjectConfigPath(), config)
}

// SaveGlobalConfig saves configuration to global config file
func SaveGlobalConfig(config *Config) error {
	globalConfigPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get global config path: %w", err)
	}
	return SaveConfig(globalConfigPath, config)
}
