// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	"{{CHIFRA}}/pkg/crud"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types/abis"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
	// EXISTING_CODE
	// EXISTING_CODE
)

func (a *App) GetAbisPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*abis.AbisPage, error) {
	collection := abis.GetAbisCollection(payload)
	return getCollectionPage[*abis.AbisPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) AbisCrud(
	payload *types.Payload,
	op crud.Operation,
	item *abis.Abi,
) error {
	collection := abis.GetAbisCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetAbisSummary(payload *types.Payload) types.Summary {
	collection := abis.GetAbisCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadAbis(payload *types.Payload) error {
	collection := abis.GetAbisCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

// EXISTING_CODE
// EXISTING_CODE
