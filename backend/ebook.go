package backend

import (
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
)

func (a *App) EbookInfo(enid string) (info *services.EbookDetail, err error) {
	return app.EbookDetail(enid)
}

func (a *App) EbookCommentList(enid string, page, limit int) (list *services.EbookCommentList, err error) {
	list, err = app.EbookCommentList(enid, "like_count", page, limit)
	return
}

func (a *App) EbookShelfAdd(enids []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = app.EbookShelfAdd(enids)
	return
}

func (a *App) EbookShelfRemove(enids []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = app.EbookShelfRemove(enids)
	return
}
