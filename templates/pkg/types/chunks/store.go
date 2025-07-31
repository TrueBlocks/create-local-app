// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package chunks

import (
	"fmt"
	"sync"

	// EXISTING_CODE
	// EXISTING_CODE
	"{{CHIFRA}}/pkg/output"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/logging"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/preferences"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/store"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
)

// EXISTING_CODE
// EXISTING_CODE

type Bloom = sdk.Bloom
type Index = sdk.Index
type Manifest = sdk.Manifest
type Stats = sdk.Stats

var (
	bloomsStore   *store.Store[Bloom]
	bloomsStoreMu sync.Mutex

	indexStore   *store.Store[Index]
	indexStoreMu sync.Mutex

	manifestStore   *store.Store[Manifest]
	manifestStoreMu sync.Mutex

	statsStore   *store.Store[Stats]
	statsStoreMu sync.Mutex
)

func (c *ChunksCollection) getBloomsStore(facet types.DataFacet) *store.Store[Bloom] {
	bloomsStoreMu.Lock()
	defer bloomsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := bloomsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: false, // Set to false to avoid weird output issues
					Chain:   chain,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksBlooms(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksBlooms, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks blooms SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Bloom {
			if it, ok := item.(*Bloom); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Bloom) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		bloomsStore = theStore
	}

	return theStore
}

func (c *ChunksCollection) getIndexStore(facet types.DataFacet) *store.Store[Index] {
	indexStoreMu.Lock()
	defer indexStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := indexStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chain,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksIndex(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksIndex, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks index SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Index {
			if it, ok := item.(*Index); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Index) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// TODO: Do we need a mapping function for chunks. I think yes
			// if item != nil {
			// 	return item.Range, true
			// }
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		indexStore = theStore
	}

	return theStore
}

func (c *ChunksCollection) getManifestStore(facet types.DataFacet) *store.Store[Manifest] {
	manifestStoreMu.Lock()
	defer manifestStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := manifestStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chain,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksManifest(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksManifest, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks manifest SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Manifest {
			if it, ok := item.(*Manifest); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Manifest) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		manifestStore = theStore
	}

	return theStore
}

func (c *ChunksCollection) getStatsStore(facet types.DataFacet) *store.Store[Stats] {
	statsStoreMu.Lock()
	defer statsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := statsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chain,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksStats(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksStats, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks stats SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Stats {
			if it, ok := item.(*Stats); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Stats) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		statsStore = theStore
	}

	return theStore
}

func (c *ChunksCollection) GetStoreName(dataFacet types.DataFacet, chain, address string) string {
	_ = chain
	_ = address
	name := ""
	switch dataFacet {
	case ChunksStats:
		name = "chunks-stats"
	case ChunksIndex:
		name = "chunks-index"
	case ChunksBlooms:
		name = "chunks-blooms"
	case ChunksManifest:
		name = "chunks-manifest"
	default:
		return ""
	}
	return name
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetChunksCount(payload *types.Payload) (int, error) {
	chain := preferences.GetLastChain()
	countOpts := sdk.ChunksOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
	}
	if countResult, _, err := countOpts.ChunksCount(); err != nil {
		return 0, fmt.Errorf("ChunksCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

var (
	collections   = make(map[store.CollectionKey]*ChunksCollection)
	collectionsMu sync.Mutex
)

func GetChunksCollection(payload *types.Payload) *ChunksCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	pl := *payload
	pl.Address = ""

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewChunksCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// EXISTING_CODE
