package note

import (
	"github.com/xeeetu/gRPC/internal/repository"
	def "github.com/xeeetu/gRPC/internal/service"
)

var _ def.NoteService = (*serv)(nil)

type serv struct {
	noteRepository repository.NoteRepository
}

func NewService(noteRepository repository.NoteRepository) *serv {
	return &serv{
		noteRepository: noteRepository,
	}
}
