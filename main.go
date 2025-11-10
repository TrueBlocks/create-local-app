package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/TrueBlocks/create-local-app/pkg/config"
	"github.com/TrueBlocks/create-local-app/pkg/processor"
	"github.com/TrueBlocks/create-local-app/pkg/templates"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

//go:embed templates/system/*.tar.gz
var systemTemplatesFS embed.FS

//go:embed VERSION
var versionContent string

func main() {
	version := strings.TrimSpace(versionContent)
	built := file.MustGetLatestFileTime("VERSION")
	args, err := config.ParseArgs(version, built.Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Initialize user configuration directory structure
	if err := config.InitializeUserConfig(); err != nil {
		fmt.Printf("Error initializing user config: %v\n", err)
		os.Exit(1)
	}

	// Initialize system templates on first run or if version changed
	if err := templates.InitializeSystemTemplates(systemTemplatesFS, version); err != nil {
		fmt.Printf("Error initializing system templates: %v\n", err)
		os.Exit(1)
	}

	// Handle list templates mode
	if args.IsList {
		if err := templates.ListTemplates(); err != nil {
			fmt.Printf("Error listing templates: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle remove template mode
	if args.IsRemove {
		configDir, err := config.GetUserConfigDir()
		if err != nil {
			fmt.Printf("Error getting user config directory: %v\n", err)
			os.Exit(1)
		}

		templatePath := filepath.Join(configDir, "templates", "contributed", args.TemplateName)

		// Check if template exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			fmt.Printf("Error: Template '%s' not found in contributed templates.\n", args.TemplateName)
			os.Exit(1)
		}

		// Ask for confirmation
		fmt.Printf("Are you sure you want to remove template '%s'? This action cannot be undone. (y/N): ", args.TemplateName)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Template removal cancelled.")
			os.Exit(0)
		}

		// Remove the template directory
		if err := os.RemoveAll(templatePath); err != nil {
			fmt.Printf("Error removing template: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Template '%s' successfully removed from contributed templates.\n", args.TemplateName)
		return
	}

	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get project directory:", err)
		os.Exit(1)
	}

	if !args.IsCreate && !args.IsRemove {
		checkLocalFolder(projectDir, args) // may not return

	} else {
		wailsJsonPath := filepath.Join(projectDir, "wails.json")
		if _, err := os.Stat(wailsJsonPath); os.IsNotExist(err) {
			fmt.Println("Error: wails.json not found in the current directory.")
			fmt.Println("Create template mode requires a valid Wails project directory.")
			os.Exit(1)
		}
	}

	appConfig, configPath, err := config.LoadProjectConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}
	fmt.Println("CONFIG_PATH=", configPath)

	// Developer mode: Special case for jrush's development environment
	developerMode := false
	if strings.Contains(configPath, "jrush") {
		// Check if config file is missing or has empty values
		configMissing := false
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configMissing = true
		} else if appConfig.Organization == "" || appConfig.ProjectName == "" || appConfig.Github == "" || appConfig.Domain == "" {
			configMissing = true
		}

		if configMissing {
			// Set TrueBlocks defaults
			appConfig.Organization = "TrueBlocks, LLC"
			appConfig.ProjectName = "dalledress"
			appConfig.Github = "github.com/TrueBlocks/trueblocks-core"
			appConfig.Domain = "trueblocks.io"
			// Remove auto mode to allow interactive confirmation of defaults
			args.IsAuto = false
			developerMode = true
			// fmt.Println("Developer mode detected - using TrueBlocks defaults")
		}
	}

	if args.IsAuto {
		if appConfig.Organization == "" || appConfig.ProjectName == "" || appConfig.Github == "" || appConfig.Domain == "" {
			fmt.Println("Error: Auto mode requires default values in config file.")
			fmt.Println("Run without --auto first to create config file with defaults.")
			os.Exit(1)
		}
	}

	organization := appConfig.Organization
	projectName := appConfig.ProjectName
	github := appConfig.Github
	domain := appConfig.Domain
	chifra := "github.com/TrueBlocks/trueblocks-core/src/apps/chifra"

	// Interactive prompting for missing values (regular mode or --create mode or developer mode)
	// For --create mode, always prompt to allow project-specific configuration
	shouldPrompt := (!args.IsRemove && !args.IsAuto) &&
		((!args.IsCreate) || args.IsCreate || developerMode)

	if shouldPrompt {
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Organization [%s]: ", appConfig.Organization)
		organizationInput, _ := reader.ReadString('\n')
		organizationInput = strings.TrimSpace(organizationInput)
		if organizationInput != "" {
			organization = organizationInput
		}

		fmt.Printf("Project Name [%s]: ", appConfig.ProjectName)
		projectNameInput, _ := reader.ReadString('\n')
		projectNameInput = strings.TrimSpace(projectNameInput)
		if projectNameInput != "" {
			projectName = projectNameInput
		}

		fmt.Printf("Github [%s]: ", appConfig.Github)
		githubInput, _ := reader.ReadString('\n')
		githubInput = strings.TrimSpace(githubInput)
		if githubInput != "" {
			github = githubInput
		}

		fmt.Printf("Domain [%s]: ", appConfig.Domain)
		domainInput, _ := reader.ReadString('\n')
		domainInput = strings.TrimSpace(domainInput)
		if domainInput != "" {
			domain = domainInput
		}
	} else {
		if args.IsCreate {
			fmt.Println("Running in create template mode with default values:")
		} else {
			fmt.Println("Running in auto mode with default values:")
		}
		fmt.Println("Organization:", organization)
		fmt.Println("Project Name:", projectName)
		fmt.Println("Github:", github)
		fmt.Println("Domain:", domain)
	}

	// Validate required fields
	if organization == "" {
		fmt.Println("Error: Organization is required.")
		os.Exit(1)
	}
	if projectName == "" {
		fmt.Println("Error: Project Name is required.")
		os.Exit(1)
	}
	if github == "" {
		fmt.Println("Error: Github is required.")
		os.Exit(1)
	}
	if domain == "" {
		fmt.Println("Error: Domain is required.")
		os.Exit(1)
	}

	publisherName := "YourCompany"
	publisherEmail := "your_email@your_company.com"

	parts := strings.Split(organization, ",")
	orgName := strings.TrimSpace(parts[0])
	slug := strings.ToLower(orgName) + "-" + projectName

	// Determine which template to use and resolve template name for saving in config
	var resolvedTemplateName string

	if !args.IsCreate {
		// Check for --template flag first, then TEMPLATE_SOURCE environment variable, then saved config
		if args.UseTemplate != "" {
			resolvedTemplateName = args.UseTemplate
		} else if templateSource := os.Getenv("TEMPLATE_SOURCE"); templateSource != "" {
			// Only save template name if it resolves to a known template (not a full path)
			if _, err := templates.GetTemplateDir(templateSource); err == nil {
				resolvedTemplateName = templateSource
			}
		} else if appConfig.Template != "" {
			// Use template from saved configuration
			resolvedTemplateName = appConfig.Template
		} else {
			// Using default template
			resolvedTemplateName = "default"
		}
	} else {
		// In create mode, we're creating a template with the specified name
		resolvedTemplateName = args.TemplateName
	}

	// Save config if we prompted for values or if template was explicitly specified
	if shouldPrompt || resolvedTemplateName != "" {
		// Preserve existing config values and only update what has changed
		newConfig := &config.Config{
			Organization:  organization,
			ProjectName:   projectName,
			Github:        github,
			Domain:        domain,
			Template:      resolvedTemplateName,
			PreserveFiles: appConfig.PreserveFiles, // Preserve existing PreserveFiles
			ViewConfig:    appConfig.ViewConfig,    // Preserve existing ViewConfig
		}

		if args.IsCreate {
			// In create mode, save to project-local config
			if err := config.SaveProjectConfig(newConfig); err != nil {
				fmt.Println("Failed to save project config file:", err)
				os.Exit(1)
			}
		} else {
			// In regular mode, save to project-local config to establish project-specific settings
			// This ensures each project gets its own config file
			if err := config.SaveProjectConfig(newConfig); err != nil {
				fmt.Println("Failed to save project config file:", err)
				os.Exit(1)
			}
			// Also update global config for convenience as fallback defaults
			if err := config.SaveGlobalConfig(newConfig); err != nil {
				fmt.Println("Failed to save global config file:", err)
				os.Exit(1)
			}
		}
	}

	// Get template directory
	var templateDir string
	if args.IsCreate {
		// In create mode, we write to the contributed template directory
		configDir, err := config.GetUserConfigDir()
		if err != nil {
			fmt.Println("Failed to get user config directory:", err)
			os.Exit(1)
		}
		templateDir = filepath.Join(configDir, "templates", "contributed", args.TemplateName)
		fmt.Println("Creating template at:", templateDir)
	} else {
		var templateSource string
		if args.UseTemplate != "" {
			templateSource = args.UseTemplate
		} else if envTemplate := os.Getenv("TEMPLATE_SOURCE"); envTemplate != "" {
			templateSource = envTemplate
		} else if appConfig.Template != "" {
			templateSource = appConfig.Template
		}

		if templateSource != "" {
			// First try to resolve as a template name (contributed or system)
			if resolvedDir, err := templates.GetTemplateDir(templateSource); err == nil {
				templateDir = resolvedDir
				fmt.Printf("Using template '%s' from: %s\n", templateSource, templateDir)
			} else {
				// Fall back to treating it as a full path
				templateDir, err = filepath.Abs(templateSource)
				if err != nil {
					fmt.Printf("Failed to resolve template source '%s' as template name or path: %v\n", templateSource, err)
					os.Exit(1)
				}
				// Verify the path exists
				if _, err := os.Stat(templateDir); os.IsNotExist(err) {
					fmt.Printf("Template directory '%s' does not exist\n", templateDir)
					os.Exit(1)
				}
				fmt.Printf("Using custom template directory: %s\n", templateDir)
			}
		} else {
			// Use default system template
			templateDir, err = templates.GetDefaultTemplateDir()
			if err != nil {
				fmt.Println("Failed to get default template directory:", err)
				os.Exit(1)
			}
			fmt.Println("Using default template directory:", templateDir)
		}
	}

	fmt.Println("TEMPLATE_DIR: ", templateDir)
	fmt.Println("PROJECT_DIR:  ", projectDir)
	fmt.Println("ORGANIZATION: ", organization)
	fmt.Println("ORG_NAME:     ", orgName)
	fmt.Println("SLUG:         ", slug)
	fmt.Println("PROJECT_NAME: ", projectName)
	fmt.Println("GITHUB:       ", github)
	fmt.Println("DOMAIN:       ", domain)
	fmt.Println("CHIFRA:       ", chifra)

	// Create template variables with safety checks for empty strings
	projectProper := projectName
	if len(projectName) > 0 {
		projectProper = strings.ToUpper(projectName[0:1]) + projectName[1:]
	}

	templateVars := &processor.TemplateVars{
		ProjectName:    projectName,
		ProjectProper:  projectProper,
		PublisherName:  publisherName,
		PublisherEmail: publisherEmail,
		Organization:   organization,
		OrgName:        orgName,
		OrgLower:       strings.ToLower(orgName),
		Slug:           slug,
		Github:         github,
		Domain:         domain,
		Chifra:         chifra,
	}

	if !args.IsCreate && !args.IsRemove {
		err = filepath.Walk(templateDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(templateDir, path)
			targetPath := filepath.Join(projectDir, relPath)

			if info.IsDir() {
				return os.MkdirAll(targetPath, os.ModePerm)
			}

			if processor.ShouldPreserve(relPath, appConfig) {
				if _, err := os.Stat(targetPath); err == nil {
					fmt.Printf("Preserving existing file: %s\n", relPath)
					return nil
				}
			}

			input, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			content := processor.ApplyTemplateVars(string(input), templateVars)
			return os.WriteFile(targetPath, []byte(content), info.Mode())
		})
	} else {
		// Set environment variable to prevent macOS resource fork files
		originalCopyFile := os.Getenv("COPYFILE_DISABLE")
		os.Setenv("COPYFILE_DISABLE", "1")
		defer func() {
			if originalCopyFile == "" {
				os.Unsetenv("COPYFILE_DISABLE")
			} else {
				os.Setenv("COPYFILE_DISABLE", originalCopyFile)
			}
		}()

		fmt.Println("Create template mode: Updating template from current project")

		filesToCopy := make(map[string]bool)
		err = filepath.Walk(projectDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if yes, err := processor.IsExcluded(path, info); yes {
				// fmt.Printf("Skipping directory or file: %s\n", path)
				return err
			}

			relPath, _ := filepath.Rel(projectDir, path)
			if relPath != "" {
				filesToCopy[relPath] = true
			}

			return nil
		})

		if err != nil {
			fmt.Println("Error scanning current directory:", err)
			os.Exit(1)
		}

		fmt.Printf("Found %d files/directories in source\n", len(filesToCopy))

		// Only clean template directory if it already exists
		if _, err := os.Stat(templateDir); err == nil {
			err = filepath.Walk(templateDir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				relPath, _ := filepath.Rel(templateDir, path)
				if relPath == "" || relPath == "." || relPath == ".wails-template.json" {
					return nil
				}

				if !filesToCopy[relPath] {
					fmt.Printf("Removing file or folder from template: %s\n", relPath)
					os.Remove(path)
				}
				return nil
			})
			if err != nil {
				fmt.Println("Error cleaning template directory:", err)
				os.Exit(1)
			}
		}

		fmt.Println("Copying files to template with replacements...")
		err = filepath.Walk(projectDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if yes, err := processor.IsExcluded(path, info); yes {
				return err
			}

			relPath, err := filepath.Rel(projectDir, path)
			if err != nil {
				fmt.Printf("Error getting relative path for %s: %v\n", path, err)
				return nil
			}

			targetPath := filepath.Join(templateDir, relPath)

			targetDir := filepath.Dir(targetPath)
			if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
				fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
				return nil
			}

			// Skip reading and writing for directories
			if info.IsDir() {
				return nil
			}

			input, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil
			}

			content := processor.ReverseTemplateVars(string(input), templateVars)

			if err := os.WriteFile(targetPath, []byte(content), info.Mode()); err != nil {
				fmt.Printf("Error writing file %s: %v\n", targetPath, err)
			}

			return nil
		})
	}

	if err != nil {
		fmt.Println("Error processing files:", err)
		os.Exit(1)
	}

	if !args.IsCreate && !args.IsRemove && !args.IsAuto {
		fmt.Println("Running 'yarn install' in", projectDir)
		cmd := exec.Command("yarn", "install")
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: 'yarn install' failed: %v\n", err)
		}
		cmd = exec.Command("wails", "generate", "modules")
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: 'wails generate modules' failed: %v\n", err)
		}

		os.Remove(filepath.Join(projectDir, ".wails-template.json"))
		fmt.Println("✅ Project created at", projectDir)
		fmt.Println("✅ Next steps ==> Run:")
		fmt.Println()
		fmt.Println(colors.BrightBlue, "  yarn lint && yarn test && yarn start", colors.Off)
		fmt.Println()
	} else {
		fmt.Println("✅ Template updated from project at", projectDir)
	}
}

