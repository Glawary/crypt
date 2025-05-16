package grpc

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	cfg    *GRPCConfig
	notify chan error
}

type RegistrarGRPCFunc func(*GRPCServer) error

func NewGRPCServer(cfg *GRPCConfig) *GRPCServer {
	grpcServer := &GRPCServer{
		server: grpc.NewServer(
			grpc.MaxSendMsgSize(20*1024*1024),
			grpc.MaxRecvMsgSize(20*1024*1024),
		),
		cfg:    cfg,
		notify: make(chan error, 1),
	}

	grpcServer.Start()

	return grpcServer
}

func (rec *GRPCServer) GetNotify() <-chan error {
	return rec.notify
}

func (rec *GRPCServer) GetGRPCServer() *grpc.Server {
	return rec.server
}

func (rec *GRPCServer) Start() {
	go func() {
		listen, err := net.Listen("tcp", rec.cfg.Url)
		if err != nil {
			log.Fatalf("failed to establish GRPC connection: %v", err)
		}
		log.Printf("RUN GRPC SERVER on %s", rec.cfg.Url)
		rec.notify <- rec.server.Serve(listen)
	}()
}

func (rec *GRPCServer) Shutdown() {
	if rec.server == nil {
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		rec.server.GracefulStop()
		close(quit)
	}()

	select {
	case <-ticker.C:
		rec.server.Stop()
	case <-quit:
		ticker.Stop()
	}
}
