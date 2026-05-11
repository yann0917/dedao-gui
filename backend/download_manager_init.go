package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/yann0917/dedao-gui/backend/downloadmgr"
)

type downloadManagerState string

const (
	downloadManagerStateUninitialized downloadManagerState = "uninitialized"
	downloadManagerStateInitializing  downloadManagerState = "initializing"
	downloadManagerStateReady         downloadManagerState = "ready"
	downloadManagerStateFailed        downloadManagerState = "failed"
)

type downloadManagerInitializer func(context.Context) (*downloadmgr.Repository, *downloadmgr.Manager, error)

func (a *App) ensureDownloadManager() error {
	ctx := a.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	waitedForInit := false
	for {
		a.downloadInitMu.Lock()
		switch a.downloadState {
		case downloadManagerStateReady:
			a.downloadInitMu.Unlock()
			return nil
		case downloadManagerStateInitializing:
			waitCh := a.downloadInitCh
			a.downloadInitMu.Unlock()
			if waitCh != nil {
				<-waitCh
			}
			waitedForInit = true
			continue
		case downloadManagerStateFailed:
			if waitedForInit {
				err := a.downloadInitErr
				a.downloadInitMu.Unlock()
				if err == nil {
					return fmt.Errorf("下载任务管理器初始化失败")
				}
				return err
			}
			fallthrough
		default:
			waitCh := make(chan struct{})
			a.downloadState = downloadManagerStateInitializing
			a.downloadInitCh = waitCh
			a.downloadInitMu.Unlock()

			repo, manager, err := a.initDownloadManager(ctx)

			a.downloadInitMu.Lock()
			if err != nil {
				a.DownloadRepo = nil
				a.DownloadManager = nil
				a.downloadInitErr = err
				a.downloadState = downloadManagerStateFailed
			} else {
				a.DownloadRepo = repo
				a.DownloadManager = manager
				a.downloadInitErr = nil
				a.downloadState = downloadManagerStateReady
			}
			a.downloadInitCh = nil
			close(waitCh)
			a.downloadInitMu.Unlock()

			if err != nil {
				return err
			}
			return nil
		}
	}
}

func (a *App) initDownloadManager(ctx context.Context) (*downloadmgr.Repository, *downloadmgr.Manager, error) {
	if a.downloadInitFn != nil {
		return a.downloadInitFn(ctx)
	}

	repo, err := downloadmgr.NewRepository()
	if err != nil {
		return nil, nil, fmt.Errorf("初始化下载任务数据库失败: %w", err)
	}

	manager := downloadmgr.NewManager(ctx, repo, downloadmgr.Config{
		WorkerCount:  readEnvInt("DOWNLOAD_TASK_WORKER_COUNT", 1),
		PollInterval: time.Duration(readEnvInt("DOWNLOAD_TASK_POLL_MS", 1000)) * time.Millisecond,
	})
	manager.RegisterExecutor(downloadmgr.BizTypeCourse, downloadmgr.NewCourseExecutor(ctx))
	manager.RegisterExecutor(downloadmgr.BizTypeOdob, downloadmgr.NewOdobExecutor(ctx))
	manager.RegisterExecutor(downloadmgr.BizTypeEbook, downloadmgr.NewEbookExecutor(ctx))
	if err = manager.Start(); err != nil {
		closeDownloadRepo(repo)
		return nil, nil, fmt.Errorf("启动下载任务管理器失败: %w", err)
	}

	return repo, manager, nil
}

func (a *App) shutdownDownloadManager() {
	a.downloadInitMu.Lock()
	manager := a.DownloadManager
	repo := a.DownloadRepo
	a.DownloadManager = nil
	a.DownloadRepo = nil
	a.downloadState = downloadManagerStateUninitialized
	a.downloadInitErr = nil
	a.downloadInitCh = nil
	a.downloadInitMu.Unlock()

	if manager != nil {
		manager.Stop()
	}
	closeDownloadRepo(repo)
}

func closeDownloadRepo(repo *downloadmgr.Repository) {
	if repo == nil || repo.DB() == nil {
		return
	}
	sqlDB, err := repo.DB().DB()
	if err != nil || sqlDB == nil {
		return
	}
	_ = sqlDB.Close()
}
