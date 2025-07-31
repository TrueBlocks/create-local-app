// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const namesFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.ALL,
    label: toProperCase(types.DataFacet.ALL),
    isDefault: true,
  },
  {
    id: types.DataFacet.CUSTOM,
    label: toProperCase(types.DataFacet.CUSTOM),
    isDefault: false,
  },
  {
    id: types.DataFacet.PREFUND,
    label: toProperCase(types.DataFacet.PREFUND),
    isDefault: false,
    dividerBefore: true,
  },
  {
    id: types.DataFacet.REGULAR,
    label: toProperCase(types.DataFacet.REGULAR),
    isDefault: false,
  },
  {
    id: types.DataFacet.BADDRESS,
    label: toProperCase(types.DataFacet.BADDRESS),
    isDefault: false,
  },
];

export const DEFAULT_FACET = types.DataFacet.ALL;
export const ROUTE = '/names' as const;
