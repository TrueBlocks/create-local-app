package templates

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/create-local-app/pkg/config"
)

// GetTemplateDir returns the path to a template directory
// It looks first in ~/.create-local-app/templates/contributed/templateName, then ~/.create-local-app/templates/system/templateName
func GetTemplateDir(templateName string) (string, error) {
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		return "", err
	}

	// Check contributed templates first
	contributedPath := filepath.Join(configDir, "templates", "contributed", templateName)
	if info, err := os.Stat(contributedPath); err == nil && info.IsDir() {
		return contributedPath, nil
	}

	// Check system templates
	systemPath := filepath.Join(configDir, "templates", "system", templateName)
	if info, err := os.Stat(systemPath); err == nil && info.IsDir() {
		return systemPath, nil
	}

	return "", fmt.Errorf("template '%s' not found in contributed or system templates", templateName)
}

// GetDefaultTemplateDir returns the default system template directory
func GetDefaultTemplateDir() (string, error) {
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		return "", err
	}

	defaultPath := filepath.Join(configDir, "templates", "system", "default")
	if info, err := os.Stat(defaultPath); err != nil || !info.IsDir() {
		return "", fmt.Errorf("default template not found at %s - run initialization to set up templates", defaultPath)
	}

	return defaultPath, nil
}

// InitializeSystemTemplates extracts embedded system templates to the user config directory
// It checks versions to determine if templates need to be updated
func InitializeSystemTemplates(embeddedFS embed.FS, currentVersion string) error {
	// Get user config directory for destination
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		return err
	}
	destTemplatesDir := filepath.Join(configDir, "templates", "system")
	versionFilePath := filepath.Join(configDir, "VERSION")

	// Check if system templates need updating by comparing versions
	needsUpdate := true
	if _, err := os.Stat(destTemplatesDir); err == nil {
		// System templates directory exists, check if it has content
		entries, err := os.ReadDir(destTemplatesDir)
		if err == nil && len(entries) > 0 {
			// Templates exist, check existing version
			if existingVersionBytes, err := os.ReadFile(versionFilePath); err == nil {
				existingVersion := strings.TrimSpace(string(existingVersionBytes))
				if existingVersion == strings.TrimSpace(currentVersion) {
					// Versions match, no update needed
					needsUpdate = false
				}
			}
		}
	}

	// Always write/update the VERSION file in config directory
	if err := os.WriteFile(versionFilePath, []byte(currentVersion+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write VERSION file: %w", err)
	}

	// Skip update if not needed
	if !needsUpdate {
		return nil
	}

	// Remove existing system templates to ensure clean extraction
	if err := os.RemoveAll(destTemplatesDir); err != nil {
		return fmt.Errorf("failed to remove existing system templates at %s: %w (check permissions)", destTemplatesDir, err)
	}

	// Create the system templates directory
	if err := os.MkdirAll(destTemplatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create system templates directory %s: %w (check permissions)", destTemplatesDir, err)
	}

	// Extract each tar.gz file from the embedded filesystem
	entries, err := embeddedFS.ReadDir("templates/system")
	if err != nil {
		return fmt.Errorf("failed to read embedded templates directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tar.gz") {
			if err := extractTarGz(embeddedFS, "templates/system/"+entry.Name(), destTemplatesDir); err != nil {
				return fmt.Errorf("failed to extract %s: %w", entry.Name(), err)
			}
		}
	}

	fmt.Printf("Successfully initialized system templates at %s\n", destTemplatesDir)
	return nil
}

// extractTarGz extracts a tar.gz file from the embedded filesystem to the destination directory
func extractTarGz(embeddedFS embed.FS, tarPath, destDir string) error {
	// Open the embedded tar.gz file
	tarFile, err := embeddedFS.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open embedded tar file %s: %w", tarPath, err)
	}
	defer tarFile.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(tarFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Skip macOS resource fork files (._filename)
		if strings.Contains(header.Name, "/._") || strings.HasPrefix(filepath.Base(header.Name), "._") {
			continue
		}

		// Create the full destination path
		destPath := filepath.Join(destDir, header.Name)

		// Handle different file types
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		case tar.TypeReg:
			// Extract regular file
			if err := extractRegularFile(tarReader, destPath, header); err != nil {
				return fmt.Errorf("failed to extract file %s: %w", destPath, err)
			}
		default:
			// Skip other file types for now
			continue
		}
	}

	return nil
}

// extractRegularFile extracts a regular file from the tar reader
func extractRegularFile(tarReader *tar.Reader, destPath string, header *tar.Header) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Create the file
	outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer outFile.Close()

	// Copy file contents
	if _, err := io.Copy(outFile, tarReader); err != nil {
		return fmt.Errorf("failed to write file contents: %w", err)
	}

	return nil
}

// ListTemplates lists all available templates (system and contributed)
func ListTemplates() error {
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	fmt.Println("Available Templates:")
	fmt.Println()

	// List system templates
	systemDir := filepath.Join(configDir, "templates", "system")
	if systemTemplates, err := listTemplatesInDir(systemDir); err == nil && len(systemTemplates) > 0 {
		fmt.Println("System Templates:")
		for _, template := range systemTemplates {
			fmt.Printf("  %s\n", template)
		}
		fmt.Println()
	}

	// List contributed templates
	contributedDir := filepath.Join(configDir, "templates", "contributed")
	if contributedTemplates, err := listTemplatesInDir(contributedDir); err == nil && len(contributedTemplates) > 0 {
		fmt.Println("Contributed Templates:")
		for _, template := range contributedTemplates {
			fmt.Printf("  %s\n", template)
		}
		fmt.Println()
	}

	return nil
}

// listTemplatesInDir lists templates in a specific directory
func listTemplatesInDir(dir string) ([]string, error) {
	var templates []string

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return templates, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}

	return templates, nil
}
