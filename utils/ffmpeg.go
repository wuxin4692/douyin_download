package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"douyin-tool/config"
)

// EnsureFFmpeg 检查 ffmpeg 是否可用，不存在则返回错误
func EnsureFFmpeg() error {
	ffmpegPath := GetFFmpegPath()
	ffprobePath := GetFFprobePath()

	// 检查是否存在
	if _, err := os.Stat(ffmpegPath); err == nil {
		if _, err := os.Stat(ffprobePath); err == nil {
			return nil
		}
	}

	// 不再自动下载，返回错误提示用户手动配置
	return fmt.Errorf("ffmpeg 未找到，请手动下载并放置到: %s\n或将其添加到系统 PATH", config.AppConfig.FFmpegDir)
}

// GetFFmpegPath 获取 ffmpeg 路径
func GetFFmpegPath() string {
	ext := ".exe"
	if runtime.GOOS != "windows" {
		ext = ""
	}
	name := "ffmpeg" + ext

	// 先检查根目录
	path := filepath.Join(config.AppConfig.FFmpegDir, name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// Windows 下搜索子目录
	if runtime.GOOS == "windows" {
		entries, err := os.ReadDir(config.AppConfig.FFmpegDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					subPath := filepath.Join(config.AppConfig.FFmpegDir, entry.Name(), name)
					if _, err := os.Stat(subPath); err == nil {
						return subPath
					}
					// 搜索 bin 子目录
					binPath := filepath.Join(config.AppConfig.FFmpegDir, entry.Name(), "bin", name)
					if _, err := os.Stat(binPath); err == nil {
						return binPath
					}
				}
			}
		}
	}

	// 检查系统 PATH 中是否存在
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path
	}

	return path
}

// GetFFprobePath 获取 ffprobe 路径
func GetFFprobePath() string {
	ext := ".exe"
	if runtime.GOOS != "windows" {
		ext = ""
	}
	name := "ffprobe" + ext

	// 先检查根目录
	path := filepath.Join(config.AppConfig.FFmpegDir, name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// Windows 下搜索子目录
	if runtime.GOOS == "windows" {
		entries, err := os.ReadDir(config.AppConfig.FFmpegDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					subPath := filepath.Join(config.AppConfig.FFmpegDir, entry.Name(), name)
					if _, err := os.Stat(subPath); err == nil {
						return subPath
					}
					// 搜索 bin 子目录
					binPath := filepath.Join(config.AppConfig.FFmpegDir, entry.Name(), "bin", name)
					if _, err := os.Stat(binPath); err == nil {
						return binPath
					}
				}
			}
		}
	}

	// 检查系统 PATH 中是否存在
	if path, err := exec.LookPath("ffprobe"); err == nil {
		return path
	}

	return path
}

// GetFFmpegStatus 获取 ffmpeg 状态
func GetFFmpegStatus() map[string]interface{} {
	ffmpegPath := GetFFmpegPath()
	_, ffmpegErr := os.Stat(ffmpegPath)
	_, ffprobeErr := os.Stat(GetFFprobePath())

	return map[string]interface{}{
		"installed": ffmpegErr == nil && ffprobeErr == nil,
		"path":      config.AppConfig.FFmpegDir,
	}
}

// RunFFmpeg 执行 ffmpeg 命令
func RunFFmpeg(args []string) (string, error) {
	cmd := GetFFmpegPath()
	return runCommand(cmd, args...)
}

// RunFFprobe 执行 ffprobe 命令
func RunFFprobe(args []string) (string, error) {
	cmd := GetFFprobePath()
	return runCommand(cmd, args...)
}

func runCommand(name string, args ...string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%s: %s", err.Error(), stderr.String())
		}
		return "", err
	}

	return stdout.String(), nil
}
