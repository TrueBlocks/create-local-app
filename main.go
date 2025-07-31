package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/TrueBlocks/create-local-app/pkg/config"
	"github.com/TrueBlocks/create-local-app/pkg/processor"
	"github.com/TrueBlocks/create-local-app/pkg/templates"
)

//go:embed templates/system/*.tar.gz
var systemTemplatesFS embed.FS

//go:embed VERSION
var versionContent string

func main() {
	version := strings.TrimSpace(versionContent)
	args, err := config.ParseArgs(version)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Initialize user configuration directory structure
	if err := config.InitializeUserConfig(); err != nil {
		fmt.Printf("Error initializing user config: %v\n", err)
		os.Exit(1)
	}

	// Initialize system templates on first run or if missing
	if err := templates.InitializeSystemTemplates(systemTemplatesFS); err != nil {
		fmt.Printf("Error initializing system templates: %v\n", err)
		os.Exit(1)
	}

	if args.IsReverse {
		err := templates.CreateTemplate(args.TemplateName)
		if err != nil {
			fmt.Printf("Error creating template: %v\n", err)
			os.Exit(1)
		}
		return
	}

	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get project directory:", err)
		os.Exit(1)
	}

	if !args.IsReverse {
		dirEntries, err := os.ReadDir(projectDir)
		if err != nil {
			fmt.Println("Failed to read project directory:", err)
			os.Exit(1)
		}
		// Remove ".git" if present (.git is okay since we don't replace it)
		dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool { return e.Name() == ".git" })
		if len(dirEntries) > 0 && !args.IsAuto {
			if !args.IsForce {
				fmt.Println("The current directory (" + projectDir + ") contains files.")
				fmt.Println("Proceeding will overwrite existing files in an unrecoverable way.")
				fmt.Println("Use --force flag to proceed without this check.")
				os.Exit(1)
			}
		}
	} else {
		wailsJsonPath := filepath.Join(projectDir, "wails.json")
		if _, err := os.Stat(wailsJsonPath); os.IsNotExist(err) {
			fmt.Println("Error: wails.json not found in the current directory.")
			fmt.Println("Reverse mode requires a valid Wails project directory.")
			os.Exit(1)
		}
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		fmt.Println("Failed to get config path:", err)
		os.Exit(1)
	}
	fmt.Println("CONFIG_PATH=", configPath)

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
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

	if !args.IsReverse && !args.IsAuto {
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
		if args.IsReverse {
			fmt.Println("Running in reverse mode with default values:")
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

	if !args.IsReverse && !args.IsAuto {
		newConfig := &config.Config{
			Organization: organization,
			ProjectName:  projectName,
			Github:       github,
			Domain:       domain,
		}
		if err := config.SaveConfig(configPath, newConfig); err != nil {
			fmt.Println("Failed to save config file:", err)
			os.Exit(1)
		}
	}

	// Get template directory
	var templateDir string
	customTemplateDir := os.Getenv("TEMPLATE_SOURCE")
	if customTemplateDir != "" {
		// Use custom template directory if specified
		templateDir, err = filepath.Abs(customTemplateDir)
		if err != nil {
			fmt.Println("Failed to resolve custom template directory:", err)
			os.Exit(1)
		}
		fmt.Println("Using custom template directory:", templateDir)
	} else {
		// Use default system template
		templateDir, err = templates.GetDefaultTemplateDir()
		if err != nil {
			fmt.Println("Failed to get default template directory:", err)
			os.Exit(1)
		}
		fmt.Println("Using default template directory:", templateDir)
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

	if !args.IsReverse {
		err = filepath.Walk(templateDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(templateDir, path)
			targetPath := filepath.Join(projectDir, relPath)

			if info.IsDir() {
				return os.MkdirAll(targetPath, os.ModePerm)
			}

			input, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			content := processor.ApplyTemplateVars(string(input), templateVars)
			return os.WriteFile(targetPath, []byte(content), info.Mode())
		})
	} else {
		fmt.Println("Reverse mode: Updating template from current project")

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

	if !args.IsReverse && !args.IsAuto {
		os.Remove(filepath.Join(projectDir, ".wails-template.json"))
		fmt.Println("✅ Project created at", projectDir)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  cd frontend && yarn install && cd ..")
		fmt.Println("  wails dev")
	} else {
		fmt.Println("✅ Template updated from project at", projectDir)
	}
}
