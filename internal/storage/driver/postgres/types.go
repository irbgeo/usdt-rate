package postgres

type StartOpts struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

type rate struct {
	TokenA    string `db:"token_a"`
	TokenB    string `db:"token_b"`
	Bid       string `db:"bid"`
	Ask       string `db:"ask"`
	Timestamp int64  `db:"timestamp"`
}
