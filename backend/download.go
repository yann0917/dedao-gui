package backend

import (
	"github.com/yann0917/dedao/app"
)

func (a *App) CourseDownload(id, aid, dType int) (err error) {
	var d app.CourseDownload
	d.ID = id
	d.AID = aid
	d.DownloadType = dType
	err = d.Download()
	return
}

func (a *App) OdobDownload(id, dType int) (err error) {
	var d app.OdobDownload
	d.ID = id
	d.DownloadType = dType
	err = d.Download()
	return
}

func (a *App) EbookDownload(id, dType int) (err error) {
	var d app.EBookDownload
	d.ID = id
	d.DownloadType = dType
	err = d.Download()
	return
}
