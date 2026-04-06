package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"douyin-tool/config"
)

func SaveToFile(content string, filename string) (string, error) {
	// 确保文件名有 .md 后缀
	if !strings.HasSuffix(strings.ToLower(filename), ".md") {
		filename += ".md"
	}

	// 使用配置的下载目录
	dir := config.AppConfig.DownloadDir
	if dir == "" {
		dir = "downloads"
	}

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 生成完整路径
	filePath := filepath.Join(dir, filename)

	// 如果文件已存在，添加时间戳
	if _, err := os.Stat(filePath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		name := strings.TrimSuffix(filename, ".md")
		filename = fmt.Sprintf("%s_%s.md", name, timestamp)
		filePath = filepath.Join(dir, filename)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return filePath, nil
}
