// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

// EXISTING_CODE
import (
	"{{CHIFRA}}/pkg/base"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types/monitors"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types/names"
)

// EXISTING_CODE

func (a *App) Reload(payload *types.Payload) error {
	lastView := a.GetAppPreferences().LastView

	switch lastView {
	case "/exports":
		return a.ReloadExports(payload)
	case "/monitors":
		return a.ReloadMonitors(payload)
	case "/abis":
		return a.ReloadAbis(payload)
	case "/names":
		return a.ReloadNames(payload)
	case "/chunks":
		return a.ReloadChunks(payload)
	case "/contracts":
		return a.ReloadContracts(payload)
	case "/status":
		return a.ReloadStatus(payload)
	}

	return nil
}

// EXISTING_CODE
func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	collection := names.GetNamesCollection(&types.Payload{})
	return collection.NameFromAddress(base.HexToAddress(address))
}

func (a *App) MonitorsClean(payload *types.Payload, addresses []string) error {
	collection := monitors.GetMonitorsCollection(payload)
	return collection.Clean(addresses)
}

// EXISTING_CODE
