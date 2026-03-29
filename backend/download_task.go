package backend

import (
	"errors"

	"github.com/yann0917/dedao-gui/backend/app"
	"github.com/yann0917/dedao-gui/backend/downloadmgr"
)

func (a *App) CreateDownloadTask(req downloadmgr.CreateTaskRequest) (task *downloadmgr.DownloadTask, err error) {
	if a.DownloadRepo == nil {
		return nil, errors.New("下载任务管理器未初始化")
	}
	if req.BizType == "" {
		return nil, errors.New("bizType 不能为空")
	}
	if req.BizID <= 0 {
		return nil, errors.New("bizId 必须大于 0")
	}
	if req.DownloadType <= 0 {
		return nil, errors.New("downloadType 必须大于 0")
	}
	if req.SaveDir == "" {
		req.SaveDir = app.OutputDir
	}
	task, err = a.DownloadRepo.CreateTask(req)
	return
}

func (a *App) ListDownloadTasks(query downloadmgr.ListTaskQuery) (result downloadmgr.ListTaskResult, err error) {
	if a.DownloadRepo == nil {
		err = errors.New("下载任务管理器未初始化")
		return
	}
	result, err = a.DownloadRepo.ListTasks(query)
	return
}

func (a *App) GetDownloadTask(id string) (task *downloadmgr.DownloadTask, err error) {
	if a.DownloadRepo == nil {
		return nil, errors.New("下载任务管理器未初始化")
	}
	task, err = a.DownloadRepo.GetTask(id)
	return
}

func (a *App) CancelDownloadTask(id string) (err error) {
	if a.DownloadManager == nil {
		return errors.New("下载任务管理器未初始化")
	}
	return a.DownloadManager.CancelTask(id)
}

func (a *App) RetryDownloadTask(id string) (err error) {
	if a.DownloadRepo == nil {
		return errors.New("下载任务管理器未初始化")
	}
	return a.DownloadRepo.RetryTask(id)
}

func (a *App) PauseDownloadTask(id string) (err error) {
	if a.DownloadManager == nil {
		return errors.New("下载任务管理器未初始化")
	}
	return a.DownloadManager.PauseTask(id)
}

func (a *App) ResumeDownloadTask(id string) (err error) {
	if a.DownloadManager == nil {
		return errors.New("下载任务管理器未初始化")
	}
	return a.DownloadManager.ResumeTask(id)
}

func (a *App) ClearDownloadTasks(clearAll bool) (err error) {
	if a.DownloadManager == nil {
		return errors.New("下载任务管理器未初始化")
	}
	return a.DownloadManager.ClearTasks(clearAll)
}
