package grpc

type GRPCConfig struct {
	Url string `env:"GRPC_URL" envDefault:"localhost:8080"`
}
