package config

import (
	"errors"
	"net"
	"os"
)

var _ GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if host == "" {
		return nil, errors.New("grpc host not found")
	}
	port := os.Getenv(grpcPortEnvName)
	if port == "" {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{host: host, port: port}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
