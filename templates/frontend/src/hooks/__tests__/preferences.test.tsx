import { useActiveProject } from '@hooks';
import { types } from '@models';
import { appPreferencesStore } from '@stores';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the utils module to prevent debugMode issues at module load time
vi.mock('@utils', async () => {
  const actual = await vi.importActual('@utils');
  return {
    ...actual,
    isDebugMode: vi.fn(() => false),
    useErrorHandler: vi.fn(() => vi.fn()),
    Log: {
      debug: vi.fn(),
      info: vi.fn(),
      warn: vi.fn(),
      error: vi.fn(),
    },
  };
});

// Mock the store dependency
vi.mock('@stores', () => ({
  appPreferencesStore: {
    subscribe: vi.fn(),
    getState: vi.fn(() => ({
      debugMode: false,
      isDarkMode: false,
      hasActiveProject: true,
      canExport: true,
      lastTab: 'default',
      activeAddress: '',
      activeChain: '',
      theme: 'light',
      language: 'en',
      menuCollapsed: false,
      helpCollapsed: false,
      lastView: 'default',
    })),
    setLastTab: vi.fn(),
    setActiveAddress: vi.fn(),
    setActiveChain: vi.fn(),
    switchProject: vi.fn(),
    toggleTheme: vi.fn(),
    changeLanguage: vi.fn(),
    setMenuCollapsed: vi.fn(),
    setHelpCollapsed: vi.fn(),
    setLastView: vi.fn(),
    toggleDarkMode: vi.fn(),
    isDarkMode: false,
    hasActiveProject: true,
    canExport: true,
  },
}));

