package controller

type Pair struct {
	TokenA string
	TokenB string
}

type Rate struct {
	Pair
	Bid       string
	Ask       string
	Timestamp int64
}
