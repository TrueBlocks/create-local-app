import {
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from 'react';

import { sdk } from '@models';

import { ViewStateKey, viewStateKeyToString } from '.';
import { createEmptySortSpec } from '../utils/sortSpec';

const EMPTY_SORT = createEmptySortSpec();

// Pagination interfaces
export interface PaginationState {
  currentPage: number;
  pageSize: number;
  totalItems: number;
}

// Navigation context for calculating post-deletion targets
export interface NavigationContext {
  currentPage: number;
  pageSize: number;
  totalItems: number;
  deletingRowIndex: number;
  deletingRowId: string;
  currentPageData: Record<string, unknown>[];
  // Performance hints
  isOnlyRowOnPage: boolean;
  isFirstRowOnPage: boolean;
  isLastRowOnPage: boolean;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

// Navigation target for post-action positioning
export interface NavigationTarget {
  type: 'page' | 'row' | 'none';
  page?: number;
  rowId?: string;
}

// Navigation strategy function type
export type NavigationStrategy = (
  context: NavigationContext,
) => NavigationTarget | null;

export interface ViewSortState {
  [key: string]: sdk.SortSpec | null;
}

export interface ViewFilterState {
  [key: string]: string;
}

export interface ViewPaginationState {
  [key: string]: PaginationState;
}

// Create stable reference for initial state to prevent new object creation
export const initialPaginationState: PaginationState = Object.freeze({
  currentPage: 0,
  pageSize: 15,
  totalItems: 0,
});

interface ViewContextType {
  currentView: string;
  setCurrentView: (view: string) => void;
  getPagination: (viewStateKey: ViewStateKey) => PaginationState;
  updatePagination: (
    viewStateKey: ViewStateKey,
    changes: Partial<PaginationState>,
  ) => void;
  getSorting: (viewStateKey: ViewStateKey) => sdk.SortSpec | null;
  updateSorting: (
    viewStateKey: ViewStateKey,
    sort: sdk.SortSpec | null,
  ) => void;
  getFiltering: (viewStateKey: ViewStateKey) => string;
  updateFiltering: (viewStateKey: ViewStateKey, filter: string) => void;
}

export const ViewContext = createContext<ViewContextType>({
  currentView: '',
  setCurrentView: () => {},
  getPagination: () => initialPaginationState,
  updatePagination: () => {},
  getSorting: () => null,
  updateSorting: () => {},
  getFiltering: () => '',
  updateFiltering: () => {},
});

export const ViewContextProvider = ({ children }: { children: ReactNode }) => {
  const [currentView, setCurrentView] = useState('');
  const [viewPagination, setViewPagination] = useState<ViewPaginationState>({});
  const [viewSorting, setViewSorting] = useState<ViewSortState>({});
  const [viewFiltering, setViewFiltering] = useState<ViewFilterState>({});

  const getPagination = useCallback(
    (viewStateKey: ViewStateKey) => {
      const key = viewStateKeyToString(viewStateKey);
      return viewPagination[key] || initialPaginationState;
    },
    [viewPagination],
  );

  const updatePagination = useCallback(
    (viewStateKey: ViewStateKey, changes: Partial<PaginationState>) => {
      setViewPagination((prev) => {
        const key = viewStateKeyToString(viewStateKey);
        const currentPagination = prev[key] || { ...initialPaginationState };
        return {
          ...prev,
          [key]: {
            ...currentPagination,
            ...changes,
          },
        };
      });
    },
    [],
  );

  const getSorting = useCallback(
    (viewStateKey: ViewStateKey) => {
      const key = viewStateKeyToString(viewStateKey);
      return viewSorting[key] || null;
    },
    [viewSorting],
  );

  const updateSorting = useCallback(
    (viewStateKey: ViewStateKey, sort: sdk.SortSpec | null) => {
      setViewSorting((prev) => {
        const key = viewStateKeyToString(viewStateKey);
        return {
          ...prev,
          [key]: sort,
        };
      });
    },
    [],
  );

  const getFiltering = useCallback(
    (viewStateKey: ViewStateKey) => {
      return viewFiltering[viewStateKey.viewName] || '';
    },
    [viewFiltering],
  );

  const updateFiltering = useCallback(
    (viewStateKey: ViewStateKey, filter: string) => {
      setViewFiltering((prev) => ({
        ...prev,
        [viewStateKey.viewName]: filter,
      }));
    },
    [],
  );

  const contextValue = useMemo(
    () => ({
      currentView,
      setCurrentView,
      getPagination,
      updatePagination,
      getSorting,
      updateSorting,
      getFiltering,
      updateFiltering,
    }),
    [
      currentView,
      setCurrentView,
      getPagination,
      updatePagination,
      getSorting,
      updateSorting,
      getFiltering,
      updateFiltering,
    ],
  );

  return (
    <ViewContext.Provider value={contextValue}>{children}</ViewContext.Provider>
  );
};

export const useViewContext = () => {
  const context = useContext(ViewContext);
  if (!context) {
    throw new Error('useViewContext must be used within a ViewContextProvider');
  }
  return context;
};

// Hook for sorting state (per-facet)
// Uses full ViewStateKey (viewName + tabName) to scope sorting per facet
export const useSorting = (viewStateKey: ViewStateKey) => {
  const { getSorting, updateSorting } = useViewContext();

  const setSorting = useCallback(
    (sort: sdk.SortSpec | null) => {
      updateSorting(viewStateKey, sort);
    },
    [viewStateKey, updateSorting],
  );

  // Get the sort value, but never return null - instead return a default empty SortSpec
  const sortValue = getSorting(viewStateKey);
  const sort = sortValue || EMPTY_SORT;

  return useMemo(() => ({ sort, setSorting }), [sort, setSorting]);
};

// Hook for filtering state (per-view)
// Uses only viewName from ViewStateKey to scope filtering per view
export const useFiltering = (viewStateKey: ViewStateKey) => {
  const { getFiltering, updateFiltering } = useViewContext();

  const setFiltering = useCallback(
    (filter: string) => {
      updateFiltering(viewStateKey, filter);
    },
    [viewStateKey, updateFiltering],
  );

  const filter = getFiltering(viewStateKey);

  return useMemo(() => ({ filter, setFiltering }), [filter, setFiltering]);
};
