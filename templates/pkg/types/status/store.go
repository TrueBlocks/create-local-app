// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package status

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

type Cache = sdk.Cache
type Chain = sdk.Chain
type Status = sdk.Status

var (
	cachesStore   *store.Store[Cache]
	cachesStoreMu sync.Mutex

	chainsStore   *store.Store[Chain]
	chainsStoreMu sync.Mutex

	statusStore   *store.Store[Status]
	statusStoreMu sync.Mutex
)

func (c *StatusCollection) getCachesStore(facet types.DataFacet) *store.Store[Cache] {
	cachesStoreMu.Lock()
	defer cachesStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := cachesStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.StatusOptions{
				Globals:   sdk.Globals{Chain: chain},
				RenderCtx: ctx,
			}
			if _, _, err := opts.StatusCaches(); err != nil {
				wrappedErr := types.NewSDKError("status", StatusCaches, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Status Caches SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Cache {
			if it, ok := item.(*Cache); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Cache) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		cachesStore = theStore
	}

	return theStore
}

func (c *StatusCollection) getChainsStore(facet types.DataFacet) *store.Store[Chain] {
	chainsStoreMu.Lock()
	defer chainsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := chainsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.StatusOptions{
				Globals:   sdk.Globals{Chain: chain},
				RenderCtx: ctx,
			}
			if _, _, err := opts.StatusChains(); err != nil {
				wrappedErr := types.NewSDKError("status", StatusChains, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Status Chains SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Chain {
			if it, ok := item.(*Chain); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Chain) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		chainsStore = theStore
	}

	return theStore
}

func (c *StatusCollection) getStatusStore(facet types.DataFacet) *store.Store[Status] {
	statusStoreMu.Lock()
	defer statusStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := statusStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			opts := sdk.StatusOptions{
				Globals:   sdk.Globals{Chain: chain},
				RenderCtx: ctx,
			}
			if _, _, err := opts.StatusHealthcheck(); err != nil {
				wrappedErr := types.NewSDKError("status", StatusStatus, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Status Healthcheck SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Status {
			if it, ok := item.(*Status); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Status) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		statusStore = theStore
	}

	return theStore
}

func (c *StatusCollection) GetStoreName(dataFacet types.DataFacet, chain, address string) string {
	_ = chain
	_ = address
	name := ""
	switch dataFacet {
	case StatusStatus:
		name = "status-status"
	case StatusCaches:
		name = "status-caches"
	case StatusChains:
		name = "status-chains"
	default:
		return ""
	}
	return name
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetStatusCount(payload *types.Payload) (int, error) {
	chain := preferences.GetLastChain()
	countOpts := sdk.StatusOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
	}
	if countResult, _, err := countOpts.StatusCount(); err != nil {
		return 0, fmt.Errorf("StatusCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

var (
	collections   = make(map[store.CollectionKey]*StatusCollection)
	collectionsMu sync.Mutex
)

func GetStatusCollection(payload *types.Payload) *StatusCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	pl := *payload
	pl.Address = ""

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewStatusCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// EXISTING_CODE
