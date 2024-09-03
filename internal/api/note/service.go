package note

import (
	"github.com/xeeetu/gRPC/internal/service"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

func NewImplementation(noteService service.NoteService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
