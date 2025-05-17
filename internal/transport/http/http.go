package http

import (
	"github.com/go-chi/chi/v5"

	swagger "github.com/swaggo/http-swagger/v2"

	"github.com/Glawary/crypt/internal/swag"
	"github.com/Glawary/crypt/internal/usecase"
	"github.com/Glawary/crypt/pkg/http"
)

type Server struct {
	cryptService *usecase.CryptService
}

func InitServer(cfg *http.HttpConfig, cryptService *usecase.CryptService) (*http.HTTPServer, error) {
	s := &Server{
		cryptService: cryptService,
	}

	r := chi.NewRouter()
	r.Route("/api/v1", func(router chi.Router) {
		router.Get("/swagger/*", swagger.Handler(func(cfg *swagger.Config) {
			cfg.InstanceName = swag.SwaggerInfo.InstanceName()
		}))
		router.Get("/list", s.ListCrypto)
	})

	server := http.NewHTTPServer(cfg, r)
	return server, nil
}
