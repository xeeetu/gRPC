package tests

import (
	"context"
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

func TestCreate(t *testing.T) {
	//t.Parallel()
	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx       context.Context
		modelInfo *model.NoteInfo
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.Name()
		content = gofakeit.Name()

		repoErr       = fmt.Errorf("repository error")
		repoErrInServ = fmt.Errorf("s.noteRepository.Create: %w", repoErr)

		modelInfo = &model.NoteInfo{
			Title:   title,
			Content: content,
		}
	)

	tests := []struct {
		name                   string
		args                   args
		want                   int64
		err                    error
		noteRepositoryMockFunc noteRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:       ctx,
				modelInfo: modelInfo,
			},
			want: id,
			err:  nil,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.CreateMock.Expect(ctx, modelInfo).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:       ctx,
				modelInfo: modelInfo,
			},
			want: 0,
			err:  repoErrInServ,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.CreateMock.Expect(ctx, modelInfo).Return(0, repoErr)
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

			resHandler, err := serv.Create(tt.args.ctx, tt.args.modelInfo)
			require.Equal(t, tt.want, resHandler)
			require.Equal(t, tt.err, err)
		})
	}
}
