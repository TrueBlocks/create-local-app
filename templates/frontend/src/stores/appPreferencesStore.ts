import * as App from '@app';
import { preferences, types } from '@models';

// State interface matching the current useActiveProject state
interface AppPreferencesState {
  // Core project state
  lastProject: string;
  lastChain: string;
  lastAddress: string;
  lastContract: string;

  // UI state
  lastTheme: string;
  lastLanguage: string;
  lastView: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  lastTab: Record<string, types.DataFacet>;

  // Debug state
  debugMode: boolean;

  // Loading state
  loading: boolean;
}

// Initial state
const initialState: AppPreferencesState = {
  lastProject: '',
  lastChain: '',
  lastAddress: '0xf503017d7baf7fbc0fff7492b751025c6a78179b',
  lastContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
  lastTheme: 'dark',
  lastLanguage: 'en',
  lastView: '/',
  menuCollapsed: true,
  helpCollapsed: true,
  lastTab: {},
  debugMode: false,
  loading: true,
};

class AppPreferencesStore {
  private state: AppPreferencesState = { ...initialState };
  private listeners = new Set<() => void>();
  private isTestMode = false;

  // Set test mode to skip backend calls
  setTestMode = (testMode: boolean): void => {
    this.isTestMode = testMode;
    if (testMode) {
      // In test mode, set loading to false and use defaults
      this.setState({ loading: false });
    }
  };

  // Get current state (required by useSyncExternalStore)
  getState = (): AppPreferencesState => {
    return this.state;
  };

  // Subscribe to state changes (required by useSyncExternalStore)
  subscribe = (listener: () => void): (() => void) => {
    this.listeners.add(listener);

    // Return unsubscribe function
    return () => {
      this.listeners.delete(listener);
    };
  };

  // Notify all listeners of state changes
  private notify = (): void => {
    this.listeners.forEach((listener) => listener());
  };

  // Internal method to update state and notify listeners
  private setState = (updates: Partial<AppPreferencesState>): void => {
    this.state = { ...this.state, ...updates };
    this.notify();
  };

  // Helper function to update backend preferences
  private updatePreferences = async (
    updates: Partial<preferences.AppPreferences>,
  ): Promise<void> => {
    // Skip backend calls in test mode
    if (this.isTestMode) {
      return;
    }

    try {
      const currentPrefs = await App.GetAppPreferences();
      const updatedPrefs = preferences.AppPreferences.createFrom({
        ...currentPrefs,
        ...updates,
      });
      await App.SetAppPreferences(updatedPrefs);
    } catch (error) {
      console.error('Failed to update preferences:', error);
      throw error;
    }
  };

  // Initialize the store by loading state from backend
  initialize = async (): Promise<void> => {
    // Skip initialization in test mode
    if (this.isTestMode) {
      return;
    }

    try {
      this.setState({ loading: true });
      let retries = 0;
      const maxRetries = 10;
      let prefs = null;
      while (retries < maxRetries) {
        try {
          prefs = await App.GetAppPreferences();
          break; // Success, exit retry loop
        } catch (error) {
          retries++;
          console.log(
            `App preferences load attempt ${retries}/${maxRetries} failed:`,
            error,
          );
          if (retries >= maxRetries) {
            console.warn(
              'Failed to load app preferences after max retries, using defaults',
            );
            // Use safe defaults if backend is not ready
            this.setState({
              lastProject: '',
              lastChain: '',
              lastAddress: '0xf503017d7baf7fbc0fff7492b751025c6a78179b',
              lastTheme: 'dark',
              lastLanguage: 'en',
              lastView: '/',
              menuCollapsed: false,
              helpCollapsed: false,
              lastTab: {},
              debugMode: false,
              loading: false,
            });
            return;
          }
          // Wait before retry
          await new Promise((resolve) => setTimeout(resolve, 200));
        }
      }
      if (prefs) {
        this.setState({
          lastProject: prefs.lastProject || '',
          lastChain: prefs.lastChain || '',
          lastAddress:
            prefs.lastAddress || '0xf503017d7baf7fbc0fff7492b751025c6a78179b',
          lastContract:
            prefs.lastContract || '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
          lastTheme: prefs.lastTheme || 'dark',
          lastLanguage: prefs.lastLanguage || 'en',
          lastView: prefs.lastView || '/',
          menuCollapsed: prefs.menuCollapsed || false,
          helpCollapsed: prefs.helpCollapsed || false,
          lastTab: (prefs.lastTab || {}) as Record<string, types.DataFacet>,
          debugMode:
            (prefs as preferences.AppPreferences & { debugMode?: boolean })
              .debugMode || false,
          loading: false,
        });
      }
    } catch (error) {
      console.error('Failed to load app preferences:', error);
      this.setState({ loading: false });
    }
  };

