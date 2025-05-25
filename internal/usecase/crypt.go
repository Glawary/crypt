package usecase

import (
	"context"
	"encoding/json"
	"math"

	sq "github.com/Masterminds/squirrel"

	"github.com/Glawary/crypt/internal/client"
	"github.com/Glawary/crypt/internal/usecase/model"
	"github.com/Glawary/crypt/pkg/postgres"
)

type CryptService struct {
	db *postgres.Instance
}

func NewCryptService() *CryptService {
	return &CryptService{
		db: client.GetDBInstance(),
	}
}

func (rec *CryptService) ListCryptoCurrency(ctx context.Context, filter *model.Filter) ([]*model.Crypto, error) {
	qu := sq.Select(
		"cd.cryptocurrency_id as cryptocurrency_id",
		"cc.cryptocurrency_ticker as cryptocurrency_ticker",
		"ce.cryptoexchange_name as cryptoexchange_name",
		"cd.data_olhcv as data_olhcv",
		"cd.data_order_book as data_order_book",
	).From("cryptocurrency_data as cd").
		LeftJoin("cryptocurrency as cc ON cd.cryptocurrency_id = cc.cryptocurrency_id").
		LeftJoin("cryptoexchange as ce ON ce.cryptoexchange_id = cd.cryptoexchange_id")

	qu = rec.applyFilter(qu, filter)
	qu = qu.OrderBy("cd.cryptocurrency_id", "ce.cryptoexchange_id")
	quSql, args, err := qu.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := rec.db.GetSqlxDB().QueryContext(ctx, quSql, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var res []*model.Crypto
	var crypt *model.Crypto
	var dataArr []*model.DataCrypto
	var dataCrypto *model.DataCrypto
	var cryptocurrencyId, cryptocurrencyNewId int
	var ticker, tickerNew string
	for rows.Next() {
		crypt = &model.Crypto{}
		dataCrypto = &model.DataCrypto{}
		err = rows.Scan(
			&cryptocurrencyNewId,
			&tickerNew,
			&dataCrypto.CryptoExchangeName,
			&dataCrypto.DataOlhcv,
			&dataCrypto.DataOrderBook,
		)
		if err != nil {
			return nil, err
		}
		if cryptocurrencyId != cryptocurrencyNewId {
			if cryptocurrencyId != 0 {
				crypt.CryptocurrencyId = cryptocurrencyId
				crypt.CryptocurrencyTicker = ticker
				crypt.Data = dataArr
				res = append(res, crypt)
			}
			dataArr = []*model.DataCrypto{}
		}
		cryptocurrencyId = cryptocurrencyNewId
		ticker = tickerNew
		dataArr = append(dataArr, dataCrypto)
	}
	crypt = &model.Crypto{}
	crypt.CryptocurrencyId = cryptocurrencyId
	crypt.CryptocurrencyTicker = ticker
	crypt.Data = dataArr
	res = append(res, crypt)
	return filterResult(res, filter), nil
}

func (rec *CryptService) applyFilter(qu sq.SelectBuilder, filter *model.Filter) sq.SelectBuilder {
	if filter == nil {
		return qu
	}

	if filter.CryptoExchangeName != "" {
		qu = qu.Where(sq.Eq{"ce.cryptoexchange_name": filter.CryptoExchangeName})
	}
	return qu
}

func filterResult(res []*model.Crypto, filter *model.Filter) []*model.Crypto {
	out := make([]*model.Crypto, 0, len(res))
	var dataCrypto []*model.DataCrypto
	for _, item := range res {
		dataCrypto = make([]*model.DataCrypto, 0, len(item.Data))
		for _, data := range item.Data {
			var olhcv [][]float64
			_ = json.Unmarshal([]byte(data.DataOlhcv), &olhcv)
			if !(len(olhcv) == 0 || (filter.PriceFrom > 0 && olhcv[len(olhcv)-1][4] < filter.PriceFrom) ||
				(filter.PriceTo > 0 && olhcv[len(olhcv)-1][4] > filter.PriceTo) || (filter.FindBrush && !detectBrush(olhcv))) {
				var orderBook *model.DataOrderBook
				_ = json.Unmarshal(data.DataOrderBook, &orderBook)
				if orderBook != nil && len(orderBook.Bids) > 0 && len(orderBook.Asks) > 0 {
					spread := math.Round(((float64(orderBook.Asks[0][0])-float64(orderBook.Bids[0][0]))/float64(orderBook.Bids[0][0]))*100*1000) / 1000
					if spread > 0 {
						data.Spread = spread
					} else {
						data.Spread = 0
					}
				}
				data.Last = olhcv[len(olhcv)-1][4]
				data.DataOrderBook = nil
				dataCrypto = append(dataCrypto, data)
			}
		}
		if len(dataCrypto) > 0 {
			item.Data = dataCrypto
			out = append(out, item)
		}
	}
	return out
}

func detectBrush(olhcv [][]float64) bool {
	var sumClose float64
	if len(olhcv) == 0 {
		return false
	}
	closes := make([]float64, 0, len(olhcv))
	for _, item := range olhcv {
		closes = append(closes, item[4])
		sumClose += item[4]
	}

	var sumVolume float64
	volumes := make([]float64, 0, len(olhcv))
	for _, item := range olhcv {
		volumes = append(volumes, item[5]*item[4])
		sumVolume += item[5]
	}

	avgClose := sumClose / float64(len(olhcv))
	upperBound := avgClose * (1.01)
	lowerBound := avgClose * (0.99)

	outliers := 0
	for _, item := range closes {
		if !(lowerBound <= item && item <= upperBound) {
			outliers += 1
		}
	}
	if outliers > 2 {
		return false
	}

	avgVolume := sumVolume / float64(len(volumes))
	spikeCount := 0
	for _, item := range volumes {
		if item > avgVolume*2 {
			spikeCount += 1
		}
	}
	if spikeCount > 3 {
		return false
	}

	lastClose := closes[len(closes)-1]
	if !(lowerBound <= lastClose && lastClose <= upperBound) {
		return false
	}
	return true
}
