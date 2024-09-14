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
)

func TestList(t *testing.T) {
	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.ListRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("service error")

		limit  = int64(gofakeit.Number(1, 10))
		offset = int64(gofakeit.Number(1, 10))

		req = &desc.ListRequest{
			Limit:  limit,
			Offset: offset,
		}

		res = &desc.ListResponse{
			Notes: []*desc.Note{},
		}

		resModel = []*model.Note{}
	)

	tests := []struct {
		name                string
		args                args
		want                *desc.ListResponse
		err                 error
		noteServiceMockFunc noteServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMockFunc: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)

				mock.ListMock.Expect(ctx, offset, limit).Return(resModel, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			noteServiceMockFunc: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)

				mock.ListMock.Expect(ctx, offset, limit).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.noteServiceMockFunc(mc)
			api := note.NewImplementation(mock)

			resHandler, err := api.List(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
