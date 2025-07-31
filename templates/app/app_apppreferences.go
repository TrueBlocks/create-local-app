package app

import (
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/preferences"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
)

// GetAppPreferences returns the current AppPreferences
func (a *App) GetAppPreferences() *preferences.AppPreferences {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	appPrefsCopy := a.Preferences.App
	return &appPrefsCopy
}

// SetAppPreferences sets the AppPreferences
func (a *App) SetAppPreferences(appPrefs *preferences.AppPreferences) error {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App = *appPrefs
	return preferences.SetAppPreferences(appPrefs)
}

// SetHelpCollapsed sets the help collapsed state
func (a *App) SetHelpCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.HelpCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetLanguage returns the current language setting
func (a *App) GetLanguage() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastLanguage
}

// SetLanguage sets the language setting
func (a *App) SetLanguage(language string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastLanguage = language
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetLastAddress returns the last active address
func (a *App) GetLastAddress() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastAddress
}

// SetLastAddress sets the last active address
func (a *App) SetLastAddress(address string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastAddress = address
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetLastChain returns the last chain setting
func (a *App) GetLastChain() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastChain
}

// SetLastChain sets the last chain setting
func (a *App) SetLastChain(chain string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastChain = chain
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetLastContract returns the last contract setting
func (a *App) GetLastContract() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastContract
}

// SetLastContract sets the last contract setting
func (a *App) SetLastContract(contract string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastContract = contract
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetLastTab returns the last tab for a given route
func (a *App) GetLastTab(route string) types.DataFacet {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return types.DataFacet(a.Preferences.App.LastTab[route])
}

// SetLastTab sets the last tab for a given route
func (a *App) SetLastTab(route string, tab types.DataFacet) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastTab[route] = string(tab)
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// SetLastView sets the last view and updates LastViewNoWizard if not wizard
func (a *App) SetLastView(view string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastView = view
	if view != "/wizard" {
		a.Preferences.App.LastViewNoWizard = view
	}
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// SetMenuCollapsed sets the menu collapsed state
func (a *App) SetMenuCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.MenuCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetTheme returns the current theme setting
func (a *App) GetTheme() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastTheme
}

// SetTheme sets the theme setting
func (a *App) SetTheme(theme string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastTheme = theme
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetDebugMode returns the current debug mode setting
func (a *App) GetDebugMode() bool {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.DebugMode
}

// SetDebugMode sets the debug mode setting
func (a *App) SetDebugMode(debugMode bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.DebugMode = debugMode
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}
