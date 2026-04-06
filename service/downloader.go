package service

import (
	"fmt"
	"path/filepath"
	"strings"

	"douyin-tool/config"
	"douyin-tool/utils"
)

type DownloadResult struct {
	VideoPath string `json:"video_path"`
	Title     string `json:"title"`
	Size      int64  `json:"size"`
}

// DownloadVideo downloads video to local
func DownloadVideo(videoInfo *VideoInfo, taskID string) (*DownloadResult, error) {
	UpdateTask(taskID, "processing", 10, "开始下载视频...")

	outputDir := config.AppConfig.DownloadDir
	filename := sanitizeFilename(videoInfo.Title + ".mp4")
	if filename == ".mp4" {
		filename = videoInfo.VideoID + ".mp4"
	}

	videoPath, err := utils.DownloadFile(videoInfo.DownloadURL, outputDir, true,
		func(downloaded, total int64) {
			progress := int(float64(downloaded) / float64(total) * 90)
			UpdateTask(taskID, "processing", 10+progress, fmt.Sprintf("下载中... %d%%", 10+progress))
		})

	if err != nil {
		SetTaskError(taskID, err.Error())
		return nil, fmt.Errorf("下载失败: %w", err)
	}

	// 重命名文件
	finalPath := filepath.Join(outputDir, filename)
	if videoPath != finalPath {
		// 文件已下载到 outputDir，需要重命名
	}

	size := utils.GetFileSize(videoPath)

	UpdateTask(taskID, "processing", 100, "下载完成")

	return &DownloadResult{
		VideoPath: videoPath,
		Title:     videoInfo.Title,
		Size:      size,
	}, nil
}

func sanitizeFilename(name string) string {
	if name == "" {
		return ""
	}
	
	// 移除或替换非法字符
	invalid := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// 移除首尾空格和点
	result = strings.Trim(result, " .")
	
	if result == "" {
		return ""
	}
	
	return result + ".mp4"
}
