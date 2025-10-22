package processor

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/TrueBlocks/create-local-app/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
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

	fileSkips := []string{
		".create-local-app.json",
		".DS_Store",
		".env",
		"shit",
		"Thumbs.db",
	}
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

// ShouldPreserve determines if an existing file should be preserved and not replaced
func ShouldPreserve(filePath string, cfg *config.Config) bool {
	if cfg == nil || len(cfg.PreserveFiles) == 0 {
		return false
	}

	// Only preserve files that actually exist
	if !file.FileExists(filePath) {
		return false
	}

	normalizedPath := filepath.ToSlash(filePath)
	for _, preserveFile := range cfg.PreserveFiles {
		if strings.HasSuffix(normalizedPath, preserveFile) {
			return true
		}
	}

	return false
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
		{"Block explorer", "{{101}}"},
		{"Block Explorer", "{{102}}"},
		{"class Explorer", "{{103}}"},
		{"dalleDressStore", "{{104}}"},
		{"Explorers", "{{105}}"},
		{"explorers: Explorer[];", "{{106}}"},
		{"ExportApprovals", "{{107}}"},
		{"ExportsOpenApprovals", "{{108}}"},
		{"getOpenapprovalsFields", "{{109}}"},
		{"getOpenApprovalsStore", "{{110}}"},
		{"json:\"dalledress\"", "{{111}}"},
		{"local block explorer", "{{112}}"},
		{"Local Explorer", "{{113}}"},
		{"localExplorer", "{{114}}"},
		{"LocalExplorer", "{{115}}"},
		{"new Explorer", "{{116}}"},
		{"Open Approvals", "{{117}}"},
		{"openapprovals", "{{118}}"},
		{"OpenApprovals", "{{119}}"},
		{"OPENAPPROVALS", "{{120}}"},
		{"openapprovalsFacet", "{{121}}"},
		{"openapprovalsStore", "{{122}}"},
		{"openapprovalsStoreMu", "{{123}}"},
		{"pageData?.dalledress", "{{124}}"},
		{"pageData.dalledress", "{{125}}"},
		{"remote block explorer", "{{126}}"},
		{"Remote Explorer", "{{127}}"},
		{"remoteExplorer", "{{128}}"},
		{"RemoteExplorer", "{{129}}"},
		{"SortOpenApprovals", "{{130}}"},
		{"this.explorers = this.convertValues(source[\"explorers\"], Explorer", "{{131}}"},
		{"TokensApprovals", "{{132}}"},
		{"dalledressStore", "{{133}}"},
		{"dalledressStoreMu", "{{134}}"},
		{"dresses-dalledress", "{{135}}"},
		{"\"dalledress\":", "{{136}}"},
		{"getDalledressFields", "{{137}}"},
		{"Store:         \"dalledress\"", "{{138}}"},
		{"this.dalledress = this.convertValues(source[\"dalledress\"], model.DalleDress);", "{{139}}"},
		{"dalledress: model.DalleDress[];", "{{140}}"},
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
