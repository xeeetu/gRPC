package note

import (
	"context"

	"github.com/xeeetu/gRPC/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	id, err := s.noteRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}
	return id, nil
}
