package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Glawary/crypt/generated"
	"github.com/Glawary/crypt/internal/usecase"
	"github.com/Glawary/crypt/internal/usecase/model"
)

type CryptServer struct {
	pb.UnimplementedCryptoServiceServer
	cryptService *usecase.CryptService
}

func NewCryptServer(cryptService *usecase.CryptService) *CryptServer {
	return &CryptServer{
		cryptService: cryptService,
	}
}

func (rec *CryptServer) ListCryptoCurrencies(ctx context.Context, req *pb.ListCryptoCurrenciesRequest) (*pb.ListCryptoCurrenciesResponse, error) {
	res, err := rec.cryptService.ListCryptoCurrency(ctx, &model.Filter{CryptexchangeName: req.GetFilter().GetCryptoexchangeName()})
	if err != nil {
		return nil, err
	}
	currencies := make([]*pb.CryptoCurrency, 0, len(res))
	for _, val := range res {
		data := make([]*pb.CryptoCurrencyInfo, 0, len(val.Data))
		for _, info := range val.Data {
			data = append(data, &pb.CryptoCurrencyInfo{
				CryptoexchangeName: info.CryptoexchangeName,
				Ohlcv:              info.DataOhlcv,
			})
		}
		currencies = append(currencies, &pb.CryptoCurrency{
			CryptocurrencyId:              int32(val.CryptocurrencyId),
			CryptocurrencyTicker:          val.CryptocurrencyTicker,
			CryptocurrencyCreateTimestamp: timestamppb.New(val.CryptocurrencyCreateTimestamp),
			Data:                          data,
		})
	}
	return &pb.ListCryptoCurrenciesResponse{Currencies: currencies}, nil
}
