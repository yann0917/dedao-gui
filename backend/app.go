package backend

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/options"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yann0917/dedao-gui/backend/downloadmgr"
	"github.com/yann0917/dedao-gui/backend/utils"
)

// App struct
type App struct {
	Ctx             context.Context
	DownloadRepo    *downloadmgr.Repository
	DownloadManager *downloadmgr.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
	repo, err := downloadmgr.NewRepository()
	if err != nil {
		fmt.Printf("初始化下载任务数据库失败: %v\n", err)
		return
	}
	manager := downloadmgr.NewManager(ctx, repo, downloadmgr.Config{
		WorkerCount:  readEnvInt("DOWNLOAD_TASK_WORKER_COUNT", 1),
		PollInterval: time.Duration(readEnvInt("DOWNLOAD_TASK_POLL_MS", 1000)) * time.Millisecond,
	})
	manager.RegisterExecutor(downloadmgr.BizTypeCourse, downloadmgr.NewCourseExecutor(ctx))
	manager.RegisterExecutor(downloadmgr.BizTypeOdob, downloadmgr.NewOdobExecutor(ctx))
	manager.RegisterExecutor(downloadmgr.BizTypeEbook, downloadmgr.NewEbookExecutor(ctx))
	if err = manager.Start(); err != nil {
		fmt.Printf("启动下载任务管理器失败: %v\n", err)
		return
	}
	a.DownloadRepo = repo
	a.DownloadManager = manager

	if err := wailsruntime.InitializeNotifications(ctx); err != nil {
		fmt.Printf("初始化通知系统失败: %v\n", err)
	} else if wailsruntime.IsNotificationAvailable(ctx) {
		authorized, err := wailsruntime.CheckNotificationAuthorization(ctx)
		if err != nil {
			fmt.Printf("检查通知权限失败: %v\n", err)
		} else if !authorized {
			authorized, err = wailsruntime.RequestNotificationAuthorization(ctx)
			if err != nil {
				fmt.Printf("请求通知权限失败: %v\n", err)
			} else if authorized {
				fmt.Println("通知权限已授权")
			}
		}
		registerNotificationCategory(ctx)
		wailsruntime.OnNotificationResponse(ctx, a.onNotificationResponse)
	}
}

func readEnvInt(name string, defaultValue int) int {
	raw := os.Getenv(name)
	if raw == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *App) Shutdown(ctx context.Context) {
	if a.DownloadManager != nil {
		a.DownloadManager.Stop()
	}
	setupCleanupOnExit()
}

func (a *App) DomReady(ctx context.Context) {
	// fmt.Println(a.Ctx)
	// fmt.Println("dom ready")
}

func (a *App) OnSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	fmt.Println("OnSecondInstanceLaunch", secondInstanceData)
}

func setupCleanupOnExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("正在关闭程序...")

		// 获取 BadgerDB 实例并关闭
		db, err := utils.GetBadgerDB(utils.GetDefaultBadgerDBPath())
		if err == nil && db != nil {
			if err := db.Close(); err != nil {
				fmt.Printf("关闭数据库时出错: %v\n", err)
			} else {
				fmt.Println("数据库已安全关闭")
			}
		}

		os.Exit(0)
	}()
}

func registerNotificationCategory(ctx context.Context) {
	category := wailsruntime.NotificationCategory{
		ID: "download-category",
		Actions: []wailsruntime.NotificationAction{
			{ID: "open_folder", Title: "打开文件夹"},
		},
	}
	if err := wailsruntime.RegisterNotificationCategory(ctx, category); err != nil {
		fmt.Printf("注册通知类别失败: %v\n", err)
	} else {
		fmt.Println("通知类别已注册")
	}
}

func (a *App) onNotificationResponse(result wailsruntime.NotificationResult) {
	if result.Error != nil {
		fmt.Printf("通知回调错误: %v\n", result.Error)
		return
	}
	response := result.Response
	if response.ActionIdentifier == "open_folder" {
		if saveDir, ok := response.UserInfo["saveDir"].(string); ok && saveDir != "" {
			openFolder(saveDir)
		}
	}
}

func openFolder(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	if err := cmd.Run(); err != nil {
		fmt.Printf("打开文件夹失败: %v\n", err)
	}
}