func checkLocalFolder(projectDir string, args *config.Args) {
	dirEntries, err := os.ReadDir(projectDir)
	if err != nil {
		fmt.Println("Failed to read project directory:", err)
		os.Exit(1)
	}

	checkFunc := func(e os.DirEntry) bool {
		if e.IsDir() && e.Name() == "frontend" {
			checkLocalFolder(filepath.Join(projectDir, "frontend"), args) // may not return
			return true
		}

		// If only these files are present, allow processing...
		okayFiles := []string{
			".git",
			".gitignore",
			"README.md",
			"LICENSE",
			".create-local-app.json",
			".env",
			"go.mod",
		}
		for _, okFile := range okayFiles {
			if e.Name() == okFile {
				// logger.InfoBB("File", e.Name(), "found. Okay")
				return true
			}
		}
		// logger.InfoBB("File", e.Name(), "found. Not okay")
		return false
	}

	dirEntries = slices.DeleteFunc(dirEntries, checkFunc)
	if len(dirEntries) > 0 && !args.IsAuto {
		if !args.IsForce {
			fmt.Println("The current directory (" + projectDir + ") contains files.")
			fmt.Println("Proceeding will overwrite existing files in an unrecoverable way.")
			fmt.Println("Use --force flag to proceed without this check.")
			os.Exit(1)
		}
	}
}
