package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/xeeetu/gRPC/internal/api/note"
	"github.com/xeeetu/gRPC/internal/service"
	serviceMocks "github.com/xeeetu/gRPC/internal/service/mocks"
	"github.com/xeeetu/gRPC/model"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUpdate(t *testing.T) {
	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("service error")

		id      = gofakeit.Int64()
		title   = gofakeit.Name()
		content = gofakeit.Name()

		updateInfo = &desc.UpdateNoteInfo{
			Title:   wrapperspb.String(title),
			Content: wrapperspb.String(content),
		}

		req = &desc.UpdateRequest{
			Id:   id,
			Info: updateInfo,
		}

		modelUpdateInfo = &model.UpdateNoteInfo{
			Title:   &title,
			Content: &content,
		}

		modelUpdate = &model.UpdateNote{
			ID:   id,
			Info: modelUpdateInfo,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *emptypb.Empty
		err                 error
		noteServiceMockFunc noteServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			noteServiceMockFunc: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.UpdateMock.Expect(ctx, modelUpdate).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  serviceErr,
			noteServiceMockFunc: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.UpdateMock.Expect(ctx, modelUpdate).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.noteServiceMockFunc(mc)
			api := note.NewImplementation(mock)

			resHandler, err := api.Update(ctx, req)
			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
