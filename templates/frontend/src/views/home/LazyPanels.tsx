import { lazy } from 'react';

// Lazy load panel components
export const LazyProjectsPanel = lazy(() =>
  import('./panels/ProjectsPanel').then((module) => ({
    default: module.ProjectsPanel,
  })),
);

export const LazyNamesPanel = lazy(() =>
  import('./panels/NamesPanel').then((module) => ({
    default: module.NamesPanel,
  })),
);

export const LazyMonitorsPanel = lazy(() =>
  import('./panels/MonitorsPanel').then((module) => ({
    default: module.MonitorsPanel,
  })),
);

export const LazyExportsPanel = lazy(() =>
  import('./panels/ExportsPanel').then((module) => ({
    default: module.ExportsPanel,
  })),
);

export const LazyChunksPanel = lazy(() =>
  import('./panels/ChunksPanel').then((module) => ({
    default: module.ChunksPanel,
  })),
);

export const LazyAbisPanel = lazy(() =>
  import('./panels/AbisPanel').then((module) => ({
    default: module.AbisPanel,
  })),
);
