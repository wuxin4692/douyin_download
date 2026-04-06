package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	_ "embed"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建应用实例
	app := NewApp()

	// 创建应用选项
	goApp := &options.App{
		Title:     "抖音工具箱",
		Width:     1200,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 25, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	}

	// 运行应用
	err := wails.Run(goApp)
	if err != nil {
		log.Fatal(err)
	}
}
