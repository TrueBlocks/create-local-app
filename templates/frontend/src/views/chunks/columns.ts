// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import { FormField } from '@components';
import { types } from '@models';

// EXISTING_CODE
// EXISTING_CODE

// Column configurations for the Chunks data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.STATS:
      return getStatsColumns();
    case types.DataFacet.INDEX:
      return getIndexColumns();
    case types.DataFacet.BLOOMS:
      return getBloomsColumns();
    case types.DataFacet.MANIFEST:
      return getManifestColumns();
    default:
      return [];
  }
};

const getBloomsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    type: 'blkrange',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    type: 'hash',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nBlooms',
    name: 'nBlooms',
    header: 'Blooms',
    label: 'Blooms',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nInserted',
    name: 'nInserted',
    header: 'Inserted',
    label: 'Inserted',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'byteWidth',
    name: 'byteWidth',
    header: 'Byte Width',
    label: 'Byte Width',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
];

const getIndexColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    type: 'blkrange',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    type: 'hash',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nAddresses',
    name: 'nAddresses',
    header: 'Addresses',
    label: 'Addresses',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nAppearances',
    name: 'nAppearances',
    header: 'Appearances',
    label: 'Appearances',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    type: 'number',
    width: '150px',
    readOnly: true,
  },
];

const getManifestColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'version',
    name: 'version',
    header: 'Version',
    label: 'Version',
    type: 'text',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'chain',
    name: 'chain',
    header: 'Chain',
    label: 'Chain',
    type: 'text',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'specification',
    name: 'specification',
    header: 'Specification',
    label: 'Specification',
    type: 'ipfshash',
    width: '100px',
    readOnly: true,
  },
];

const getStatsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    type: 'blkrange',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nAddrs',
    name: 'nAddrs',
    header: 'Addrs',
    label: 'Addrs',
    type: 'number',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'nApps',
    name: 'nApps',
    header: 'Apps',
    label: 'Apps',
    type: 'number',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'nBlocks',
    name: 'nBlocks',
    header: 'Blocks',
    label: 'Blocks',
    type: 'number',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'nBlooms',
    name: 'nBlooms',
    header: 'Blooms',
    label: 'Blooms',
    type: 'number',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'recWid',
    name: 'recWid',
    header: 'Rec Wid',
    label: 'Rec Wid',
    type: 'number',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'bloomSz',
    name: 'bloomSz',
    header: 'Bloom Sz',
    label: 'Bloom Sz',
    type: 'number',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'chunkSz',
    name: 'chunkSz',
    header: 'Chunk Sz',
    label: 'Chunk Sz',
    type: 'number',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'addrsPerBlock',
    name: 'addrsPerBlock',
    header: 'Addrs Per Block',
    label: 'Addrs Per Block',
    type: 'float64',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'appsPerBlock',
    name: 'appsPerBlock',
    header: 'Apps Per Block',
    label: 'Apps Per Block',
    type: 'float64',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'appsPerAddr',
    name: 'appsPerAddr',
    header: 'Apps Per Addr',
    label: 'Apps Per Addr',
    type: 'float64',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'ratio',
    name: 'ratio',
    header: 'Ratio',
    label: 'Ratio',
    type: 'float64',
    width: '100px',
    readOnly: true,
  },
];

// EXISTING_CODE
// EXISTING_CODE
