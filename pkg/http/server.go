package http

import (
	"context"
	"log"
	"net"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	cfg    *HttpConfig
	notify chan error
}

func NewHTTPServer(cfg *HttpConfig, handler http.Handler) *HTTPServer {
	httpServer := &HTTPServer{
		server: &http.Server{
			Addr:    cfg.Url,
			Handler: handler,
		},
		cfg:    cfg,
		notify: make(chan error, 1),
	}

	httpServer.Start()

	return httpServer
}

func (rec *HTTPServer) GetNotify() <-chan error {
	return rec.notify
}

func (rec *HTTPServer) GetGRPCServer() *http.Server {
	return rec.server
}

func (rec *HTTPServer) Start() {
	go func() {
		listen, err := net.Listen("tcp", rec.cfg.Url)
		if err != nil {
			log.Fatalf("failed to establish HTTP connection: %v", err)
		}
		log.Printf("RUN HTTP SERVER on %s", rec.cfg.Url)
		rec.notify <- rec.server.Serve(listen)
	}()
}

func (rec *HTTPServer) Shutdown() {
	_ = rec.server.Shutdown(context.Background())
}
