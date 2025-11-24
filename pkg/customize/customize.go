package customize

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/TrueBlocks/create-local-app/pkg/config"
	"github.com/TrueBlocks/goMaker/v6/types"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/colors"
)

// RunCustomize executes the interactive customize workflow
func RunCustomize() error {
	// Step 9: Validate .create-local-app.json exists
	if _, err := os.Stat(".create-local-app.json"); os.IsNotExist(err) {
		return fmt.Errorf(".create-local-app.json file not found in current directory")
	}

	// Step 10: Call goMaker's ReadTomlFiles(includeDisabled=false)
	structures, err := types.ReadTomlFiles(false)
	if err != nil {
		return fmt.Errorf("failed to read TOML files: %w", err)
	}

	// Step 11: Validate ./code_gen/templates/classDefinitions exists with enabled TOML
	if len(structures) == 0 {
		return fmt.Errorf("no enabled structures found in ./code_gen/templates/classDefinitions")
	}

	// Load existing configuration
	existingConfig, err := config.LoadConfig(".create-local-app.json")
	if err != nil {
		return fmt.Errorf("failed to load existing config: %w", err)
	}

	// Step 12: Implement ViewConfig merging with decade menuOrder
	updatedConfig, err := mergeViewConfig(existingConfig, structures)
	if err != nil {
		return fmt.Errorf("failed to merge ViewConfig: %w", err)
	}

	// Step 13-15: Show opening summary and enter command loop
	fmt.Println("\n=== Customize View Configuration ===")
	fmt.Println("Configure which views are enabled/disabled in your application.")
	fmt.Println("Current configuration:")
	fmt.Println()

	// Display current configuration in table format
	displayConfigTable(updatedConfig)

	// Enter command loop
	finalConfig, hasChanges, err := runCommandLoop(updatedConfig)
	if err != nil {
		return fmt.Errorf("command loop failed: %w", err)
	}

	// If no changes were made, exit early
	if !hasChanges {
		return nil
	}

	// Step 16-17: Summary and confirmation
	if err := showSummaryAndConfirm(finalConfig); err != nil {
		return fmt.Errorf("summary and confirmation failed: %w", err)
	}

	// Step 18: Use existing config schema/save logic from create-local-app
	if err := config.SaveConfig(".create-local-app.json", finalConfig); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Configuration updated successfully!")
	return nil
}

// mergeViewConfig merges existing ViewConfig with new structures, assigning decade menuOrders to new items
func mergeViewConfig(existingConfig *config.Config, structures []types.Structure) (*config.Config, error) {
	if existingConfig.ViewConfig == nil {
		existingConfig.ViewConfig = make(map[string]config.ViewConfigEntry)
	}

	// Find the highest existing menuOrder
	highestOrder := 0
	for _, entry := range existingConfig.ViewConfig {
		if entry.MenuOrder > highestOrder {
			highestOrder = entry.MenuOrder
		}
	}

	// Ensure we start at the next decade
	nextOrder := ((highestOrder / 10) + 1) * 10

	// Add new structures that don't exist in ViewConfig
	for _, structure := range structures {
		structName := strings.ToLower(structure.Class)
		if _, exists := existingConfig.ViewConfig[structName]; !exists {
			// Special handling for 'projects' - always menuOrder 0 and enabled
			if structName == "projects" {
				existingConfig.ViewConfig[structName] = config.ViewConfigEntry{
					MenuOrder:      0,
					Disabled:       false,
					DisabledFacets: make(map[string]bool),
				}
			} else {
				existingConfig.ViewConfig[structName] = config.ViewConfigEntry{
					MenuOrder:      nextOrder,
					Disabled:       false,
					DisabledFacets: make(map[string]bool),
				}
				nextOrder += 10 // Next decade
			}
		} else if structName == "projects" {
			// Ensure projects always has menuOrder 0 and is enabled
			entry := existingConfig.ViewConfig[structName]
			entry.MenuOrder = 0
			entry.Disabled = false
			existingConfig.ViewConfig[structName] = entry
		}
	}

	return existingConfig, nil
}

