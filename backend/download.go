package backend

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/downloadmgr"
	"github.com/yann0917/dedao-gui/backend/services"
	"github.com/yann0917/dedao-gui/backend/utils"
)

func (a *App) OpenDirectoryDialog(title string) (dir string, err error) {
	home, _ := os.LookupEnv("HOME")
	dialogOptions := wailsruntime.OpenDialogOptions{
		DefaultDirectory:           home,
		Title:                      title,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	}
	dir, err = wailsruntime.OpenDirectoryDialog(a.Ctx, dialogOptions)
	app.SetOutputDir(dir)
	return
}

func (a *App) OpenFileDialog(title string) (file string, err error) {
	home, _ := os.LookupEnv("HOME")
	dialogOptions := wailsruntime.OpenDialogOptions{
		DefaultDirectory:           home,
		Title:                      title,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	}
	file, err = wailsruntime.OpenFileDialog(a.Ctx, dialogOptions)
	return
}

func (a *App) SetDir(dir []string) (err error) {
	if len(dir) > 0 {
		app.OutputDir = dir[0]
	}
	if len(dir) > 1 {
		utils.FfmpegDir = dir[1]
	}
	if len(dir) > 2 {
		utils.WkToPdfDir = dir[2]
	}
	if len(dir) > 1 {
		if err = validateExecutablePath(utils.FfmpegDir, "ffmpeg"); err != nil {
			return err
		}
	}
	if len(dir) > 2 {
		if err = validateExecutablePath(utils.WkToPdfDir, "wkhtmltopdf"); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) CourseDownload(id, aid, dType int, enid string) (err error) {
	title := "课程_" + strconv.Itoa(id)
	if enid != "" {
		if info, e := app.CourseInfoByEnid(enid); e == nil && info != nil && info.ClassInfo.Name != "" {
			title = info.ClassInfo.Name
		}
	}
	_, err = a.CreateDownloadTask(downloadmgr.CreateTaskRequest{
		BizType:      downloadmgr.BizTypeCourse,
		BizID:        id,
		Title:        title,
		EnID:         enid,
		ArticleID:    aid,
		DownloadType: dType,
		SaveDir:      app.OutputDir,
	})
	return
}

func (a *App) OdobDownload(id, dType int, data *services.Course) (err error) {
	if data == nil {
		return fmt.Errorf("缺少 odob 下载数据")
	}
	payload, err := jsoniter.MarshalToString(downloadmgr.OdobPayload{Data: *data})
	if err != nil {
		return err
	}
	_, err = a.CreateDownloadTask(downloadmgr.CreateTaskRequest{
		BizType:      downloadmgr.BizTypeOdob,
		BizID:        id,
		Title:        data.Title,
		EnID:         data.Enid,
		DownloadType: dType,
		SaveDir:      app.OutputDir,
		PayloadJSON:  payload,
	})
	return
}

func (a *App) EbookDownload(id, dType int, enid string) (err error) {
	title := "电子书_" + strconv.Itoa(id)
	if enid != "" {
		if detail, e := app.EbookDetail(enid); e == nil && detail != nil {
			if detail.Title != "" {
				title = detail.Title
			} else if detail.OperatingTitle != "" {
				title = detail.OperatingTitle
			}
		}
	}
	_, err = a.CreateDownloadTask(downloadmgr.CreateTaskRequest{
		BizType:      downloadmgr.BizTypeEbook,
		BizID:        id,
		Title:        title,
		EnID:         enid,
		DownloadType: dType,
		SaveDir:      app.OutputDir,
	})
	return
}

func validateExecutablePath(path string, label string) error {
	if path == "" {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s 路径无效", label)
	}
	if info.IsDir() {
		return fmt.Errorf("%s 必须是可执行文件", label)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("%s 没有执行权限", label)
	}
	return nil
}
