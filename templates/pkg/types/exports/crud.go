package exports

import (
	"{{CHIFRA}}/pkg/crud"
	"github.com/{{ORG_NAME}}/{{SLUG}}/pkg/types"
)

// Crud implements CRUD operations for exports - placeholder implementation
func (c *ExportsCollection) Crud(
	payload *types.Payload,
	op crud.Operation,
	item interface{},
) error {
	// Placeholder implementation - no SDK interaction yet
	// When SDK support is added, implement similar to other collections
	return nil
}
