package note

import (
	"context"
	"fmt"

	"github.com/xeeetu/gRPC/model"
)

func (s *serv) Update(ctx context.Context, info *model.UpdateNote) error {
	err := s.noteRepository.Update(ctx, info)
	if err != nil {
		return fmt.Errorf("s.noteRepository.Update: %w", err)
	}
	return nil
}
