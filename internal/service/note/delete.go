package note

import "context"

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.noteRepository.Delete(ctx, id)
	return err
}
