package model

import (
	"time"
)

type Crypto struct {
	CryptocurrencyId              int       `db:"cryptocurrency_id"`
	CryptocurrencyTicker          string    `db:"cryptocurrency_ticker"`
	CryptocurrencyCreateTimestamp time.Time `db:"cryptocurrency_create_timestamp"`
	Data                          []*DataCrypto
}

type DataCrypto struct {
	CryptoexchangeName string `db:"cryptoexchange_name"`
	DataOhlcv          []byte `db:"data_olhcv" json:"data_olhcv"`
	DataOrderBook      []byte `db:"data_order_book" json:"data_order_book"`
}
