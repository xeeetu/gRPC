package note

import (
	"context"

	"github.com/xeeetu/gRPC/model"
)

func (s *serv) Update(ctx context.Context, info *model.UpdateNote) error {
	return s.noteRepository.Update(ctx, info)
}
