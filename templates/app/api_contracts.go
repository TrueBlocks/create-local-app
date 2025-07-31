// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	"fmt"

	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/logging"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types/contracts"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
	// EXISTING_CODE
	// EXISTING_CODE
)

func (a *App) GetContractsPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*contracts.ContractsPage, error) {
	logging.LogBackend(fmt.Sprintf("🚀 GetContractsPage called with payload: DataFacet=%s, Address=%s, first=%d, pageSize=%d",
		payload.DataFacet, payload.Address, first, pageSize))

	collection := contracts.GetContractsCollection(payload)
	logging.LogBackend(fmt.Sprintf("📚 Collection retrieved: %T", collection))

	result, err := getCollectionPage[*contracts.ContractsPage](collection, payload, first, pageSize, sort, filter)
	if err != nil {
		logging.LogBackend(fmt.Sprintf("❌ GetContractsPage error: %v", err))
		return nil, err
	}

	logging.LogBackend(fmt.Sprintf("✅ GetContractsPage success: totalItems=%d, contracts=%d, isFetching=%v",
		result.TotalItems, len(result.Contracts), result.IsFetching))

	return result, nil
}

func (a *App) GetContractsSummary(payload *types.Payload) types.Summary {
	collection := contracts.GetContractsCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadContracts(payload *types.Payload) error {
	collection := contracts.GetContractsCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

// EXISTING_CODE
// EXISTING_CODE
