// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const abisFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.DOWNLOADED,
    label: toProperCase(types.DataFacet.DOWNLOADED),
    isDefault: true,
  },
  {
    id: types.DataFacet.KNOWN,
    label: toProperCase(types.DataFacet.KNOWN),
    isDefault: false,
  },
  {
    id: types.DataFacet.FUNCTIONS,
    label: toProperCase(types.DataFacet.FUNCTIONS),
    isDefault: false,
  },
  {
    id: types.DataFacet.EVENTS,
    label: toProperCase(types.DataFacet.EVENTS),
    isDefault: false,
  },
];

export const DEFAULT_FACET = types.DataFacet.DOWNLOADED;
export const ROUTE = '/abis' as const;
