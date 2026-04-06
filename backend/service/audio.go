package service

import (
	"fmt"
	"os"
	"path/filepath"

	"douyin-tool/config"
	"douyin-tool/utils"
)

type AudioResult struct {
	AudioPath string `json:"audio_path"`
	Duration  float64 `json:"duration"`
}

// ExtractAudio extracts audio from video file
func ExtractAudio(videoPath string, taskID string) (*AudioResult, error) {
	UpdateTask(taskID, "processing", 0, "开始提取音频...")

	// 确保 ffmpeg 可用
	_, err := os.Stat(utils.GetFFmpegPath())
	if err != nil {
		return nil, ErrFFmpegNotFound
	}

	// 生成输出路径
	audioDir := filepath.Join(config.AppConfig.DownloadDir, "audio")
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return nil, err
	}

	videoName := filepath.Base(videoPath)
	audioName := videoName[:len(videoName)-len(filepath.Ext(videoName))] + ".mp3"
	audioPath := filepath.Join(audioDir, audioName)

	// 使用 ffmpeg 提取音频
	UpdateTask(taskID, "processing", 50, "提取音频中...")

	_, err = utils.RunFFmpeg([]string{
		"-i", videoPath,
		"-vn",
		"-acodec", "libmp3lame",
		"-q:a", "2",
		"-y",
		audioPath,
	})

	if err != nil {
		SetTaskError(taskID, err.Error())
		return nil, fmt.Errorf("音频提取失败: %w", err)
	}

	// 获取音频时长
	duration, err := utils.GetVideoDuration(audioPath)
	if err != nil {
		duration = 0
	}

	UpdateTask(taskID, "completed", 100, "音频提取完成")

	return &AudioResult{
		AudioPath: audioPath,
		Duration:  duration,
	}, nil
}
