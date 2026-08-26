package backend

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/yann0917/dedao-gui/backend/downloadmgr"
	"github.com/yann0917/dedao-gui/backend/utils"
)

// App struct
type App struct {
	Ctx             context.Context
	DownloadRepo    *downloadmgr.Repository
	DownloadManager *downloadmgr.Manager
	downloadInitMu  sync.Mutex
	downloadInitCh  chan struct{}
	downloadInitErr error
	downloadState   downloadManagerState
	downloadInitFn  downloadManagerInitializer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		downloadState: downloadManagerStateUninitialized,
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
	if err := a.ensureDownloadManager(); err != nil {
		fmt.Printf("启动下载任务管理器失败: %v\n", err)
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
	a.shutdownDownloadManager()
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
