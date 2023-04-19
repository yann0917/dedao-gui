package backend

import (
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao/app"
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

func (a *App) SetDownloadDir(dir string) {
	app.SetOutputDir(dir)
}

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
