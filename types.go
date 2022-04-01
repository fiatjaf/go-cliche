package cliche

import "encoding/json"

type JSONRPCMessage struct {
	Version string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type JSONRPCResponse struct {
	Version string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type GetInfoResult struct {
	Keys struct {
		Pub       string `json:"pub"`
		Priv      string `json:"priv"`
		Mnemonics string `json:"mnemonics"`
	} `json:"keys"`
	BlockHeight int `json:"block_height"`
	Wallets     []struct {
		Label   string `json:"label"`
		Balance int    `json:"balance"`
	} `json:"wallets"`
	Channels []struct {
		ID      string `json:"id"`
		Balance int    `json:"balance"`
	} `json:"channels"`
	KnownChannels struct {
		Normal int `json:"normal"`
		Hosted int `json:"hosted"`
	} `json:"known_channels"`
	OutgoingPayments []interface{} `json:"outgoing_payments"`
	FiatRates        struct {
		Usd float64 `json:"usd"`
	} `json:"fiat_rates"`
	FeeRates struct {
		Num1   int `json:"1"`
		Num10  int `json:"10"`
		Num100 int `json:"100"`
	} `json:"fee_rates"`
}

type CreateInvoiceParams struct {
	Description     string `json:"description,omitempty"`
	DescriptionHash string `json:"description_hash,omitempty"`
	Msatoshi        int    `json:"msatoshi,omitempty"`
	Preimage        string `json:"preimage,omitempty"`
}

type CreateInvoiceResult struct {
	Invoice     string `json:"invoice"`
	PaymentHash string `json:"payment_hash"`
	HintsCount  int    `json:"hints_count"`
}

type PayInvoiceParams struct {
	Invoice  string `json:"invoice"`
	Msatoshi int    `json:"msatoshi,omitempty"`
}

type PayInvoiceResult struct {
	Sent        bool   `json:"sent"`
	Payee       string `json:"payee"`
	FeeReserve  int    `json:"fee_reserve"`
	PaymentHash string `json:"payment_hash"`
}

type Event struct {
	Event string `json:"event"`
}

type PaymentFailedEvent struct {
	Event
	PaymentHash string `json:"payment_hash"`
	Parts       int    `json:"parts"`
	Failure     string `json:"failure"`
}

type PaymentReceivedEvent struct {
	Event
	PaymentHash string `json:"payment_hash"`
}

type PaymentSucceededEvent struct {
	Event
	PaymentHash string `json:"payment_hash"`
	Parts       int    `json:"parts"`
	FeeMsatoshi int    `json:"fee_msatoshi"`
}
