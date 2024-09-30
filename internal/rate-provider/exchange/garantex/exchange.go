package garantex

import (
	"encoding/json"
	"fmt"
	"net/http"

	rateprovider "github.com/irbgeo/usdt-rate/internal/rate-provider"
	rateerr "github.com/irbgeo/usdt-rate/internal/utils/rate-error"
)

const (
	baseURL      = "https://garantex.org/api/v2"
	orderBookURL = baseURL + "/depth"
)

type exchange struct {
	cli *http.Client
}

func NewClient() *exchange {
	return &exchange{
		cli: &http.Client{},
	}
}

func (s *exchange) OrderBook(market string, depth int) (*rateprovider.OrderBook, error) {
	ob, err := s.orderBook(market)
	if err != nil {
		return nil, rateerr.New(err, "msg", "failed to get order book")
	}

	return toInternalOrderBook(ob, depth)
}

func (s *exchange) orderBook(market string) (*orderBook, error) {
	req, err := http.NewRequest("GET", orderBookURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("market", market)
	req.URL.RawQuery = q.Encode()

	resp, err := s.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var orderBook orderBook
	if err := json.NewDecoder(resp.Body).Decode(&orderBook); err != nil {
		return nil, rateerr.New(err, "msg", "failed to decode response")
	}

	return &orderBook, nil
}
