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
	"strings"

	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
	// EXISTING_CODE
	// EXISTING_CODE
)

type StatusPage struct {
	Facet         types.DataFacet `json:"facet"`
	Caches        []*Cache        `json:"caches"`
	Chains        []*Chain        `json:"chains"`
	Status        []*Status       `json:"status"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *StatusPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *StatusPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *StatusPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *StatusPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *StatusPage) GetState() types.LoadState {
	return p.State
}

func (c *StatusCollection) GetPage(
	payload *types.Payload,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	dataFacet := payload.DataFacet

	page := &StatusPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	switch dataFacet {
	case StatusStatus:
		facet := c.statusFacet
		var filterFunc func(*Status) bool
		if filter != "" {
			filterFunc = func(item *Status) bool {
				return c.matchesStatusFilter(item, filter)
			}
		}
		sortFunc := func(items []Status, sort sdk.SortSpec) error {
			return sdk.SortStatus(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("status", dataFacet, "GetPage", err)
		} else {
			status := make([]*Status, 0, len(result.Items))
			for i := range result.Items {
				status = append(status, &result.Items[i])
			}
			page.Status, page.TotalItems, page.State = status, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case StatusCaches:
		facet := c.cachesFacet
		var filterFunc func(*Cache) bool
		if filter != "" {
			filterFunc = func(item *Cache) bool {
				return c.matchesCacheFilter(item, filter)
			}
		}
		sortFunc := func(items []Cache, sort sdk.SortSpec) error {
			return sdk.SortCaches(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("status", dataFacet, "GetPage", err)
		} else {
			cache := make([]*Cache, 0, len(result.Items))
			for i := range result.Items {
				cache = append(cache, &result.Items[i])
			}
			page.Caches, page.TotalItems, page.State = cache, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case StatusChains:
		facet := c.chainsFacet
		var filterFunc func(*Chain) bool
		if filter != "" {
			filterFunc = func(item *Chain) bool {
				return c.matchesChainFilter(item, filter)
			}
		}
		sortFunc := func(items []Chain, sort sdk.SortSpec) error {
			return sdk.SortChains(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("status", dataFacet, "GetPage", err)
		} else {
			chain := make([]*Chain, 0, len(result.Items))
			for i := range result.Items {
				chain = append(chain, &result.Items[i])
			}
			page.Chains, page.TotalItems, page.State = chain, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	default:
		return nil, types.NewValidationError("status", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	return page, nil
}

// EXISTING_CODE
func (c *StatusCollection) matchesCacheFilter(item *Cache, filter string) bool {
	_ = item
	_ = filter
	return true
}

func (c *StatusCollection) matchesChainFilter(item *Chain, filter string) bool {
	_ = item
	_ = filter
	return true
}

func (c *StatusCollection) matchesStatusFilter(item *Status, filter string) bool {
	_ = item
	_ = filter
	return true
}

// EXISTING_CODE
