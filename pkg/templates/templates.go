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
	"github.com/TrueBlocks/create-local-app/pkg/processor"
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
func InitializeSystemTemplates(embeddedFS embed.FS) error {
	// Get user config directory for destination
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		return err
	}
	destTemplatesDir := filepath.Join(configDir, "templates", "system")

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

// CreateTemplate creates a template from the current directory in the contributed folder
func CreateTemplate(templateName string) error {
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

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Get user config directory
	configDir, err := config.GetUserConfigDir()
	if err != nil {
		return err
	}

	// Create the contributed template directory
	templateDir := filepath.Join(configDir, "templates", "contributed", templateName)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory %s: %v (check permissions)", templateDir, err)
	}

	// Check if template already exists
	if _, err := os.Stat(templateDir); err == nil {
		fmt.Printf("Warning: Template %s already exists and will be overwritten.\n", templateName)
		// Remove existing template
		if err := os.RemoveAll(templateDir); err != nil {
			return fmt.Errorf("failed to remove existing template: %v", err)
		}
		// Recreate directory
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			return fmt.Errorf("failed to recreate template directory: %v", err)
		}
	}

	// Copy current directory contents to template directory
	err = copyDirectory(currentDir, templateDir)
	if err != nil {
		return fmt.Errorf("failed to copy directory contents: %v", err)
	}

	fmt.Printf("Successfully created template: %s\n", templateDir)
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the destination directory itself to avoid infinite recursion
		if strings.HasPrefix(path, dst) {
			return filepath.SkipDir
		}

		// Skip excluded files and directories
		excluded, skipErr := processor.IsExcluded(path, info)
		if excluded {
			if skipErr == filepath.SkipDir && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return CopyFile(path, dstPath)
	})
}
