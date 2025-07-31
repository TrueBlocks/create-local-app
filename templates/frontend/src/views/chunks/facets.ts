// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const chunksFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATS,
    label: toProperCase(types.DataFacet.STATS),
    isDefault: true,
  },
  {
    id: types.DataFacet.INDEX,
    label: toProperCase(types.DataFacet.INDEX),
    isDefault: false,
  },
  {
    id: types.DataFacet.BLOOMS,
    label: toProperCase(types.DataFacet.BLOOMS),
    isDefault: false,
  },
  {
    id: types.DataFacet.MANIFEST,
    label: toProperCase(types.DataFacet.MANIFEST),
    isDefault: false,
  },
];

export const DEFAULT_FACET = types.DataFacet.STATS;
export const ROUTE = '/chunks' as const;
