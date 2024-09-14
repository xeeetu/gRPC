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

func TestUpdate(t *testing.T) {
	//t.Parallel()
	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx        context.Context
		updateNote *model.UpdateNote
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.Name()
		content = gofakeit.Name()

		repoErr       = fmt.Errorf("repository error")
		repoErrInServ = fmt.Errorf("s.noteRepository.Update: %w", repoErr)

		modelUpdateInfo = &model.UpdateNoteInfo{
			Title:   &title,
			Content: &content,
		}

		modelUpdateNote = &model.UpdateNote{
			ID:   id,
			Info: modelUpdateInfo,
		}
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
				ctx:        ctx,
				updateNote: modelUpdateNote,
			},
			err: nil,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.UpdateMock.Expect(ctx, modelUpdateNote).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:        ctx,
				updateNote: modelUpdateNote,
			},
			err: repoErrInServ,
			noteRepositoryMockFunc: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repositoryMocks.NewNoteRepositoryMock(mc)

				mock.UpdateMock.Expect(ctx, modelUpdateNote).Return(repoErr)
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

			err := serv.Update(tt.args.ctx, tt.args.updateNote)
			require.Equal(t, tt.err, err)
		})
	}
}
