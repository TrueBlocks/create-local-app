import { DataFacet, DataFacetConfig, useActiveFacet } from '@hooks';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the useActiveProject hook
vi.mock('../useActiveProject', () => ({
  useActiveProject: vi.fn(),
}));

const { useActiveProject } = await import('../useActiveProject');
const mockedUseActiveProject = vi.mocked(useActiveProject);

describe('useActiveFacet Hook Tests (DataFacet implementation)', () => {
  let mockLastTab: Record<string, string>;
  let mockSetLastTab: ReturnType<typeof vi.fn>;

  const sampleFacets: DataFacetConfig[] = [
    {
      id: 'transactions' as DataFacet,
      label: 'Transactions',
      isDefault: true,
    },
    {
      id: 'receipts' as DataFacet,
      label: 'Receipts',
    },
    {
      id: 'statements' as DataFacet,
      label: 'Statements',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();

    mockLastTab = {};
    mockSetLastTab = vi.fn();

    mockedUseActiveProject.mockReturnValue({
      lastTab: mockLastTab,
      setLastTab: mockSetLastTab,
      // Mock other required properties
      lastProject: 'test-project',
      lastChain: 'mainnet',
      lastAddress: '0x123',
      lastContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      lastTheme: 'dark',
      lastLanguage: 'en',
      lastView: '/exports',
      menuCollapsed: false,
      helpCollapsed: false,
      loading: false,
    } as any);
  });

  describe('initialization and defaults', () => {
    it('should initialize with default facet when no preference saved', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('transactions');
      expect(result.current.getDefaultFacet()).toBe('transactions');
    });

    it('should use provided default facet override', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
          defaultFacet: 'statements' as DataFacet,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });

    it('should restore saved preference via backend API mapping', () => {
      mockLastTab['/exports'] = 'statements';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });
  });

  describe('facet switching', () => {
    it('should change active facet and persist backend API value to preferences', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      act(() => {
        result.current.setActiveFacet('receipts' as DataFacet);
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/exports', 'receipts');
    });
  });

  describe('facet configuration', () => {
    it('should return correct facet configuration', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      const config = result.current.getFacetConfig('receipts' as DataFacet);
      expect(config).toEqual({
        id: 'receipts',
        label: 'Receipts',
      });
    });

    it('should provide available facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.availableFacets).toHaveLength(3);
      expect(result.current.availableFacets).toEqual(sampleFacets);
    });

    it('should check if facet is active correctly', () => {
      mockLastTab['/exports'] = 'statements';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.isFacetActive('statements' as DataFacet)).toBe(
        true,
      );
      expect(result.current.isFacetActive('receipts' as DataFacet)).toBe(false);
    });
  });

  describe('backward compatibility', () => {
    it('should provide correct backend API value for current facet', () => {
      mockLastTab['/exports'] = 'statements';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.getCurrentDataFacet()).toBe('statements');
    });

    it('should fallback to facet ID when no backend API mapping exists', () => {
      const facetsWithoutApiMapping: DataFacetConfig[] = [
        {
          id: 'custom-facet' as DataFacet,
          label: 'Custom Facet',
          // No dataFacet property
        },
      ];

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: facetsWithoutApiMapping,
        }),
      );

      expect(result.current.getCurrentDataFacet()).toBe('custom-facet');
    });
  });

  describe('edge cases', () => {
    it('should handle empty facets array', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: [],
        }),
      );

      expect(result.current.activeFacet).toBe('all'); // fallback
      expect(result.current.availableFacets).toEqual([]);
    });

    it('should handle invalid saved preference', () => {
      mockLastTab['/exports'] = 'INVALID_DATAFACET';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('transactions'); // fallback to default
    });
  });
});
