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
		Balance int64  `json:"balance"`
	} `json:"wallets"`
	Channels []struct {
		ID      string `json:"id"`
		Balance int64  `json:"balance"`
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
	Msatoshi        int64  `json:"msatoshi,omitempty"`
	Preimage        string `json:"preimage,omitempty"`
}

type CreateInvoiceResult struct {
	Invoice     string `json:"invoice"`
	PaymentHash string `json:"payment_hash"`
}

type PayInvoiceParams struct {
	Invoice  string `json:"invoice"`
	Msatoshi int64  `json:"msatoshi,omitempty"`
}

type PayInvoiceResult struct {
	Sent        bool   `json:"sent"`
	Payee       string `json:"payee"`
	FeeReserve  int    `json:"fee_reserve"`
	PaymentHash string `json:"payment_hash"`
}

type CheckPaymentResult struct {
	Status      string `json:"status"`
	SeenAt      int64  `json:"seen_at"`
	Invoice     string `json:"invoice"`
	Preimage    string `json:"preimage"`
	Msatoshi    int64  `json:"msatoshi"`
	UpdatedAt   int64  `json:"updated_at"`
	IsIncoming  bool   `json:"is_incoming"`
	FeeMsatoshi int64  `json:"fee_msatoshi"`
	PaymentHash string `json:"payment_hash"`
}

type Event struct {
	Event string `json:"event"`
}

type PaymentFailedEvent struct {
	Event
	PaymentHash string   `json:"payment_hash"`
	Parts       int      `json:"parts"`
	Failure     []string `json:"failure"`
}

type PaymentReceivedEvent struct {
	Event
	PaymentHash string `json:"payment_hash"`
	Preimage    string `json:"preimage"`
	Msatoshi    int64  `json:"msatoshi"`
}

type PaymentSucceededEvent struct {
	Event
	PaymentHash string `json:"payment_hash"`
	FeeMsatoshi int64  `json:"fee_msatoshi"`
	Msatoshi    int64  `json:"msatoshi"`
	Preimage    string `json:"preimage"`
	Parts       int    `json:"parts"`
}
