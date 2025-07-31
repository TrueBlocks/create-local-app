// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package monitors

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

type Monitor = sdk.Monitor

var (
	monitorsStore   *store.Store[Monitor]
	monitorsStoreMu sync.Mutex
)

func (c *MonitorsCollection) getMonitorsStore(facet types.DataFacet) *store.Store[Monitor] {
	monitorsStoreMu.Lock()
	defer monitorsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := monitorsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			listOpts := sdk.MonitorsOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
				RenderCtx: ctx,
			}
			if _, _, err := listOpts.MonitorsList(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("monitors", MonitorsMonitors, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Monitors SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Monitor {
			if it, ok := item.(*Monitor); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Monitor) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		monitorsStore = theStore
	}

	return theStore
}

func (c *MonitorsCollection) GetStoreName(dataFacet types.DataFacet, chain, address string) string {
	_ = chain
	_ = address
	name := ""
	switch dataFacet {
	case MonitorsMonitors:
		name = "monitors-monitors"
	default:
		return ""
	}
	return name
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetMonitorsCount(payload *types.Payload) (int, error) {
	chain := preferences.GetLastChain()
	countOpts := sdk.MonitorsOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
	}
	if countResult, _, err := countOpts.MonitorsCount(); err != nil {
		return 0, fmt.Errorf("MonitorsCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

var (
	collections   = make(map[store.CollectionKey]*MonitorsCollection)
	collectionsMu sync.Mutex
)

func GetMonitorsCollection(payload *types.Payload) *MonitorsCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	pl := *payload
	pl.Address = ""

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewMonitorsCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// EXISTING_CODE
