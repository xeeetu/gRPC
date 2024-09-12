package note

import (
	"context"

	"github.com/xeeetu/gRPC/converter"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.noteService.Update(ctx, converter.ToUpdateNoteFromDesc(req.GetId(), req.GetInfo()))
	return &emptypb.Empty{}, err
}
