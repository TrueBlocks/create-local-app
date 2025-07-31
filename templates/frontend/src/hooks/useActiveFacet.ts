import { useCallback, useMemo } from 'react';

import { types } from '@models';

import { useActiveProject } from './useActiveProject';

export type DataFacet = types.DataFacet;

export interface DataFacetConfig {
  id: DataFacet;
  label: string;
  isDefault?: boolean;
  dividerBefore?: boolean;
}

export interface UseActiveFacetReturn {
  activeFacet: DataFacet;
  setActiveFacet: (facet: DataFacet) => void;
  availableFacets: DataFacetConfig[];
  getFacetConfig: (facet: DataFacet) => DataFacetConfig | undefined;
  isFacetActive: (facet: DataFacet) => boolean;
  getDefaultFacet: () => DataFacet;
  getCurrentDataFacet: () => types.DataFacet;
}

export interface UseActiveFacetParams {
  viewRoute: string;
  facets: DataFacetConfig[];
  defaultFacet?: DataFacet;
}

/**
 * Hook for managing active data facets in views
 *
 * This hook provides a clean API for:
 * - Managing the currently active facet in a view
 * - Persisting facet selection to user preferences
 * - Providing facet configuration and metadata
 * - Backward compatibility with existing DataFacet usage
 *
 * @param params Configuration for the hook
 * @returns Hook interface for facet management
 */
export const useActiveFacet = (
  params: UseActiveFacetParams,
): UseActiveFacetReturn => {
  const { viewRoute, facets, defaultFacet } = params;
  const { lastTab, setLastTab } = useActiveProject();

  // Determine the default facet for this view
  const computedDefaultFacet = useMemo((): DataFacet => {
    if (defaultFacet) return defaultFacet;

    const defaultConfig = facets.find((f) => f.isDefault);
    if (defaultConfig) return defaultConfig.id;

    // Fallback to first facet
    return facets[0]?.id || types.DataFacet.ALL;
  }, [defaultFacet, facets]);

  // Get the currently active facet from preferences or default
  const activeFacet = useMemo((): DataFacet => {
    const saved = lastTab[viewRoute];

    if (saved) {
      const savedAsFacet = saved as unknown as DataFacet;
      if (facets.some((f) => f.id === savedAsFacet)) {
        return savedAsFacet;
      }
    }
    return computedDefaultFacet;
  }, [lastTab, viewRoute, facets, computedDefaultFacet]);

  // Function to change the active facet
  const setActiveFacet = useCallback(
    (facet: DataFacet): void => {
      setLastTab(viewRoute, facet);
    },
    [setLastTab, viewRoute],
  );

  // Get configuration for a specific facet
  const getFacetConfig = useCallback(
    (facet: DataFacet): DataFacetConfig | undefined => {
      return facets.find((f) => f.id === facet);
    },
    [facets],
  );

  // Check if a facet is active
  const isFacetActive = useCallback(
    (facet: DataFacet): boolean => {
      return activeFacet === facet;
    },
    [activeFacet],
  );

  // Get the default facet for the view
  const getDefaultFacet = useCallback((): DataFacet => {
    return computedDefaultFacet;
  }, [computedDefaultFacet]);

  const getCurrentDataFacet = useCallback((): DataFacet => {
    return activeFacet;
  }, [activeFacet]);

  return {
    activeFacet,
    setActiveFacet,
    availableFacets: facets,
    getFacetConfig,
    isFacetActive,
    getDefaultFacet,
    getCurrentDataFacet,
  };
};
