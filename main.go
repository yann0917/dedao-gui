package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/yann0917/dedao-gui/backend"
	_ "github.com/yann0917/dedao-gui/backend/config"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	app := backend.NewApp()

	wailsApp := application.New(application.Options{
		Name: "dedao-gui",
		Icon: icon,
		Services: []application.Service{
			application.NewService(app),
			application.NewService(app.Notifier),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID:               "b7d0c23a-e0eb-4949-aa69-bb2f8ebe40e2",
			OnSecondInstanceLaunch: app.OnSecondInstanceLaunch,
		},
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "dedao-gui",
		Width:          1280,
		Height:         1000,
		MinWidth:       1024,
		MinHeight:      768,
		MaxWidth:       2560,
		MaxHeight:      1440,
		Frameless:      true,
		BackgroundType: application.BackgroundTypeTransparent,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTransparent,
		},
	})

	if err := wailsApp.Run(); err != nil {
		println("Error:", err.Error())
	}
}
