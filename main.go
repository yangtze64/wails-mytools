package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"wails-mytools/internal/services"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	devToolsService := services.NewDevToolsService()
	app := application.New(application.Options{
		Name:        "wails-mytools",
		Description: "A desktop toolbox built with Wails and Vue",
		Services: []application.Service{
			application.NewService(services.NewAppService()),
			application.NewService(services.NewBookmarkService()),
			application.NewService(services.NewSecretManagerService()),
			application.NewService(devToolsService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	devToolsService.SetApp(app)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "MyTools",
		Width:     1240,
		Height:    820,
		MinWidth:  1100,
		MinHeight: 720,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(15, 23, 42),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
