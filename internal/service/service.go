package service

import (
	"context"

	"github.com/xeeetu/gRPC/model"
)

type NoteService interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
