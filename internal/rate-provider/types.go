package rateprovider

type OrderBook struct {
	Timestamp int64
	Asks      []Order
	Bids      []Order
}

type Order struct {
	Price string
}
