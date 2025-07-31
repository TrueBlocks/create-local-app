package processor

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

// IsExcluded determines if a file or directory should be excluded from processing
func IsExcluded(path string, info fs.FileInfo) (bool, error) {
	_ = info // linter
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

// TemplateVars represents template replacement variables
type TemplateVars struct {
	ProjectName    string
	ProjectProper  string
	PublisherName  string
	PublisherEmail string
	Organization   string
	OrgName        string
	OrgLower       string
	Slug           string
	Github         string
	Domain         string
	Chifra         string
}

// ApplyTemplateVars applies template variable replacements to content
func ApplyTemplateVars(content string, vars *TemplateVars) string {
	content = strings.ReplaceAll(content, "{{PROJECT_NAME}}", vars.ProjectName)
	content = strings.ReplaceAll(content, "{{PROJECT_PROPER}}", vars.ProjectProper)
	content = strings.ReplaceAll(content, "{{PUBLISHER_NAME}}", vars.PublisherName)
	content = strings.ReplaceAll(content, "{{PUBLISHER_EMAIL}}", vars.PublisherEmail)
	content = strings.ReplaceAll(content, "{{ORGANIZATION}}", vars.Organization)
	content = strings.ReplaceAll(content, "{{ORG_NAME}}", vars.OrgName)
	content = strings.ReplaceAll(content, "{{ORG_LOWER}}", vars.OrgLower)
	content = strings.ReplaceAll(content, "{{SLUG}}", vars.Slug)
	content = strings.ReplaceAll(content, "{{GITHUB}}", vars.Github)
	content = strings.ReplaceAll(content, "{{DOMAIN}}", vars.Domain)
	content = strings.ReplaceAll(content, "{{CHIFRA}}", vars.Chifra)
	return content
}

// ReverseTemplateVars reverses template variable replacements (for create mode)
func ReverseTemplateVars(content string, vars *TemplateVars) string {
	content = strings.ReplaceAll(content, vars.Chifra, "{{CHIFRA}}")
	if vars.Domain != "" {
		content = strings.ReplaceAll(content, vars.Domain, "{{DOMAIN}}")
	}
	if vars.Github != "" {
		content = strings.ReplaceAll(content, vars.Github, "{{GITHUB}}")
	}
	if vars.Slug != "" {
		content = strings.ReplaceAll(content, vars.Slug, "{{SLUG}}")
	}
	if vars.OrgName != "" {
		content = strings.ReplaceAll(content, vars.OrgName, "{{ORG_NAME}}")
		content = strings.ReplaceAll(content, vars.OrgLower, "{{ORG_LOWER}}")
	}
	if vars.Organization != "" {
		content = strings.ReplaceAll(content, vars.Organization, "{{ORGANIZATION}}")
	}
	if vars.ProjectName != "" {
		content = strings.ReplaceAll(content, vars.ProjectName, "{{PROJECT_NAME}}")
		content = strings.ReplaceAll(content, vars.ProjectProper, "{{PROJECT_PROPER}}")
	}
	return content
}
