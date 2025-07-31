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
	"strings"

	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
	// EXISTING_CODE
	// EXISTING_CODE
)

// TODO: The slices should be slices to pointers
type MonitorsPage struct {
	Facet         types.DataFacet `json:"facet"`
	Monitors      []Monitor       `json:"monitors"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *MonitorsPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *MonitorsPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *MonitorsPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *MonitorsPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *MonitorsPage) GetState() types.LoadState {
	return p.State
}

func (c *MonitorsCollection) GetPage(
	payload *types.Payload,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	dataFacet := payload.DataFacet

	page := &MonitorsPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	switch dataFacet {
	case MonitorsMonitors:
		facet := c.monitorsFacet
		var filterFunc func(*Monitor) bool
		if filter != "" {
			filterFunc = func(item *Monitor) bool {
				return c.matchesMonitorFilter(item, filter)
			}
		}
		sortFunc := func(items []Monitor, sort sdk.SortSpec) error {
			return sdk.SortMonitors(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("monitors", dataFacet, "GetPage", err)
		} else {

			page.Monitors, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	default:
		return nil, types.NewValidationError("monitors", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	return page, nil
}

// EXISTING_CODE
// EXISTING_CODE
