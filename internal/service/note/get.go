package note

import (
	"context"
	"fmt"

	"github.com/xeeetu/gRPC/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.Note, error) {
	note, err := s.noteRepository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.noteRepository.Get: %w", err)
	}
	return note, nil
}
