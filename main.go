package main

import (
	"embed"

	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/yann0917/dedao-gui/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	_ "github.com/yann0917/dedao-gui/backend/config"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := backend.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "dedao-gui",
		Width:     1280,
		Height:    1000,
		MinWidth:  1024,
		MinHeight: 768,
		MaxWidth:  2560,
		MaxHeight: 1440,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:   &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:          app.Startup,
		OnShutdown:         app.Shutdown,
		OnDomReady:         app.DomReady,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.ERROR,
		WindowStartState:   options.Normal,
		Bind: []interface{}{
			app,
		},
		ErrorFormatter: func(err error) any { return err.Error() },
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "b7d0c23a-e0eb-4949-aa69-bb2f8ebe40e2",
			OnSecondInstanceLaunch: app.OnSecondInstanceLaunch,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			// WebviewIsTransparent:              true,
			// WindowIsTranslucent:               false,
			// DisableWindowIcon:                 false,
			// DisableFramelessWindowDecorations: false,
			// WebviewUserDataPath:               "",
			// WebviewBrowserPath:                "",
			Theme: windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "dedao gui downloader",
				Message: "https://github.com/yann0917/dedao-gui",
				Icon:    icon,
			},
		},
		Linux: &linux.Options{
			Icon: icon,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
