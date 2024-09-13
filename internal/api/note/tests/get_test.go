package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/xeeetu/gRPC/internal/api/note"
	"github.com/xeeetu/gRPC/internal/service"
	"github.com/xeeetu/gRPC/model"

	serviceMocks "github.com/xeeetu/gRPC/internal/service/mocks"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {
	type noteServiceMock func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("service error")

		id        = gofakeit.Int64()
		title     = gofakeit.Word()
		content   = gofakeit.Word()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		req = &desc.GetRequest{
			Id: id,
		}

		info = &desc.NoteInfo{
			Title:   title,
			Content: content,
		}

		res = &desc.GetResponse{
			Note: &desc.Note{
				Id:        id,
				Info:      info,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}

		infoRes = model.NoteInfo{
			Title:   title,
			Content: content,
		}
		noteRes = &model.Note{
			ID:        id,
			Info:      infoRes,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *desc.GetResponse
		err                 error
		noteServiceMockFunc noteServiceMock
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
				mock.GetMock.Expect(ctx, id).Return(noteRes, nil)
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
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.noteServiceMockFunc(mc)
			api := note.NewImplementation(mock)

			resHandler, err := api.Get(ctx, req)
			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
