package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Glawary/crypt/internal/client"
	"github.com/Glawary/crypt/internal/config"
	grpc_service "github.com/Glawary/crypt/internal/transport/grpc"
	"github.com/Glawary/crypt/internal/usecase"
)

func Run() {
	cfg, err := config.New("./.env")
	if err != nil {
		log.Printf("error loading config: %v", err)
	}

	_, err = client.InitDB(cfg.DB)
	if err != nil {
		log.Printf("error initializing postgres: %v", err)
	}

	cryptService := usecase.NewCryptService()
	grpcServer, err := grpc_service.InitServer(cfg.GRPCServer, cryptService)
	if err != nil {
		log.Printf("error initializing grpc: %v", err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		log.Println("Got signal:", s)
	case err := <-grpcServer.GetNotify():
		log.Println("Got error:", err)
	}

	grpcServer.Shutdown()
}
