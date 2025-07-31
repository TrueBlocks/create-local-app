<!--
Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
Use of this source code is governed by a license that can
be found in the LICENSE file.

Parts of this file were auto generated. Edit only those parts of
the code inside of 'EXISTING_CODE' tags.
-->
# Status View

// EXISTING_CODE
// EXISTING_CODE

## Facets:
- Status Facet uses Status store.
- Caches Facet uses Caches store.
- Chains Facet uses Chains store.

## Stores:

- **Caches Store (0 members):**


- **Chains Store (0 members):**


- **Status Store (18 members):**

  - cachePath: path to the cache directory
  - chain: the chain identifier
  - chainConfig: path to chain configuration
  - chainId: the chain ID
  - clientVersion: version of the client
  - hasEsKey: whether Etherscan API key is available
  - hasPinKey: whether Pinata API key is available
  - indexPath: path to the index directory
  - isApi: whether running in API mode
  - isArchive: whether node is archive node
  - isScraping: whether scraper is running
  - isTesting: whether in testing mode
  - isTracing: whether tracing is enabled
  - networkId: the network ID
  - progress: progress information
  - rootConfig: path to root configuration
  - rpcProvider: RPC provider URL
  - version: {{ORG_NAME}} version

// EXISTING_CODE
// EXISTING_CODE
