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
)

func TestDelete(t *testing.T) {
	//t.Parallel()
	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx context.Context
		id  int64
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr       = fmt.Errorf("repository error")
		repoErrInServ = fmt.Errorf("s.noteRepository.Delete: %w", repoErr)
	)

	tests := []struct {
		name                   string
		args                   args
		err                    error
		noteRepositoryMockFunc noteRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErrInServ,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
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

			err := serv.Delete(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
		})
	}
}
