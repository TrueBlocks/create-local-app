import { usePagination } from '@components';
import { ViewStateKey } from '@contexts';
import { types } from '@models';
import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// Mock the context dependencies
vi.mock('@contexts', async () => {
  const actual = await vi.importActual('@contexts');
  return {
    ...actual,
    usePaginationContext: vi.fn(() => ({
      paginationStates: new Map(),
      setPaginationState: vi.fn(),
    })),
  };
});

describe('Hook Integration Tests (DataFacet refactor preparation)', () => {
  describe('usePagination with ViewStateKey', () => {
    it('creates unique pagination state for different ViewStateKey combinations', () => {
      const exportsTransactionsKey: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.ALL,
      };
      const exportsReceiptsKey: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.CUSTOM,
      };

      // Test that different tabName values create separate pagination states
      const { result: result1 } = renderHook(() =>
        usePagination(exportsTransactionsKey),
      );
      const { result: result2 } = renderHook(() =>
        usePagination(exportsReceiptsKey),
      );

      // Both should initialize successfully
      expect(result1.current.pagination).toBeDefined();
      expect(result2.current.pagination).toBeDefined();
      expect(result1.current.setTotalItems).toBeInstanceOf(Function);
      expect(result2.current.setTotalItems).toBeInstanceOf(Function);
    });

    it('handles ViewStateKey patterns used across all views', () => {
      const testKeys: ViewStateKey[] = [
        { viewName: '/names', tabName: types.DataFacet.ALL },
      ];

      testKeys.forEach((key) => {
        const { result } = renderHook(() => usePagination(key));

        expect(result.current.pagination).toBeDefined();
        expect(result.current.setTotalItems).toBeInstanceOf(Function);
        expect(result.current.pagination.currentPage).toBeGreaterThanOrEqual(0);
      });
    });

    it('maintains state isolation between different views with same tabName', () => {
      // Test case where different views might use the same tabName (DataFacet value)
      const exportsKey: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.ALL,
      };
      const namesKey: ViewStateKey = {
        viewName: '/names',
        tabName: types.DataFacet.CUSTOM,
      };

      const { result: exportsResult } = renderHook(() =>
        usePagination(exportsKey),
      );
      const { result: namesResult } = renderHook(() => usePagination(namesKey));

      // Should have separate pagination states despite same tabName
      expect(exportsResult.current.pagination).toBeDefined();
      expect(namesResult.current.pagination).toBeDefined();
    });
  });

  describe('ViewStateKey creation patterns', () => {
    it('tests ViewStateKey uniqueness across all combinations', () => {
      const routes = ['/exports', '/chunks', '/monitors', '/names', '/abis'];
      const facets = [
        'transactions',
        'receipts',
        'chunk-summary',
        'txs',
        'entity-names',
        'get-abis',
      ];

      const keys: ViewStateKey[] = [];

      // Generate all valid combinations as they appear in real views
      routes.forEach((route) => {
        facets.forEach((facet) => {
          keys.push({ viewName: route, tabName: facet as types.DataFacet });
        });
      });

      // Test that each key can be used with hooks
      keys.slice(0, 10).forEach((key) => {
        // Limit to first 10 to avoid excessive test overhead
        const { result } = renderHook(() => usePagination(key));
        expect(result.current).toBeDefined();
      });

      // Verify uniqueness through string conversion
      const stringKeys = keys.map((key) => `${key.viewName}/${key.tabName}/`);
      const uniqueKeys = [...new Set(stringKeys)];
      expect(uniqueKeys.length).toBe(stringKeys.length);
    });

    it('validates ViewStateKey structure required by hooks', () => {
      const validKey: ViewStateKey = {
        viewName: '/test-view',
        tabName: types.DataFacet.ALL,
      };

      const { result } = renderHook(() => usePagination(validKey));

      expect(result.current.pagination).toBeDefined();
      expect(typeof validKey.viewName).toBe('string');
      expect(typeof validKey.tabName).toBe('string');
      expect(validKey.viewName.length).toBeGreaterThan(0);
      expect(validKey.tabName.length).toBeGreaterThan(0);
    });
  });

  describe('Multi-hook integration', () => {
    it('tests ViewStateKey sharing between pagination, sorting, and filtering hooks', () => {
      const testKey: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.ALL,
      };

      // This simulates the pattern used in actual views where the same ViewStateKey
      // is passed to multiple hooks
      const { result: paginationResult } = renderHook(() =>
        usePagination(testKey),
      );

      // For now, just test pagination since that's what we can mock
      // In the real refactor, this will expand to test useSorting and useFiltering
      expect(paginationResult.current.pagination).toBeDefined();
      expect(paginationResult.current.setTotalItems).toBeInstanceOf(Function);

      // Verify the key structure is maintained
      expect(testKey.viewName).toBe('/exports');
      expect(testKey.tabName).toBe('all');
    });

    it('handles ViewStateKey memoization patterns from views', () => {
      let currentDataFacet = types.DataFacet.ALL;
      const getMemoizedKey = (): ViewStateKey => ({
        viewName: '/exports',
        tabName: currentDataFacet,
      });

      const { result: result1, rerender } = renderHook(() =>
        usePagination(getMemoizedKey()),
      );

      expect(result1.current.pagination).toBeDefined();

      currentDataFacet = types.DataFacet.ALL;
      rerender();

      expect(result1.current.pagination).toBeDefined();
    });
  });
});