// displayConfigTable shows the current configuration in a formatted ASCII table
func displayConfigTable(cfg *config.Config) {
	// Create sorted list by menuOrder
	var items []struct {
		name  string
		entry config.ViewConfigEntry
	}

	for name, entry := range cfg.ViewConfig {
		// Skip 'projects' - it's not user-configurable
		if name == "projects" {
			continue
		}
		items = append(items, struct {
			name  string
			entry config.ViewConfigEntry
		}{name, entry})
	}

	// Sort by menuOrder
	sort.Slice(items, func(i, j int) bool {
		return items[i].entry.MenuOrder < items[j].entry.MenuOrder
	})

	if len(items) == 0 {
		fmt.Println("No views configured.")
		return
	}

	// Calculate column widths
	maxKeyLen := 3 // "Key"
	for _, item := range items {
		if len(item.name) > maxKeyLen {
			maxKeyLen = len(item.name)
		}
	}

	disabledWidth := 8   // "Disabled"
	menuOrderWidth := 10 // "Menu Order"

	// Print table header
	fmt.Printf("┌─%s─┬─%s─┬─%s─┐\n",
		strings.Repeat("─", maxKeyLen),
		strings.Repeat("─", disabledWidth),
		strings.Repeat("─", menuOrderWidth))

	fmt.Printf("│ %-*s │ %-*s │ %-*s │\n",
		maxKeyLen, "Key",
		disabledWidth, "Disabled",
		menuOrderWidth, "Menu Order")

	fmt.Printf("├─%s─┼─%s─┼─%s─┤\n",
		strings.Repeat("─", maxKeyLen),
		strings.Repeat("─", disabledWidth),
		strings.Repeat("─", menuOrderWidth))

	// Print table rows
	for _, item := range items {
		// Handle disabled column - center + 1 for the ✗
		var paddedDisabledDisplay string
		if item.entry.Disabled {
			leftPadding := (disabledWidth-1)/2 + 1 // Center + 1 for the ✗
			rightPadding := disabledWidth - 1 - leftPadding
			paddedDisabledDisplay = strings.Repeat(" ", leftPadding) + colors.Red + "✗" + colors.Off + strings.Repeat(" ", rightPadding)
		} else {
			paddedDisabledDisplay = strings.Repeat(" ", disabledWidth)
		}

		// Handle menu order column - right align with 3 spaces offset from center
		menuOrderStr := fmt.Sprintf("%d", item.entry.MenuOrder)
		menuOrderCenter := menuOrderWidth / 2
		menuOrderPos := menuOrderCenter + 2 // 2 spaces right of center
		menuOrderLeftPadding := menuOrderPos - len(menuOrderStr)
		if menuOrderLeftPadding < 0 {
			menuOrderLeftPadding = 0
		}
		menuOrderRightPadding := menuOrderWidth - len(menuOrderStr) - menuOrderLeftPadding
		paddedMenuOrder := strings.Repeat(" ", menuOrderLeftPadding) + menuOrderStr + strings.Repeat(" ", menuOrderRightPadding)

		fmt.Printf("│ %-*s │ %s │ %s │\n",
			maxKeyLen, item.name,
			paddedDisabledDisplay,
			paddedMenuOrder)
	}

	fmt.Printf("└─%s─┴─%s─┴─%s─┘\n",
		strings.Repeat("─", maxKeyLen),
		strings.Repeat("─", disabledWidth),
		strings.Repeat("─", menuOrderWidth))
}

// showSummaryAndConfirm displays final configuration and asks for confirmation
func showSummaryAndConfirm(cfg *config.Config) error {
	fmt.Println("\n=== Configuration Summary ===")

	displayConfigTable(cfg)

	// Renumber menuOrders to eliminate gaps (optional step from risk notes)
	renumberMenuOrders(cfg)

	fmt.Println("\nSave these changes? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	if strings.TrimSpace(strings.ToLower(response)) != "y" {
		return fmt.Errorf("configuration changes cancelled by user")
	}

	return nil
}

// renumberMenuOrders eliminates gaps in menuOrder numbering
func renumberMenuOrders(cfg *config.Config) {
	// Create sorted list by current menuOrder
	var items []struct {
		name  string
		entry config.ViewConfigEntry
	}

	for name, entry := range cfg.ViewConfig {
		items = append(items, struct {
			name  string
			entry config.ViewConfigEntry
		}{name, entry})
	}

	// Sort by current menuOrder
	sort.Slice(items, func(i, j int) bool {
		return items[i].entry.MenuOrder < items[j].entry.MenuOrder
	})

	// Renumber starting from 10, incrementing by 10
	for i, item := range items {
		newOrder := (i + 1) * 10
		entry := item.entry
		entry.MenuOrder = newOrder
		cfg.ViewConfig[item.name] = entry
	}
}
