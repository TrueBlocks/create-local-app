package store

import (
	"{{CHIFRA}}/pkg/base"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
)

type CollectionKey struct {
	Chain   string       // may be empty
	Address base.Address // may be empty
}

func GetCollectionKey(payload *types.Payload) CollectionKey {
	return CollectionKey{
		Chain:   payload.Chain,
		Address: base.HexToAddress(payload.Address)}
}
