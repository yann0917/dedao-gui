package backend

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/yann0917/dedao-gui/backend/utils"
)

// App struct
type App struct {
	Ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	// fmt.Println(a.Ctx)
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
