package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Glawary/crypt/internal/client"
	"github.com/Glawary/crypt/internal/config"
	"github.com/Glawary/crypt/internal/interface"
	grpc_service "github.com/Glawary/crypt/internal/transport/grpc"
	http_service "github.com/Glawary/crypt/internal/transport/http"
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
	var server _interface.Server
	if cfg.GRPCServer.Url != "" {
		server, err = grpc_service.InitServer(cfg.GRPCServer, cryptService)
		if err != nil {
			log.Printf("error initializing http: %v", err)
		}
	} else if cfg.HttpServer.Url != "" {
		server, err = http_service.InitServer(cfg.HttpServer, cryptService)
		if err != nil {
			log.Printf("error initializing http: %v", err)
		}
	}
	if err != nil {
		log.Printf("error initializing grpc: %v", err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		log.Println("Got signal:", s)
	case err := <-server.GetNotify():
		log.Println("Got error:", err)
	}
	server.Shutdown()
}
