package model

type Crypto struct {
	CryptocurrencyId     int    `db:"cryptocurrency_id"`
	CryptocurrencyTicker string `db:"cryptocurrency_ticker"`
	Data                 []*DataCrypto
}

type DataCrypto struct {
	CryptoExchangeName string  `db:"cryptoexchange_name"`
	DataOlhcv          []byte  `db:"data_olhcv" json:"data_olhcv"`
	DataOrderBook      []byte  `db:"data_order_book" json:"data_order_book,omitempty"`
	Last               float64 `db:"last" json:"last"`
	Spread             float64 `db:"spread" json:"spread"`
}

type DataOrderBook struct {
	Bids [][]float64 `db:"bids" json:"bids"`
	Asks [][]float64 `db:"asks" json:"asks"`
}

type Filter struct {
	CryptoExchangeName string  `json:"cryptoexchange_name"`
	PriceFrom          float64 `json:"price_from"`
	PriceTo            float64 `json:"price_to"`
	FindBrush          bool    `json:"find_brush"`
}
