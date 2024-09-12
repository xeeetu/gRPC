package note

import (
	"context"

	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.noteService.Delete(ctx, req.GetId())
	return &emptypb.Empty{}, err
}
