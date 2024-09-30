package garantex

type orderBook struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []order `json:"asks"`
	Bids      []order `json:"bids"`
}

type order struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}
