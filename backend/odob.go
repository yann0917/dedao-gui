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

func (a *App) OdobShelfAdd(enids []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = app.OdobShelfAdd(enids)
	if err != nil {
		return
	}
	return
}
