package handler

import (
	"encoding/json"
	"os"
	"path/filepath"

	"douyin-tool/config"
	"douyin-tool/service"
)

// ParseVideoSync 同步解析视频
func ParseVideoSync(url string, taskID string) (interface{}, error) {
	videoInfo, err := service.ParseShareUrl(url)
	if err != nil {
		return nil, err
	}
	service.UpdateTask(taskID, "completed", 100, "解析完成")
	return videoInfo, nil
}

// DownloadVideoSync 同步下载视频
func DownloadVideoSync(url string, taskID string) (interface{}, error) {
	// 解析链接
	videoInfo, err := service.ParseShareUrl(url)
	if err != nil {
		return nil, err
	}

	// 下载视频
	result, err := service.DownloadVideo(videoInfo, taskID)
	if err != nil {
		return nil, err
	}

	service.UpdateTask(taskID, "completed", 100, "下载完成")
	return result, nil
}

// ExtractAudioSync 同步提取音频
func ExtractAudioSync(videoPath string, taskID string) (interface{}, error) {
	result, err := service.ExtractAudio(videoPath, taskID)
	if err != nil {
		return nil, err
	}
	service.UpdateTask(taskID, "completed", 100, "音频提取完成")
	return result, nil
}

// SaveSettingsToFile 保存设置到文件
func SaveSettingsToFile(cfg *config.Config) error {
	settingsPath := filepath.Join(getConfigDir(), "settings.json")
	data, _ := json.Marshal(map[string]interface{}{
		"silicon_flow_key": cfg.SiliconFlowKey,
		"minimax_key":     cfg.MiniMaxKey,
		"download_dir":    cfg.DownloadDir,
	})
	return os.WriteFile(settingsPath, data, 0644)
}
