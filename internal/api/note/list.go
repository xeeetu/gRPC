package note

import (
	"context"

	"github.com/xeeetu/gRPC/converter"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
)

func (i *Implementation) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	notes, err := i.noteService.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	descNotes := make([]*desc.Note, 0, len(notes))
	for _, note := range notes {
		descNotes = append(descNotes, converter.ToNoteFromService(note))
	}
	return &desc.ListResponse{Notes: descNotes}, nil
}
