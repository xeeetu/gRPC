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

func TestCreate(t *testing.T) {
	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.Name()
		content = gofakeit.Name()

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.NoteInfo{
				Title:   title,
				Content: content,
			},
		}

		info = &model.NoteInfo{
			Title:   title,
			Content: content,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *desc.CreateResponse
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
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
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
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteServiceMock := tt.noteServiceMockFunc(mc)
			api := note.NewImplementation(noteServiceMock)

			resHandler, err := api.Create(ctx, req)
			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