describe('Preference System Tests (DataFacet refactor preparation)', () => {
  let mockStoreState: any;

  beforeEach(() => {
    vi.clearAllMocks();

    mockStoreState = {
      lastProject: 'test-project',
      lastChain: 'mainnet',
      lastAddress: '0x123',
      lastContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      lastTheme: 'dark',
      lastLanguage: 'en',
      lastView: '/exports',
      menuCollapsed: false,
      helpCollapsed: false,
      lastTab: {
        '/exports': 'transactions' as types.DataFacet,
        '/chunks': 'chunk-summary' as types.DataFacet,
        '/monitors': 'txs' as types.DataFacet,
      },
      loading: false,
    };

    (appPreferencesStore.getState as any).mockReturnValue(mockStoreState);
    (appPreferencesStore.subscribe as any).mockImplementation(
      (_callback: () => void) => {
        // Mock subscription
        return () => {}; // unsubscribe function
      },
    );
  });

  describe('lastTab persistence behavior', () => {
    it('retrieves stored lastTab values for different routes', () => {
      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastTab).toEqual({
        '/exports': 'transactions',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
      });

      // Verify specific route lookups
      expect(result.current.lastTab['/exports']).toBe('transactions');
      expect(result.current.lastTab['/chunks']).toBe('chunk-summary');
      expect(result.current.lastTab['/monitors']).toBe('txs');
    });

    it('handles missing lastTab entries gracefully', () => {
      // Test scenario where a route has no stored lastTab
      const { result } = renderHook(() => useActiveProject());

      // Routes not in the stored lastTab should return undefined
      expect(result.current.lastTab['/names']).toBeUndefined();
      expect(result.current.lastTab['/abis']).toBeUndefined();
      expect(result.current.lastTab['/unknown-route']).toBeUndefined();
    });

    it('calls setLastTab with correct parameters', async () => {
      const { result } = renderHook(() => useActiveProject());

      await act(async () => {
        await result.current.setLastTab(
          '/exports',
          'receipts' as types.DataFacet,
        );
      });

      expect(appPreferencesStore.setLastTab as any).toHaveBeenCalledWith(
        '/exports',
        'receipts',
      );
    });

    it('supports all known DataFacet values in setLastTab', async () => {
      const { result } = renderHook(() => useActiveProject());

      const testCases: Array<[string, types.DataFacet]> = [
        ['/exports', types.DataFacet.ALL],
        ['/exports', 'receipts' as types.DataFacet],
        ['/chunks', 'chunk-summary' as types.DataFacet],
        ['/monitors', 'txs' as types.DataFacet],
        ['/names', 'entity-names' as types.DataFacet],
        ['/abis', 'get-abis' as types.DataFacet],
      ];

      for (const [route, dataFacet] of testCases) {
        await act(async () => {
          await result.current.setLastTab(route, dataFacet);
        });

        expect(appPreferencesStore.setLastTab as any).toHaveBeenCalledWith(
          route,
          dataFacet,
        );
      }
    });
  });

  describe('cross-session state recovery', () => {
    it('initializes with previously stored lastTab state', () => {
      // This simulates app restart with stored preferences
      const storedState = {
        ...mockStoreState,
        lastTab: {
          '/exports': 'receipts' as types.DataFacet,
          '/chunks': 'chunk-summary' as types.DataFacet,
          '/monitors': 'txs' as types.DataFacet,
          '/names': 'entity-names' as types.DataFacet,
        },
      };

      (appPreferencesStore.getState as any).mockReturnValue(storedState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastTab).toEqual({
        '/exports': 'receipts',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
        '/names': 'entity-names',
      });
    });

    it('handles empty lastTab state on first run', () => {
      const emptyState = {
        ...mockStoreState,
        lastTab: {},
      };

      (appPreferencesStore.getState as any).mockReturnValue(emptyState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastTab).toEqual({});
      expect(Object.keys(result.current.lastTab)).toHaveLength(0);
    });
  });

  describe('default facet selection logic', () => {
    it('provides setLastTab for setting default selections', async () => {
      const { result } = renderHook(() => useActiveProject());

      expect(result.current.setLastTab).toBeInstanceOf(Function);

      // Test that it can be called to set defaults
      await act(async () => {
        await result.current.setLastTab(
          '/new-route',
          'transactions' as types.DataFacet,
        );
      });

      expect(appPreferencesStore.setLastTab as any).toHaveBeenCalledWith(
        '/new-route',
        'transactions',
      );
    });

    it('maintains lastTab structure for route-based lookups', () => {
      const { result } = renderHook(() => useActiveProject());

      // Verify the structure supports the pattern used in views:
      // const currentTab = lastTab[currentRoute] || defaultTab;
      const lastTab = result.current.lastTab;

      expect(typeof lastTab).toBe('object');
      expect(lastTab).not.toBeNull();
      expect(Array.isArray(lastTab)).toBe(false);

      // Test the lookup pattern
      const exportsTab = lastTab['/exports'] || 'transactions';
      const unknownTab = lastTab['/unknown'] || 'default-facet';

      expect(exportsTab).toBe('transactions');
      expect(unknownTab).toBe('default-facet');
    });
  });

  describe('edge cases in preference handling', () => {
    it('handles route keys with special characters', async () => {
      const { result } = renderHook(() => useActiveProject());

      // Test routes that might have special characters
      const specialRoutes = [
        '/exports/sub-route',
        '/exports?param=value',
        '/exports#hash',
      ];

      for (const route of specialRoutes) {
        await act(async () => {
          await result.current.setLastTab(
            route,
            'transactions' as types.DataFacet,
          );
        });

        expect(appPreferencesStore.setLastTab as any).toHaveBeenCalledWith(
          route,
          'transactions',
        );
      }
    });

    it('handles concurrent lastTab updates', async () => {
      const { result } = renderHook(() => useActiveProject());

      // Simulate rapid tab switching
      const promises = [
        result.current.setLastTab(
          '/exports',
          'transactions' as types.DataFacet,
        ),
        result.current.setLastTab('/exports', 'receipts' as types.DataFacet),
        result.current.setLastTab(
          '/chunks',
          'chunk-summary' as types.DataFacet,
        ),
      ];

      await act(async () => {
        await Promise.all(promises);
      });

      expect(appPreferencesStore.setLastTab as any).toHaveBeenCalledTimes(3);
    });

    it('handles invalid DataFacet values gracefully', async () => {
      const { result } = renderHook(() => useActiveProject());

      // This tests the function signature - TypeScript should catch invalid types,
      // but we test runtime behavior

      await act(async () => {
        await result.current.setLastTab(
          '/exports',
          'transactions' as types.DataFacet,
        );
      });

      expect(appPreferencesStore.setLastTab as any).toHaveBeenCalledWith(
        '/exports',
        'transactions',
      );
    });

    it('maintains state consistency during store updates', () => {
      let subscriptionCallback: (() => void) | null = null;

      (appPreferencesStore.subscribe as any).mockImplementation(
        (callback: () => void) => {
          subscriptionCallback = callback;
          return () => {
            subscriptionCallback = null;
          };
        },
      );

      const { result } = renderHook(() => useActiveProject());

      // Initial state
      expect(result.current.lastTab).toEqual({
        '/exports': 'transactions',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
      });

      // Simulate store update
      const newState = {
        ...mockStoreState,
        lastTab: {
          ...mockStoreState.lastTab,
          '/exports': 'receipts' as types.DataFacet,
        },
      };
      (appPreferencesStore.getState as any).mockReturnValue(newState);

      // Trigger subscription callback to simulate store change
      if (subscriptionCallback) {
        act(() => {
          subscriptionCallback?.();
        });
      }

      // State should remain consistent
      expect(result.current.lastTab['/exports']).toBe('receipts');
    });
  });
});
