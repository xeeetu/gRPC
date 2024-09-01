package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xeeetu/gRPC/internal/config"
	"github.com/xeeetu/gRPC/internal/config/env"
	"github.com/xeeetu/gRPC/internal/repository"
	"github.com/xeeetu/gRPC/internal/repository/note"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	noteRepository repository.NoteRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.noteRepository.Create(ctx, req.GetInfo())
	if err != nil {
		return nil, err
	}

	log.Printf("created note with id %d", id)

	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := s.noteRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, title: %s, content: %s, created_at: %v, updated_at: %v\n", noteObj.GetId(), noteObj.Info.Title, noteObj.Info.Content, noteObj.CreatedAt, noteObj.UpdatedAt)

	return &desc.GetResponse{Note: noteObj}, nil
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

	noteRepo := note.NewRepository(pool)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{noteRepository: noteRepo})

	log.Printf("server listening on %s", grpcConfig.Address())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
