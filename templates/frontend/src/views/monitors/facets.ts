// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const monitorsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.MONITORS,
    label: toProperCase(types.DataFacet.MONITORS),
    isDefault: true,
  },
];

export const DEFAULT_FACET = types.DataFacet.MONITORS;
export const ROUTE = '/monitors' as const;
