// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const statusFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATUS,
    label: toProperCase(types.DataFacet.STATUS),
    isDefault: true,
  },
  {
    id: types.DataFacet.CACHES,
    label: toProperCase(types.DataFacet.CACHES),
    isDefault: false,
  },
  {
    id: types.DataFacet.CHAINS,
    label: toProperCase(types.DataFacet.CHAINS),
    isDefault: false,
  },
];

export const DEFAULT_FACET = types.DataFacet.STATUS;
export const ROUTE = '/status' as const;
