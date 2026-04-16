package services

import (
	"os"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type DevToolsService struct {
	app *application.App
}

type DevEnvironment struct {
	GoVersion        string `json:"goVersion"`
	OS               string `json:"os"`
	Arch             string `json:"arch"`
	FrontendDevURL   string `json:"frontendDevUrl"`
	WorkingDirectory string `json:"workingDirectory"`
}

func NewDevToolsService() *DevToolsService {
	return &DevToolsService{}
}

func (s *DevToolsService) SetApp(app *application.App) {
	s.app = app
}

func (s *DevToolsService) OpenDevTools() {
	if s.app == nil {
		return
	}

	window := s.app.Window.Current()
	if window == nil {
		return
	}

	window.OpenDevTools()
}

func (s *DevToolsService) Environment() DevEnvironment {
	workingDirectory, _ := os.Getwd()

	return DevEnvironment{
		GoVersion:        runtime.Version(),
		OS:               runtime.GOOS,
		Arch:             runtime.GOARCH,
		FrontendDevURL:   os.Getenv("FRONTEND_DEVSERVER_URL"),
		WorkingDirectory: workingDirectory,
	}
}
