package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

//go:embed templates.tar.gz
var templatesData []byte

// extractTemplates extracts the embedded templates.tar.gz to a temporary directory
func extractTemplates() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "create-local-app-templates-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Create gzip reader
	gzipReader, err := gzip.NewReader(strings.NewReader(string(templatesData)))
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("failed to read tar header: %w", err)
		}

		targetPath := filepath.Join(tempDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				os.RemoveAll(tempDir)
				return "", fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			// Ensure the directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				os.RemoveAll(tempDir)
				return "", fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}

			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				os.RemoveAll(tempDir)
				return "", fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				os.RemoveAll(tempDir)
				return "", fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			file.Close()
		}
	}

	return tempDir, nil
}

func main() {
	autoMode, reverseMode, err := parseArgs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get project directory:", err)
		os.Exit(1)
	}

	if !reverseMode {
		dirEntries, err := os.ReadDir(projectDir)
		if err != nil {
			fmt.Println("Failed to read project directory:", err)
			os.Exit(1)
		}
		if len(dirEntries) > 0 && !autoMode {
			fmt.Println("The current directory (" + projectDir + ") contains files.")
			fmt.Println("Proceeding will overwrite existing files in an unrecoverable way.")
			fmt.Print("Are you sure you want to proceed? (Y/n): ")
			reader := bufio.NewReader(os.Stdin)
			confirmation, _ := reader.ReadString('\n')
			confirmation = strings.TrimSpace(confirmation)
			if confirmation != "" && confirmation != "y" && confirmation != "Y" && confirmation != "yes" && confirmation != "Yes" {
				fmt.Println("Operation cancelled.")
				os.Exit(0)
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

	exePath, err := os.Executable()
	execDir := filepath.Dir(exePath)
	if err != nil {
		fmt.Println("Failed to get executable path:", err)
		os.Exit(1)
	}
	configPath := filepath.Join(execDir, ".wails-template.json")
	fmt.Println("CONFIG_PATH=", configPath)

	defaultOrg := ""
	defaultProject := ""
	defaultGithub := ""
	defaultDomain := ""
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var config map[string]string
			if json.Unmarshal(data, &config) == nil {
				defaultOrg = config["Organization"]
				defaultProject = config["ProjectName"]
				defaultGithub = config["Github"]
				defaultDomain = config["Domain"]
			}
		}
	}

	if autoMode {
		if defaultOrg == "" || defaultProject == "" || defaultGithub == "" || defaultDomain == "" {
			fmt.Println("Error: Auto mode requires default values in config file.")
			fmt.Println("Run without --auto first to create config file with defaults.")
			os.Exit(1)
		}
	}

	organization := defaultOrg
	projectName := defaultProject
	github := defaultGithub
	domain := defaultDomain
	chifra := "github.com/TrueBlocks/trueblocks-core/src/apps/chifra"

	if !reverseMode && !autoMode {
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Organization [%s]: ", defaultOrg)
		organizationInput, _ := reader.ReadString('\n')
		organizationInput = strings.TrimSpace(organizationInput)
		if organizationInput != "" {
			organization = organizationInput
		}

		fmt.Printf("Project Name [%s]: ", defaultProject)
		projectNameInput, _ := reader.ReadString('\n')
		projectNameInput = strings.TrimSpace(projectNameInput)
		if projectNameInput != "" {
			projectName = projectNameInput
		}

		fmt.Printf("Github [%s]: ", defaultGithub)
		githubInput, _ := reader.ReadString('\n')
		githubInput = strings.TrimSpace(githubInput)
		if githubInput != "" {
			github = githubInput
		}

		fmt.Printf("Domain [%s]: ", defaultDomain)
		domainInput, _ := reader.ReadString('\n')
		domainInput = strings.TrimSpace(domainInput)
		if domainInput != "" {
			domain = domainInput
		}
	} else {
		if reverseMode {
			fmt.Println("Running in reverse mode with default values:")
		} else {
			fmt.Println("Running in auto mode with default values:")
		}
		fmt.Println("Organization:", organization)
		fmt.Println("Project Name:", projectName)
		fmt.Println("Github:", github)
		fmt.Println("Domain:", domain)
	}

	publisherName := "YourCompany"
	publisherEmail := "your_email@your_company.com"

	parts := strings.Split(organization, ",")
	orgName := strings.TrimSpace(parts[0])
	slug := strings.ToLower(orgName) + "-" + projectName

	if !reverseMode && !autoMode {
		config := map[string]string{
			"Organization": organization,
			"ProjectName":  projectName,
			"Github":       github,
			"Domain":       domain,
		}
		configData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("Failed to marshal config to JSON:", err)
			os.Exit(1)
		}
		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			fmt.Println("Failed to write config file:", err)
			os.Exit(1)
		}
	}

	// Extract embedded templates to temporary directory
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
		// Extract embedded templates
		templateDir, err = extractTemplates()
		if err != nil {
			fmt.Println("Failed to extract embedded templates:", err)
			os.Exit(1)
		}
		// Clean up temporary directory when done
		defer os.RemoveAll(templateDir)
		fmt.Println("Using embedded templates (extracted to temp dir)")
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

	if !reverseMode {
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

			content := string(input)
			content = strings.ReplaceAll(content, "{{PROJECT_NAME}}", projectName)
			content = strings.ReplaceAll(content, "{{PROJECT_PROPER}}", strings.ToUpper(projectName[0:1])+projectName[1:])
			content = strings.ReplaceAll(content, "{{PUBLISHER_NAME}}", publisherName)
			content = strings.ReplaceAll(content, "{{PUBLISHER_EMAIL}}", publisherEmail)
			content = strings.ReplaceAll(content, "{{ORGANIZATION}}", organization)
			content = strings.ReplaceAll(content, "{{ORG_NAME}}", orgName)
			content = strings.ReplaceAll(content, "{{ORG_LOWER}}", strings.ToLower(orgName))
			content = strings.ReplaceAll(content, "{{SLUG}}", slug)
			content = strings.ReplaceAll(content, "{{GITHUB}}", github)
			content = strings.ReplaceAll(content, "{{DOMAIN}}", domain)
			content = strings.ReplaceAll(content, "{{CHIFRA}}", chifra)

			return os.WriteFile(targetPath, []byte(content), info.Mode())
		})
	} else {
		fmt.Println("Reverse mode: Updating template from current project")

		filesToCopy := make(map[string]bool)
		err = filepath.Walk(projectDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if yes, err := isExcluded(path, info); yes {
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

			if yes, err := isExcluded(path, info); yes {
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

			content := string(input)

			content = strings.ReplaceAll(content, chifra, "{{CHIFRA}}")
			if domain != "" {
				content = strings.ReplaceAll(content, domain, "{{DOMAIN}}")
			}
			if github != "" {
				content = strings.ReplaceAll(content, github, "{{GITHUB}}")
			}
			if slug != "" {
				content = strings.ReplaceAll(content, slug, "{{SLUG}}")
			}
			if orgName != "" {
				content = strings.ReplaceAll(content, orgName, "{{ORG_NAME}}")
				content = strings.ReplaceAll(content, strings.ToLower(orgName), "{{ORG_LOWER}}")
			}
			if organization != "" {
				content = strings.ReplaceAll(content, organization, "{{ORGANIZATION}}")
			}
			if projectName != "" {
				content = strings.ReplaceAll(content, projectName, "{{PROJECT_NAME}}")
				proper := strings.ToUpper(projectName[0:1]) + projectName[1:]
				content = strings.ReplaceAll(content, proper, "{{PROJECT_PROPER}}")
			}

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

	if !reverseMode && !autoMode {
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

func parseArgs() (isAuto bool, isReverse bool, err error) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--reverse":
			return false, true, nil
		case "--auto":
			return true, false, nil
		default:
			return false, false, fmt.Errorf("unknown argument: %s (valid options: --reverse, --auto)", os.Args[1])
		}
	}
	return false, false, nil
}

// Update the isExcluded function to accept path and FileInfo
func isExcluded(path string, info fs.FileInfo) (bool, error) {
	baseName := filepath.Base(path)
	folderName := filepath.Base(filepath.Dir(path))

	if folderName == ".git" || folderName == "node_modules" || folderName == "dist" {
		return true, filepath.SkipDir
	}

	if baseName == ".env" || baseName == ".DS_Store" || baseName == "shit" {
		return true, nil
	}

	if strings.Contains(path, "/build") && baseName != "appicon.png" {
		return true, nil
	}

	if strings.Contains(path, "/ai") {
		keep := []string{"README.md", "RulesOfEngagement.md", ".gitignore", "ToDoList.md"}
		if slices.Contains(keep, baseName) {
			return false, nil
		}
		return baseName != "ai", nil
	}

	if strings.Contains(path, "/output") {
		return baseName != ".gitignore", nil
	}

	if strings.Contains(path, "/book/book") && baseName != "book.toml" {
		return true, nil
	}

	return false, nil
}
