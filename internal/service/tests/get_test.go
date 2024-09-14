package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/xeeetu/gRPC/internal/repository"
	repositoryMocks "github.com/xeeetu/gRPC/internal/repository/mocks"
	noteService "github.com/xeeetu/gRPC/internal/service/note"
	"github.com/xeeetu/gRPC/model"
)

func TestGet(t *testing.T) {
	//t.Parallel()
	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx context.Context
		id  int64
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		title     = gofakeit.Name()
		content   = gofakeit.Name()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		repoErr       = fmt.Errorf("repository error")
		repoErrInServ = fmt.Errorf("s.noteRepository.Get: %w", repoErr)

		modelInfo = model.NoteInfo{
			Title:   title,
			Content: content,
		}

		modelNote = &model.Note{
			ID:        id,
			Info:      modelInfo,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)

	tests := []struct {
		name                   string
		args                   args
		want                   *model.Note
		err                    error
		noteRepositoryMockFunc noteRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: modelNote,
			err:  nil,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.GetMock.Expect(ctx, id).Return(modelNote, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErrInServ,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()

			mock := tt.noteRepositoryMockFunc(mc)
			serv := noteService.NewService(mock)

			resHandler, err := serv.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
