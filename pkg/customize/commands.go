package customize

import (
	"fmt"
	"io"
	"strings"

	"github.com/TrueBlocks/create-local-app/pkg/config"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/colors"
	"github.com/chzyer/readline"
)

// runCommandLoop handles the interactive command interface
func runCommandLoop(cfg *config.Config) (*config.Config, bool, error) {
	hasChanges := false

	showHelp()

	// Initialize readline
	rl, err := readline.New("customize> ")
	if err != nil {
		return nil, false, fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			if err == io.EOF || err == readline.ErrInterrupt {
				fmt.Println("\nExiting customize mode...")
				return cfg, hasChanges, nil
			}
			return nil, false, fmt.Errorf("failed to read user input: %w", err)
		}

		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "order":
			if len(parts) < 3 {
				fmt.Printf("%sError: order command requires view name and direction(s)%s\n", colors.Red, colors.Off)
				continue
			}
			changed, err := handleOrderCommand(cfg, parts[1:])
			if err != nil {
				fmt.Printf("%sError: %s%s\n", colors.Red, err.Error(), colors.Off)
				continue
			}
			if changed {
				hasChanges = true
				displayConfigTable(cfg)
			}

		case "disable", "enable":
			if len(parts) < 2 {
				fmt.Printf("%sError: %s command requires view name(s)%s\n", colors.Red, command, colors.Off)
				continue
			}
			changed, err := handleDisableCommand(cfg, command, parts[1:])
			if err != nil {
				fmt.Printf("%sError: %s%s\n", colors.Red, err.Error(), colors.Off)
				continue
			}
			if changed {
				hasChanges = true
				displayConfigTable(cfg)
			}

		case "summary":
			fmt.Printf("Command: %s\n", command)
			displayConfigTable(cfg)

		case "help", "h":
			showHelp()

		case "q", "quit":
			fmt.Println("Exiting customize mode...")
			return cfg, hasChanges, nil

		case "":
			// Empty input, continue loop
			continue

		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Available commands: disable, enable, order, summary, help, q/quit")
		}
	}
}

// handleDisableCommand processes disable/enable commands with view names
func handleDisableCommand(cfg *config.Config, command string, args []string) (bool, error) {
	// Join all args and split by comma for comma-separated support
	allArgs := strings.Join(args, " ")
	viewNames := strings.Split(allArgs, ",")

	// Trim whitespace from each view name
	for i, name := range viewNames {
		viewNames[i] = strings.TrimSpace(name)
	}

	// Remove empty entries
	var cleanNames []string
	for _, name := range viewNames {
		if name != "" {
			cleanNames = append(cleanNames, name)
		}
	}

	if len(cleanNames) == 0 {
		return false, fmt.Errorf("%s command requires view name(s)", command)
	}

	// Check for "enable all" command with exact case
	if len(cleanNames) == 1 && cleanNames[0] == "all" {
		if command == "enable" {
			return handleEnableAllCommand(cfg)
		} else {
			return false, fmt.Errorf("'disable all' is not supported - specify individual view names")
		}
	}

	// Validate all view names exist and check for 'projects'
	var invalidNames []string
	var protectedNames []string
	for _, name := range cleanNames {
		if name == "projects" {
			protectedNames = append(protectedNames, name)
		} else if _, exists := cfg.ViewConfig[name]; !exists {
			invalidNames = append(invalidNames, name)
		}
	}

	if len(protectedNames) > 0 {
		return false, fmt.Errorf("'%s' cannot be modified - it is always enabled", strings.Join(protectedNames, ", "))
	}

	if len(invalidNames) > 0 {
		return false, fmt.Errorf("invalid view name(s): %s", strings.Join(invalidNames, ", "))
	}

	// For disable operations, validate at least one view will remain enabled
	if command == "disable" {
		if err := validateAtLeastOneEnabled(cfg, cleanNames); err != nil {
			return false, err
		}
	}

	// Apply changes
	disabled := (command == "disable")
	var changedViews []string

	for _, name := range cleanNames {
		entry := cfg.ViewConfig[name]
		if entry.Disabled != disabled {
			entry.Disabled = disabled
			cfg.ViewConfig[name] = entry
			changedViews = append(changedViews, name)
		}
	}

	if len(changedViews) == 0 {
		fmt.Printf("No changes made - all views already %sd\n", command)
		return false, nil
	}

	// Generate summary message
	action := "disabled"
	if command == "enable" {
		action = "enabled"
	}

	if len(changedViews) == 1 {
		fmt.Printf("%s %s\n", changedViews[0], action)
	} else {
		fmt.Printf("%d views %s: %s\n", len(changedViews), action, strings.Join(changedViews, ", "))
	}

	return true, nil
}

// handleEnableAllCommand processes "enable all" commands
func handleEnableAllCommand(cfg *config.Config) (bool, error) {
	// Enable all views
	var changedViews []string
	for name, entry := range cfg.ViewConfig {
		if entry.Disabled {
			entry.Disabled = false
			cfg.ViewConfig[name] = entry
			changedViews = append(changedViews, name)
		}
	}

	if len(changedViews) == 0 {
		fmt.Println("No changes made - all views already enabled")
		return false, nil
	}

	fmt.Printf("All views enabled (%d views changed)\n", len(changedViews))
	return true, nil
}

