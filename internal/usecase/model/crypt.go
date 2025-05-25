package model

type Crypto struct {
	CryptocurrencyId     int    `db:"cryptocurrency_id"`
	CryptocurrencyTicker string `db:"cryptocurrency_ticker"`
	Data                 []*DataCrypto
}

type DataCrypto struct {
	CryptoExchangeName string `db:"cryptoexchange_name"`
	DataOhlcv          []byte `db:"data_olhcv" json:"data_olhcv"`
}

type Filter struct {
	CryptoExchangeName string  `json:"cryptoexchange_name"`
	PriceFrom          float64 `json:"price_from"`
	PriceTo            float64 `json:"price_to"`
	FindBrush          bool    `json:"find_brush"`
}
