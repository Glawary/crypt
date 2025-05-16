package usecase

import (
	"context"
	"time"

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
		"cc.cryptocurrency_create_timestamp as cryptocurrency_create_timestamp",
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
	var timestamp, timestampNew time.Time
	for rows.Next() {
		crypt = &model.Crypto{}
		dataCrypto = &model.DataCrypto{}
		err = rows.Scan(
			&cryptocurrencyNewId,
			&tickerNew,
			&dataCrypto.CryptoexchangeName,
			&dataCrypto.DataOhlcv,
			&dataCrypto.DataOrderBook,
			&timestampNew,
		)
		if err != nil {
			return nil, err
		}
		if cryptocurrencyId != cryptocurrencyNewId {
			if cryptocurrencyId != 0 {
				crypt.CryptocurrencyId = cryptocurrencyId
				crypt.CryptocurrencyTicker = ticker
				crypt.CryptocurrencyCreateTimestamp = timestamp
				crypt.Data = dataArr
				res = append(res, crypt)
			}
			dataArr = []*model.DataCrypto{}
		}
		cryptocurrencyId = cryptocurrencyNewId
		ticker = tickerNew
		timestamp = timestampNew
		dataArr = append(dataArr, dataCrypto)
	}
	crypt = &model.Crypto{}
	crypt.CryptocurrencyId = cryptocurrencyId
	crypt.CryptocurrencyTicker = ticker
	crypt.CryptocurrencyCreateTimestamp = timestamp
	crypt.Data = dataArr
	res = append(res, crypt)
	return res, nil
}

func (rec *CryptService) applyFilter(qu sq.SelectBuilder, filter *model.Filter) sq.SelectBuilder {
	if filter == nil {
		return qu
	}

	if filter.CryptexchangeName != "" {
		qu = qu.Where(sq.Eq{"ce.cryptoexchange_name": filter.CryptexchangeName})
	}
	return qu
}
