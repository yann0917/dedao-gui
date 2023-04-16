package backend

import (
	"github.com/yann0917/dedao/app"
	"github.com/yann0917/dedao/services"
)

func (a *App) EbookInfo(enid string) (info *services.EbookDetail, err error) {
	return app.EbookDetailByEnid(enid)
}

func (a *App) EbookCommentList(enid string, page, limit int) (list *services.EbookCommentList, err error) {
	list, err = app.EbookCommentList(enid, "like_count", page, limit)
	return
}

func (a *App) EbookShelfAdd(enids []string) (resp *services.EbookShelfAddResp, err error) {
	resp, err = app.EbookShelfAdd(enids)
	return
}
