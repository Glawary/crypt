package grpc

import (
	pb "github.com/Glawary/crypt/generated"
	"github.com/Glawary/crypt/internal/usecase"
	"github.com/Glawary/crypt/pkg/grpc"
)

type Server struct{}

func InitServer(cfg *grpc.GRPCConfig, cryptService *usecase.CryptService) (*grpc.GRPCServer, error) {
	server := grpc.NewGRPCServer(cfg)
	pb.RegisterCryptoServiceServer(server.GetGRPCServer(), NewCryptServer(cryptService))
	return server, nil
}
