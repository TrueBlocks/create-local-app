package names

import "{{CHIFRA}}/pkg/base"

func (c *NamesCollection) NameFromAddress(address base.Address) (*Name, bool) {
	return namesStore.GetItemFromMap(address)
}
