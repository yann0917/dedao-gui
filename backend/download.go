package backend

import (
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/services"
	"github.com/yann0917/dedao-gui/backend/utils"
)

func (a *App) OpenDirectoryDialog(title string) (dir string, err error) {
	home, _ := os.LookupEnv("HOME")
	dialogOptions := runtime.OpenDialogOptions{
		DefaultDirectory:           home,
		Title:                      title,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	}
	dir, err = runtime.OpenDirectoryDialog(a.Ctx, dialogOptions)
	app.SetOutputDir(dir)
	return
}

func (a *App) SetDir(dir []string) {
	app.OutputDir = dir[0]
	utils.FfmpegDir = dir[1]
	utils.WkToPdfDir = dir[2]
}

func (a *App) CourseDownload(id, aid, dType int, enid string) (err error) {
	var d app.CourseDownload
	d.Ctx = a.Ctx
	d.ID = id
	d.AID = aid
	d.EnId = enid
	d.DownloadType = dType
	err = d.Download()
	return
}

func (a *App) OdobDownload(id, dType int, data *services.Course) (err error) {
	var d app.OdobDownload
	d.Ctx = a.Ctx
	d.ID = id
	d.DownloadType = dType
	d.Data = data
	err = d.Download()
	return
}

func (a *App) EbookDownload(id, dType int, enid string) (err error) {
	var d app.EBookDownload
	d.Ctx = a.Ctx
	d.ID = id
	d.DownloadType = dType
	d.EnID = enid
	err = d.Download()
	return
}
