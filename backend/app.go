package backend

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
	"github.com/yann0917/dedao-gui/backend/downloadmgr"
	"github.com/yann0917/dedao-gui/backend/utils"
)

// App struct
type App struct {
	Ctx             context.Context
	Notifier        *notifications.NotificationService
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
		Notifier:      notifications.New(),
		downloadState: downloadManagerStateUninitialized,
	}
}

// ServiceStartup is called when the service starts. The context is saved
// so we can call the runtime methods
func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.Ctx = ctx
	if err := a.ensureDownloadManager(); err != nil {
		fmt.Printf("启动下载任务管理器失败: %v\n", err)
	}
	return nil
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

func (a *App) ServiceShutdown() error {
	a.shutdownDownloadManager()
	// 退出前关闭缓存库：Close 内会回收一轮 value log，被删除的章节内容才能释放磁盘
	if err := utils.CloseBadgerDB(); err != nil {
		fmt.Printf("关闭缓存数据库时出错: %v\n", err)
	}
	return nil
}

func (a *App) OnSecondInstanceLaunch(secondInstanceData application.SecondInstanceData) {
	fmt.Println("OnSecondInstanceLaunch", secondInstanceData)
}
