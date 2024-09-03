package converter

import (
	modelRepo "github.com/xeeetu/gRPC/internal/repository/note/model"
	"github.com/xeeetu/gRPC/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {

	return &model.Note{
		ID:        note.ID,
		Info:      ToNoteInfoFromRepo(&note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNoteInfoFromRepo(noteInfo *modelRepo.Info) model.NoteInfo {
	return model.NoteInfo{
		Title:   noteInfo.Title,
		Content: noteInfo.Content,
	}
}
