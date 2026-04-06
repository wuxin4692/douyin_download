package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port           string
	DownloadDir    string
	FFmpegDir      string
	SiliconFlowKey string
	MiniMaxKey     string
	AIProvider     string
	AIBaseURL      string
	AIModel        string
}

var AppConfig *Config

func Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	downloadDir := filepath.Join(homeDir, "Downloads", "douyin-tool")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	// 默认 ffmpeg 目录，用户需手动放置
	ffmpegDir := filepath.Join(homeDir, ".douyin-tool", "ffmpeg")

	AppConfig = &Config{
		Port:         getEnv("PORT", "8080"),
		DownloadDir:  downloadDir,
		FFmpegDir:    ffmpegDir,
		AIProvider:   getEnv("AI_PROVIDER", "openai"),
		AIBaseURL:    getEnv("AI_BASE_URL", ""),
		AIModel:      getEnv("AI_MODEL", "gpt-3.5-turbo"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