  // Action methods that update both local state and backend
  setActiveAddress = async (address: string): Promise<void> => {
    await this.updatePreferences({ lastAddress: address });
    this.setState({ lastAddress: address });
  };

  setActiveChain = async (chain: string): Promise<void> => {
    await this.updatePreferences({ lastChain: chain });
    this.setState({ lastChain: chain });
  };

  setActiveContract = async (contract: string): Promise<void> => {
    // Skip backend calls in test mode
    if (!this.isTestMode) {
      await App.SetLastContract(contract);
    }
    this.setState({ lastContract: contract });
  };

  switchProject = async (project: string): Promise<void> => {
    await this.updatePreferences({ lastProject: project });
    this.setState({ lastProject: project });
  };

  toggleTheme = async (): Promise<void> => {
    const newTheme = this.state.lastTheme === 'dark' ? 'light' : 'dark';
    await this.updatePreferences({ lastTheme: newTheme });
    this.setState({ lastTheme: newTheme });
  };

  changeLanguage = async (language: string): Promise<void> => {
    await this.updatePreferences({ lastLanguage: language });
    this.setState({ lastLanguage: language });
  };

  setMenuCollapsed = async (collapsed: boolean): Promise<void> => {
    await this.updatePreferences({ menuCollapsed: collapsed });
    this.setState({ menuCollapsed: collapsed });
  };

  setHelpCollapsed = async (collapsed: boolean): Promise<void> => {
    await this.updatePreferences({ helpCollapsed: collapsed });
    this.setState({ helpCollapsed: collapsed });
  };

  setLastTab = async (route: string, tab: types.DataFacet): Promise<void> => {
    // Skip backend calls in test mode
    if (!this.isTestMode) {
      await App.SetLastTab(route, tab);
    }
    this.setState({
      lastTab: { ...this.state.lastTab, [route]: tab },
    });
  };

  setLastView = async (view: string): Promise<void> => {
    await this.updatePreferences({ lastView: view });
    this.setState({ lastView: view });
  };

  toggleDarkMode = async (): Promise<void> => {
    const newTheme = this.state.lastTheme === 'dark' ? 'light' : 'dark';
    await this.updatePreferences({ lastTheme: newTheme });
    this.setState({ lastTheme: newTheme });
  };

  toggleDebugMode = async (): Promise<void> => {
    const newDebugMode = !this.state.debugMode;
    await this.updatePreferences({
      debugMode: newDebugMode,
    } as Partial<preferences.AppPreferences & { debugMode: boolean }>);
    this.setState({ debugMode: newDebugMode });
  };

  // Computed getters
  get isDarkMode(): boolean {
    return this.state.lastTheme === 'dark';
  }

  get hasActiveProject(): boolean {
    return Boolean(this.state.lastProject);
  }

  get canExport(): boolean {
    return Boolean(this.state.lastProject && this.state.lastAddress);
  }
}

// Create and export a singleton instance
export const appPreferencesStore = new AppPreferencesStore();

// Initialize the store when the module loads (only when not in test environment)
if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  appPreferencesStore.initialize();
} else {
  // In test environment, enable test mode
  appPreferencesStore.setTestMode(true);
}
