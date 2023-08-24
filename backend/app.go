package backend

import (
	"context"
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
}

func (a *App) DomReady(ctx context.Context) {
	// fmt.Println(a.Ctx)
	// fmt.Println("dom ready")
}
