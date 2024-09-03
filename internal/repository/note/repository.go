package note

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xeeetu/gRPC/internal/repository"
	"github.com/xeeetu/gRPC/internal/repository/note/converter"
	modelRepo "github.com/xeeetu/gRPC/internal/repository/note/model"
	"github.com/xeeetu/gRPC/model"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.NoteRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, "INSERT INTO note (title, content) VALUES ($1, $2) RETURNING id", info.Title, info.Content).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	var note modelRepo.Note
	err := r.db.QueryRow(ctx, "SELECT id, title, content, created_at, updated_at FROM note WHERE id = $1", id).
		Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
