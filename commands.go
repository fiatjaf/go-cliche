package cliche

import "encoding/json"

func (c *Control) GetInfo() (result GetInfoResult, err error) {
	resultJson, err := c.Call("get-info", map[string]interface{}{})
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resultJson, &result)
	return result, err
}

func (c *Control) CreateInvoice(params CreateInvoiceParams) (
	result CreateInvoiceResult,
	err error,
) {
	resultJson, err := c.Call("create-invoice", params)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resultJson, &result)
	return result, err
}

func (c *Control) PayInvoice(params PayInvoiceParams) (
	result PayInvoiceResult,
	err error,
) {
	resultJson, err := c.Call("create-invoice", params)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resultJson, &result)
	return result, err
}
