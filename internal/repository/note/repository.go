package note

import (
	"context"
	"fmt"
	"strings"

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
		return 0, fmt.Errorf("reportRepo.Create: %w", err)
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	var note modelRepo.Note
	err := r.db.QueryRow(ctx, "SELECT id, title, content, created_at, updated_at FROM note WHERE id = $1", id).
		Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("reportRepo.Get,id-%d: %w", id, err)
	}

	return converter.ToNoteFromRepo(&note), nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM note WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("reportRepo.Delete, id-%d: %w", id, err)
	}
	return nil
}

func (r *repo) List(ctx context.Context, offset int64, limit int64) ([]*model.Note, error) {
	res, err := r.db.Query(ctx, "SELECT id, title, content, created_at, updated_at FROM note WHERE id >= $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("reportRepo.List, limit-%d, offset-%d: %w", limit, offset, err)
	}
	defer res.Close()

	notes := make([]*model.Note, 0)

	for res.Next() {
		var note model.Note
		err = res.Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("reportRepo.List: %w", err)
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func (r *repo) Update(ctx context.Context, updateInfo *model.UpdateNote) error {
	query := "UPDATE note SET"
	values := []any{}
	numberParam := 1

	if updateInfo.Info.Title != nil {
		query = fmt.Sprintf("%s title = $%d,", query, numberParam)
		values = append(values, *updateInfo.Info.Title)
		numberParam++
	}
	if updateInfo.Info.Content != nil {
		query = fmt.Sprintf("%s content = $%d,", query, numberParam)
		values = append(values, *updateInfo.Info.Content)
		numberParam++
	}

	if numberParam > 1 {
		query = fmt.Sprintf("%s updated_at = now(),", query)
	}

	query = strings.TrimSuffix(query, ",")
	query = fmt.Sprintf("%s WHERE id = $%d", query, numberParam)

	values = append(values, updateInfo.ID)

	_, err := r.db.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("reportRepo.Update: %w", err)
	}
	return nil
}
