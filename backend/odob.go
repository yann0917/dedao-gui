package backend

import (
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
)

// AudioDetail 听书简介
func (a *App) AudioDetail(id string) (detail *services.AudioInfoResp, err error) {
	detail, err = app.AudioDetail(id)
	if err != nil {
		return
	}

	return
}

// AudioDetailAlias 听书音频播放信息（包含 mp3_play_url）
func (a *App) AudioDetailAlias(aliasID string) (detail *services.Audio, err error) {
	detail, err = app.AudioDetailAlias(aliasID)
	if err != nil {
		return
	}
	return
}

func (a *App) OdobShelfAdd(enids []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = app.OdobShelfAdd(enids)
	if err != nil {
		return
	}
	return
}
