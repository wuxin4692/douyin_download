package main

import (
	"fmt"
	"log"
	"net/http"

	"douyin-tool/config"
	"douyin-tool/handler"
	"douyin-tool/service"
	"douyin-tool/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化 ffmpeg
	log.Println("正在检查 ffmpeg...")
	if err := utils.EnsureFFmpeg(); err != nil {
		log.Printf("ffmpeg 初始化警告: %v", err)
	} else {
		log.Println("ffmpeg 就绪")
	}

	// 初始化服务
	service.Init()

	// 设置 Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 静态文件
	r.Static("/static", "./static")

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/ffmpeg/status", handler.FFmpegStatus)
		api.GET("/settings", handler.GetSettings)
		api.POST("/settings", handler.SaveSettings)
		api.POST("/video/parse", handler.ParseVideo)
		api.POST("/video/download", handler.DownloadVideo)
		api.POST("/audio/extract", handler.ExtractAudio)
		api.POST("/task/transcribe", handler.Transcribe)
		api.POST("/task/segment", handler.Segment)
		api.GET("/task/status/:id", handler.GetTaskStatus)
		api.POST("/file/save", handler.SaveFile)
	}

	fmt.Printf(`
╔══════════════════════════════════════════════════════════╗
║           抖音工具箱 - Douyin Tool                    ║
╠══════════════════════════════════════════════════════════╣
║  服务地址: http://localhost:%s                          ║
║  下载目录: %s
╚══════════════════════════════════════════════════════════╝
`, config.AppConfig.Port, config.AppConfig.DownloadDir)

	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
