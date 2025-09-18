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

	folderSkips := []string{".git", "node_modules", "dist"}
	if slices.Contains(folderSkips, folderName) {
		return true, filepath.SkipDir
	}

	// Also check if the current file/folder name should be skipped
	if slices.Contains(folderSkips, baseName) {
		if info.IsDir() {
			return true, filepath.SkipDir
		} else {
			return true, nil
		}
	}

	fileSkips := []string{".DS_Store", "Thumbs.db", ".env", "shit", ".create-local-app.json"}
	if slices.Contains(fileSkips, baseName) {
		return true, nil
	}

	keeps := []string{"appicon.png", "Info.plist", "Info.dev.plist"}
	if strings.Contains(path, "/build/") && !slices.Contains(keeps, baseName) {
		return true, nil
	}

	if strings.Contains(path, "/ai") {
		keeps := []string{".gitignore", "README.md", "Invoker.md", "Rules.md"}
		if slices.Contains(keeps, baseName) {
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
	content = strings.ReplaceAll(content, "{{SDK}}", "github.com/TrueBlocks/trueblocks-sdk/v5")
	content = strings.ReplaceAll(content, "{{PACKAGES}}", "{{SAVEPKG}}")
	content = strings.ReplaceAll(content, "{{DALLE}}", "github.com/TrueBlocks/trueblocks-dalle/v2")
	content = strings.ReplaceAll(content, "{{APP}}", "github.com/TrueBlocks/"+vars.Slug+"/app")
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
	content = strings.ReplaceAll(content, "{{SAVEPKG}}", "github.com/TrueBlocks/"+vars.Slug+"/pkg")
	return content
}

// ReverseTemplateVars reverses template variable replacements (for create mode)
func ReverseTemplateVars(content string, vars *TemplateVars) string {
	preserves := []struct {
		old string
		new string
	}{
		{"remoteExplorer", "{{rE}}"},
		{"RemoteExplorer", "{{RE}}"},
		{"localExplorer", "{{lE}}"},
		{"LocalExplorer", "{{LE}}"},
		{"Remote Explorer", "{{RSE}}"},
		{"Local Explorer", "{{LSE}}"},
		{"remote block explorer", "{{rbe}}"},
		{"local block explorer", "{{lbe}}"},
		{"Block explorer", "{{Be}}"},
		{"Block Explorer", "{{BE}}"},
		{"Explorers", "{{E}}"},
		{"class Explorer", "{{cE}}"},
		{"new Explorer", "{{nE}}"},
		{"explorers: Explorer[];", "{{EE}}"},
		{"this.explorers = this.convertValues(source[\"explorers\"], Explorer", "{{EE2}}"},
	}
	for _, pr := range preserves {
		content = strings.ReplaceAll(content, pr.old, pr.new)
	}

	content = strings.ReplaceAll(content, "github.com/TrueBlocks/trueblocks-sdk/v5", "{{SDK}}")
	content = strings.ReplaceAll(content, "github.com/TrueBlocks/trueblocks-dalle/v2", "{{DALLE}}")
	content = strings.ReplaceAll(content, "github.com/TrueBlocks/"+vars.Slug+"/pkg", "{{PACKAGES}}")
	content = strings.ReplaceAll(content, "github.com/TrueBlocks/"+vars.Slug+"/app", "{{APP}}")
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

	for _, pr := range preserves {
		content = strings.ReplaceAll(content, pr.new, pr.old)
	}

	return content
}
