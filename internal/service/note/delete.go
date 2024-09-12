package note

import (
	"context"
	"fmt"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.noteRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("s.noteRepository.Delete: %w", err)
	}
	return nil
}
