package app

import "github.com/yann0917/dedao-gui/backend/services"

// TopicAll 推荐话题列表
func TopicAll(page, limit int) (list *services.TopicAll, err error) {
	list, err = getService().TopicAll(page, limit)
	if err != nil {
		return
	}
	return
}

// TopicDetail Topic Detail
func TopicDetail(id string) (detail *services.TopicDetail, err error) {
	detail, err = getService().TopicDetail(id)
	if err != nil {
		return
	}
	return
}

// TopicNotesList Topic NotesList
func TopicNotesList(id string, page, limit int) (list *services.NotesList, err error) {
	list, err = getService().TopicNotesList(id, page, limit)
	if err != nil {
		return
	}
	return
}

// TopicNotesTimeline Topic timeline
func TopicNotesTimeline(maxID string) (list *services.NotesTimeline, err error) {
	list, err = getService().TopicNotesTimeline(maxID)
	if err != nil {
		return
	}
	return
}
