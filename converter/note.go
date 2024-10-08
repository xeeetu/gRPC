package converter

import (
	"github.com/xeeetu/gRPC/model"
	desc "github.com/xeeetu/gRPC/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromService(note *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}
	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromService(&note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNoteInfoFromService(noteInfo *model.NoteInfo) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   noteInfo.Title,
		Content: noteInfo.Content,
	}
}

func ToNoteInfoFromDesc(noteInfo *desc.NoteInfo) *model.NoteInfo {
	return &model.NoteInfo{
		Title:   noteInfo.Title,
		Content: noteInfo.Content,
	}
}

func ToUpdateNoteInfoFromDesc(noteInfo *desc.UpdateNoteInfo) *model.UpdateNoteInfo {
	var updateNoteInfo = &model.UpdateNoteInfo{}

	if noteInfo.Title != nil {
		str := noteInfo.Title.GetValue()
		updateNoteInfo.Title = &str
	}
	if noteInfo.Content != nil {
		str := noteInfo.Content.GetValue()
		updateNoteInfo.Content = &str
	}

	return updateNoteInfo
}

func ToUpdateNoteFromDesc(id int64, noteInfo *desc.UpdateNoteInfo) *model.UpdateNote {
	return &model.UpdateNote{ID: id, Info: ToUpdateNoteInfoFromDesc(noteInfo)}
}
