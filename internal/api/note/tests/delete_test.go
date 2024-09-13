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
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelete(t *testing.T) {
	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("service error")

		id = gofakeit.Int64()

		req = &desc.DeleteRequest{
			Id: id,
		}

		res = &emptypb.Empty{}
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
				ctx: context.Background(),
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMockFunc: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)

				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want: res,
			err:  serviceErr,
			noteServiceMockFunc: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)

				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.noteServiceMockFunc(mc)
			api := note.NewImplementation(mock)

			resHandler, err := api.Delete(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
