package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	swagger "github.com/swaggo/http-swagger/v2"

	"github.com/Glawary/crypt/internal/swag"
	"github.com/Glawary/crypt/internal/usecase"
	http2 "github.com/Glawary/crypt/pkg/http"
)

type Server struct {
	cryptService *usecase.CryptService
}

func InitServer(cfg *http2.HttpConfig, cryptService *usecase.CryptService) (*http2.HTTPServer, error) {
	s := &Server{
		cryptService: cryptService,
	}

	r := chi.NewRouter()
	r.Route("/api/v1", func(router chi.Router) {
		router.Use(useCors().Handler)
		router.Get("/swagger/*", swagger.Handler(func(cfg *swagger.Config) {
			cfg.InstanceName = swag.SwaggerInfo.InstanceName()
		}))
		router.Get("/list", s.ListCrypto)
	})

	server := http2.NewHTTPServer(cfg, r)
	return server, nil
}

func useCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		AllowedMethods: []string{
			http.MethodGet,
		},
	})
}
