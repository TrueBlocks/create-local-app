package app

import (
	"strings"
	"sync"

	"{{CHIFRA}}/pkg/base"
	sdk "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5"
)

var ensLock sync.Mutex

func (a *App) ConvertToAddress(addr string) (base.Address, bool) {
	if !strings.HasSuffix(addr, ".eth") {
		ret := base.HexToAddress(addr)
		return ret, ret != base.ZeroAddr
	}

	ensLock.Lock()
	ensAddr, exists := a.ensMap[addr]
	ensLock.Unlock()

	if exists {
		return ensAddr, true
	}

	// Try to get an ENS or return the same input
	opts := sdk.NamesOptions{
		Terms: []string{addr},
	}
	if names, _, err := opts.Names(); err != nil {
		// msgs.Send(a.ctx, msgs.Error, msgs.NewErrorMsg(err))
		return base.ZeroAddr, false
	} else {
		if len(names) > 0 {
			ensLock.Lock()
			defer ensLock.Unlock()
			a.ensMap[addr] = names[0].Address
			return names[0].Address, true
		} else {
			ret := base.HexToAddress(addr)
			return ret, ret != base.ZeroAddr
		}
	}
}
