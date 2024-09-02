package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xeeetu/gRPC/internal/config"
	"github.com/xeeetu/gRPC/internal/config/env"

	noteAPI "github.com/xeeetu/gRPC/internal/api/note"
	noteRepository "github.com/xeeetu/gRPC/internal/repository/note"
	noteService "github.com/xeeetu/gRPC/internal/service/note"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	noteRepo := noteRepository.NewRepository(pool)
	noteSrv := noteService.NewService(noteRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, noteAPI.NewImplementation(noteSrv))

	log.Printf("server listening on %s", grpcConfig.Address())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