// validateAtLeastOneEnabled ensures at least one view will remain enabled after disabling
func validateAtLeastOneEnabled(cfg *config.Config, viewsToDisable []string) error {
	enabledCount := 0
	for _, entry := range cfg.ViewConfig {
		if !entry.Disabled {
			enabledCount++
		}
	}

	// Count how many currently enabled views would be disabled
	disableCount := 0
	for _, name := range viewsToDisable {
		if entry, exists := cfg.ViewConfig[name]; exists && !entry.Disabled {
			disableCount++
		}
	}

	if enabledCount-disableCount < 1 {
		return fmt.Errorf("cannot disable all views - at least one must remain enabled")
	}

	return nil
}

// showHelp displays the available commands and their descriptions
func showHelp() {
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  disable <view>     - Disable specific view(s) (comma-separated)")
	fmt.Println("  enable <view>      - Enable specific view(s) (comma-separated)")
	fmt.Println("  enable all         - Enable all views")
	fmt.Println("  order <view> <dir> - Reorder view (dir: up/down, can repeat)")
	fmt.Println("  summary            - Show current configuration table")
	fmt.Println("  h / help           - Show this help message")
	fmt.Println("  q / quit           - Exit customize mode")
	fmt.Println()
}

// handleOrderCommand processes order commands with view name and direction sequence
func handleOrderCommand(cfg *config.Config, args []string) (bool, error) {
	if len(args) < 2 {
		return false, fmt.Errorf("order command requires view name and at least one direction")
	}

	viewName := strings.TrimSpace(args[0])
	directions := args[1:]

	// Trim all directions
	for i, dir := range directions {
		directions[i] = strings.TrimSpace(dir)
	}

	// Validate view name exists and is not 'projects'
	if viewName == "projects" {
		return false, fmt.Errorf("'projects' cannot be reordered - it always appears first")
	}
	if _, exists := cfg.ViewConfig[viewName]; !exists {
		return false, fmt.Errorf("invalid view name: %s", viewName)
	}

	// Validate all directions are "up" or "down"
	for _, dir := range directions {
		if dir != "up" && dir != "down" {
			return false, fmt.Errorf("invalid direction: %s. Use 'up' or 'down'", dir)
		}
	}

	// Get current sorted positions
	positions := getSortedPositions(cfg)

	// Find current position of the view
	currentPos := -1
	for i, pos := range positions {
		if pos.name == viewName {
			currentPos = i
			break
		}
	}

	if currentPos == -1 {
		return false, fmt.Errorf("view %s not found in current positions", viewName)
	}

	// Simulate the entire move sequence to validate bounds
	simulatedPos := currentPos
	for i, dir := range directions {
		if dir == "up" {
			if simulatedPos <= 0 {
				return false, fmt.Errorf("cannot move %s up - move %d would exceed top boundary", viewName, i+1)
			}
			simulatedPos--
		} else { // dir == "down"
			if simulatedPos >= len(positions)-1 {
				return false, fmt.Errorf("cannot move %s down - move %d would exceed bottom boundary", viewName, i+1)
			}
			simulatedPos++
		}
	}

	// All moves are valid, apply them
	startPos := currentPos
	workingPos := currentPos

	for _, dir := range directions {
		if dir == "up" {
			// Swap with previous item
			positions[workingPos], positions[workingPos-1] = positions[workingPos-1], positions[workingPos]
			workingPos--
		} else { // dir == "down"
			// Swap with next item
			positions[workingPos], positions[workingPos+1] = positions[workingPos+1], positions[workingPos]
			workingPos++
		}
	}

	// Update the config with new positions
	updateConfigFromPositions(cfg, positions)

	// Renumber using decades
	renumberMenuOrders(cfg)

	// Generate summary message
	finalPos := workingPos
	moveCount := 0
	direction := ""

	if finalPos < startPos {
		direction = "up"
		moveCount = startPos - finalPos
	} else if finalPos > startPos {
		direction = "down"
		moveCount = finalPos - startPos
	}

	if moveCount == 0 {
		fmt.Printf("No changes made - %s position unchanged\n", viewName)
		return false, nil
	}

	fmt.Printf("%s moved %s %d position(s)\n", viewName, direction, moveCount)
	return true, nil
}

// position represents a view's name and menu order for sorting
type position struct {
	name      string
	menuOrder int
}

// getSortedPositions returns views sorted by their current menuOrder
func getSortedPositions(cfg *config.Config) []position {
	var positions []position
	for name, entry := range cfg.ViewConfig {
		// Skip 'projects' - it's not part of reordering
		if name == "projects" {
			continue
		}
		positions = append(positions, position{
			name:      name,
			menuOrder: entry.MenuOrder,
		})
	}

	// Sort by menuOrder
	for i := 0; i < len(positions)-1; i++ {
		for j := i + 1; j < len(positions); j++ {
			if positions[i].menuOrder > positions[j].menuOrder {
				positions[i], positions[j] = positions[j], positions[i]
			}
		}
	}

	return positions
}

// updateConfigFromPositions updates the config with new position order (temporary menuOrders)
func updateConfigFromPositions(cfg *config.Config, positions []position) {
	for i, pos := range positions {
		// Set temporary order for renumberMenuOrders to work with
		entry := cfg.ViewConfig[pos.name]
		entry.MenuOrder = (i + 1) * 10
		cfg.ViewConfig[pos.name] = entry
	}
}
