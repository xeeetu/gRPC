package note

import (
	"context"
	"log"

	"github.com/xeeetu/gRPC/converter"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := i.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("returning note: id: %d, title: %s, content: %s, createAt: %v, updatedAt: %v", noteObj.ID, noteObj.Info.Title, noteObj.Info.Content, noteObj.CreatedAt, noteObj.UpdatedAt)

	return &desc.GetResponse{Note: converter.ToNoteFromService(noteObj)}, nil
}
