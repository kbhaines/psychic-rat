package admin

import (
	"psychic-rat/types"
)

// apiTxn wraps the error handling of multiple transactions with the API; the
// user just checks the 'err' field at the end of the transaction block.
type apiTxn struct {
	err error
}

func (a *apiTxn) addCompany(co types.Company) (c *types.Company) {
	if a.err != nil {
		return &co
	}
	c, a.err = companyAPI.AddCompany(co)
	return c
}

func (a *apiTxn) getCompany(id int) (c *types.Company) {
	if a.err != nil {
		return c
	}
	var co types.Company
	co, a.err = companyAPI.GetCompany(id)
	return &co
}

func (a *apiTxn) addItem(item types.Item) (i *types.Item) {
	if a.err != nil {
		return &item
	}
	i, a.err = itemsAPI.AddItem(item)
	return i
}

func (a *apiTxn) currencyConversion(id int, value int) (v int) {
	if a.err != nil {
		return 0
	}
	v, a.err = itemsAPI.CurrencyConversion(id, value)
	return v
}

func (a *apiTxn) getItem(id int) (i *types.Item) {
	if a.err != nil {
		return i
	}
	var item types.Item
	item, a.err = itemsAPI.GetItem(id)
	return &item

}

func (a *apiTxn) addPledge(item *types.Item, userID string, usdValue int) {
	if a.err != nil {
		return
	}
	_, a.err = pledgeAPI.AddPledge(item.ID, userID, usdValue)
}

func (a *apiTxn) deleteNewItem(id int) {
	if a.err != nil {
		return
	}
	a.err = newItemsAPI.DeleteNewItem(id)
}
