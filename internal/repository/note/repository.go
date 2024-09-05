package note

import (
	"context"
	"fmt"
	"strconv"
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

func (r *repo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM note WHERE id = $1", id)
	return err
}

func (r *repo) List(ctx context.Context, offset int64, limit int64) ([]*model.Note, error) {
	res, err := r.db.Query(ctx, "SELECT id, title, content, created_at, updated_at FROM note WHERE id >= $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	notes := make([]*model.Note, 0)

	for res.Next() {
		var note model.Note
		err = res.Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func (r *repo) Update(ctx context.Context, updateInfo *model.UpdateNote) error {
	query := "UPDATE note SET"
	values := []interface{}{}
	numberParam := 1

	if updateInfo.Info.Title != nil {
		query += " title = $" + strconv.Itoa(numberParam) + ","
		values = append(values, *updateInfo.Info.Title)
		numberParam++
	}
	if updateInfo.Info.Content != nil {
		query += " content = $" + strconv.Itoa(numberParam) + ","
		values = append(values, *updateInfo.Info.Content)
		numberParam++
	}

	if numberParam > 1 {
		query += " updated_at = now(),"
	}

	query = strings.TrimSuffix(query, ",") + " WHERE id = $" + strconv.Itoa(numberParam)
	values = append(values, updateInfo.ID)

	fmt.Println(query)
	fmt.Println(values)
	_, err := r.db.Exec(ctx, query, values...)
	return err
}
