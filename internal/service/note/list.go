package note

import (
	"context"

	"github.com/xeeetu/gRPC/model"
)

func (s *serv) List(ctx context.Context, offset int64, limit int64) ([]*model.Note, error) {
	notes, err := s.noteRepository.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
