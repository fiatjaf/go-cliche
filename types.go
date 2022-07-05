package cliche

import "encoding/json"

type JSONRPCRequest struct {
	Id     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type JSONRPCResponse struct {
	Id     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type JSONRPCNotification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type GetInfoResult struct {
	BlockHeight int           `json:"block_height"`
	Channels    []ChannelInfo `json:"channels"`
	FeeRates    struct {
		Num1   int `json:"1"`
		Num10  int `json:"10"`
		Num100 int `json:"100"`
	} `json:"fee_rates"`
	FiatRates struct {
		USD float64 `json:"usd"`
	} `json:"fiat_rates"`
	KnownChannels struct {
		Hosted int `json:"hosted"`
		Normal int `json:"normal"`
	} `json:"known_channels"`
	MainPubkey       string `json:"main_pubkey"`
	OutgoingPayments []struct {
		Hash  string `json:"hash"`
		Htlcs []struct {
			Channel  string `json:"channel"`
			Expiry   int    `json:"expiry"`
			ID       int    `json:"id"`
			Msatoshi int64  `json:"msatoshi"`
		} `json:"htlcs"`
	} `json:"outgoing_payments"`
	Wallets []struct {
		Balance int64  `json:"balance"`
		Label   string `json:"label"`
	} `json:"wallets"`
}

type ChannelInfo struct {
	ID             string `json:"id"`
	ShortChannelId string `json:"short_channel_id"`
	Balance        int    `json:"balance"`
	CanReceive     int64  `json:"can_receive"`
	CanSend        int64  `json:"can_send"`
	HostedChannel  struct {
		OverrideProposal struct {
			OurBalance   int64 `json:"our_balance"`
			TheirBalance int64 `json:"their_balance"`
		} `json:"override_proposal"`
		ResizeProposal int64 `json:"resize_proposal"`
	} `json:"hosted_channel"`
	Inflight struct {
		Incoming int `json:"incoming"`
		Outgoing int `json:"outgoing"`
		Revealed int `json:"revealed"`
	} `json:"inflight"`
	Peer struct {
		Addr      string `json:"addr"`
		OurPubkey string `json:"our_pubkey"`
		Pubkey    string `json:"pubkey"`
	} `json:"peer"`
	Policy struct {
		BaseFee         int64 `json:"base_fee"`
		CltvDelta       int   `json:"cltv_delta"`
		FeePerMillionth int   `json:"fee_per_millionth"`
		HtlcMax         int64 `json:"htlc_max"`
		HtlcMin         int64 `json:"htlc_min"`
	} `json:"policy,omitempty"`
	Status string `json:"status"`
}

type CreateInvoiceParams struct {
	Description     string `json:"description,omitempty"`
	DescriptionHash string `json:"description_hash,omitempty"`
	Msatoshi        int64  `json:"msatoshi,omitempty"`
	Preimage        string `json:"preimage,omitempty"`
	Label           string `json:"label,omitempty"`
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

type PaymentInfo struct {
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

type CheckPaymentResult struct {
	PaymentInfo
}

type ListPaymentsResult []PaymentInfo

type PaymentFailedEvent struct {
	PaymentHash string   `json:"payment_hash"`
	Parts       int      `json:"parts"`
	Failure     []string `json:"failure"`
}

type PaymentReceivedEvent struct {
	PaymentHash string `json:"payment_hash"`
	Preimage    string `json:"preimage"`
	Msatoshi    int64  `json:"msatoshi"`
}

type PaymentSucceededEvent struct {
	PaymentHash string `json:"payment_hash"`
	FeeMsatoshi int64  `json:"fee_msatoshi"`
	Msatoshi    int64  `json:"msatoshi"`
	Preimage    string `json:"preimage"`
	Parts       int    `json:"parts"`
}
