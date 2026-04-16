package services

import "runtime"

type AppService struct{}

type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	GoVersion   string `json:"goVersion"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) Info() AppInfo {
	return AppInfo{
		Name:        "MyTools",
		Version:     "0.1.0",
		Description: "Desktop toolbox",
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}
}

func (s *AppService) Health() string {
	return "ok"
}
