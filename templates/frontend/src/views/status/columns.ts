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

// Column configurations for the Status data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.STATUS:
      return getStatusColumns();
    case types.DataFacet.CACHES:
      return getCachesColumns();
    case types.DataFacet.CHAINS:
      return getChainsColumns();
    default:
      return [];
  }
};

const getCachesColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'type',
    name: 'type',
    header: 'Type',
    label: 'Type',
    type: 'text',
  },
  {
    key: 'path',
    name: 'path',
    header: 'Path',
    label: 'Path',
    type: 'text',
  },
  {
    key: 'nFiles',
    name: 'nFiles',
    header: 'Number of Files',
    label: 'Number of Files',
    type: 'number',
  },
  {
    key: 'nFolders',
    name: 'nFolders',
    header: 'Number of Folders',
    label: 'Number of Folders',
    type: 'number',
  },
  {
    key: 'sizeInBytes',
    name: 'sizeInBytes',
    header: 'Size (Bytes)',
    label: 'Size (Bytes)',
    type: 'number',
  },
  {
    key: 'lastCached',
    name: 'lastCached',
    header: 'Last Cached',
    label: 'Last Cached',
    type: 'datetime',
  },
  // EXISTING_CODE
];

const getChainsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'chain',
    name: 'chain',
    header: 'Chain',
    label: 'Chain',
    type: 'text',
  },
  {
    key: 'chainId',
    name: 'chainId',
    header: 'Chain Id',
    label: 'Chain Id',
    type: 'text',
  },
  {
    key: 'ipfsGateway',
    name: 'ipfsGateway',
    header: 'IPFS Gateway',
    label: 'IPFS Gateway',
    type: 'text',
  },
  {
    key: 'localExplorer',
    name: 'localExplorer',
    header: 'Local Explorer',
    label: 'Local Explorer',
    type: 'text',
  },
  {
    key: 'remoteExplorer',
    name: 'remoteExplorer',
    header: 'Remote Explorer',
    label: 'Remote Explorer',
    type: 'text',
  },
  {
    key: 'rpcProvider',
    name: 'rpcProvider',
    header: 'RPC Provider',
    label: 'RPC Provider',
    type: 'text',
  },
  {
    key: 'symbol',
    name: 'symbol',
    header: 'Symbol',
    label: 'Symbol',
    type: 'text',
  },
  // EXISTING_CODE
];

const getStatusColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'cachePath',
    name: 'cachePath',
    header: 'Cache Path',
    label: 'Cache Path',
    type: 'text',
    width: '200px',
  },
  {
    key: 'chain',
    name: 'chain',
    header: 'Chain',
    label: 'Chain',
    type: 'text',
    width: '200px',
  },
  {
    key: 'chainConfig',
    name: 'chainConfig',
    header: 'Chain Config',
    label: 'Chain Config',
    type: 'text',
    width: '200px',
  },
  {
    key: 'chainId',
    name: 'chainId',
    header: 'Chain Id',
    label: 'Chain Id',
    type: 'text',
    width: '200px',
  },
  {
    key: 'clientVersion',
    name: 'clientVersion',
    header: 'Client Version',
    label: 'Client Version',
    type: 'text',
    width: '200px',
  },
  {
    key: 'hasEsKey',
    name: 'hasEsKey',
    header: 'Has Es Key',
    label: 'Has Es Key',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'hasPinKey',
    name: 'hasPinKey',
    header: 'Has Pin Key',
    label: 'Has Pin Key',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'indexPath',
    name: 'indexPath',
    header: 'Index Path',
    label: 'Index Path',
    type: 'text',
    width: '200px',
  },
  {
    key: 'isApi',
    name: 'isApi',
    header: 'Api',
    label: 'Api',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'isArchive',
    name: 'isArchive',
    header: 'Archive',
    label: 'Archive',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'isScraping',
    name: 'isScraping',
    header: 'Scraping',
    label: 'Scraping',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'isTesting',
    name: 'isTesting',
    header: 'Testing',
    label: 'Testing',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'isTracing',
    name: 'isTracing',
    header: 'Tracing',
    label: 'Tracing',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'networkId',
    name: 'networkId',
    header: 'Network Id',
    label: 'Network Id',
    type: 'text',
    width: '200px',
  },
  {
    key: 'progress',
    name: 'progress',
    header: 'Progress',
    label: 'Progress',
    type: 'text',
    width: '200px',
  },
  {
    key: 'rootConfig',
    name: 'rootConfig',
    header: 'Root Config',
    label: 'Root Config',
    type: 'text',
    width: '200px',
  },
  {
    key: 'rpcProvider',
    name: 'rpcProvider',
    header: 'Rpc Provider',
    label: 'Rpc Provider',
    type: 'text',
    width: '200px',
  },
  {
    key: 'version',
    name: 'version',
    header: 'Version',
    label: 'Version',
    type: 'text',
    width: '200px',
  },
];

// EXISTING_CODE
// EXISTING_CODE
