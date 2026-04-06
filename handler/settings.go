package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"douyin-tool/config"

	"github.com/gin-gonic/gin"
)

type Settings struct {
	SiliconFlowKey string `json:"silicon_flow_key"`
	MiniMaxKey     string `json:"minimax_key"`
	DownloadDir    string `json:"download_dir"`
}

func GetSettings(c *gin.Context) {
	settings := Settings{
		SiliconFlowKey: config.AppConfig.SiliconFlowKey,
		MiniMaxKey:     config.AppConfig.MiniMaxKey,
		DownloadDir:    config.AppConfig.DownloadDir,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": settings,
	})
}

func SaveSettings(c *gin.Context) {
	var req Settings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的请求参数",
		})
		return
	}

	// 更新配置
	config.AppConfig.SiliconFlowKey = req.SiliconFlowKey
	config.AppConfig.MiniMaxKey = req.MiniMaxKey

	if req.DownloadDir != "" {
		if err := os.MkdirAll(req.DownloadDir, 0755); err == nil {
			config.AppConfig.DownloadDir = req.DownloadDir
		}
	}

	// 保存到配置文件
	settingsPath := filepath.Join(getConfigDir(), "settings.json")
	data, _ := json.Marshal(map[string]interface{}{
		"silicon_flow_key": req.SiliconFlowKey,
		"minimax_key":     req.MiniMaxKey,
		"download_dir":    config.AppConfig.DownloadDir,
	})
	os.WriteFile(settingsPath, data, 0644)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置已保存",
	})
}

func getConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".douyin-tool")
	os.MkdirAll(configDir, 0755)
	return configDir
}
