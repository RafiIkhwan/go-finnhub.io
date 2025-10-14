package entity

type RealTimeTrade struct {
	Type string         `json:"type"`
	Data []TradeDetails `json:"data"`
}

type TradeDetails struct {
	Symbol    string   `json:"s"`
	Price     float64  `json:"p"`
	Timestamp int64    `json:"t"`
	Volume    float64  `json:"v"`
	Conditions []string `json:"c,omitempty"`
}