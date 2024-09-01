package repository

import (
	"context"

	desc "github.com/xeeetu/gRPC/pkg/note_v1"
)

type NoteRepository interface {
	Create(ctx context.Context, info *desc.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.Note, error)
}
