package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"douyin-tool/config"

	"github.com/google/uuid"
)

type ProgressCallback func(downloaded, total int64)

func DownloadFile(url string, outputDir string, showProgress bool, callback ProgressCallback) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 从 Content-Disposition 获取文件名
	filename := getFilenameFromHeaders(resp.Header, url)
	if filename == "" {
		filename = uuid.New().String() + ".mp4"
	}

	// 清理文件名
	filename = sanitizeFilename(filename)

	filepath := filepath.Join(outputDir, filename)

	total := resp.ContentLength
	if total <= 0 {
		total = 100 * 1024 * 1024 // 默认 100MB
	}

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err := file.Write(buf[:n])
			if err != nil {
				return "", err
			}
			downloaded += int64(n)
			if callback != nil {
				callback(downloaded, total)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return filepath, nil
}

func getFilenameFromHeaders(header http.Header, url string) string {
	cd := header.Get("Content-Disposition")
	if cd != "" {
		parts := strings.Split(cd, "filename=")
		if len(parts) > 1 {
			filename := strings.Trim(parts[1], "\" ")
			if filename != "" {
				return filename
			}
		}
	}

	// 从 URL 提取
	paths := strings.Split(url, "/")
	if len(paths) > 0 {
		last := paths[len(paths)-1]
		if idx := strings.Index(last, "?"); idx > 0 {
			last = last[:idx]
		}
		if strings.Contains(last, ".") {
			return last
		}
	}

	return ""
}

func sanitizeFilename(filename string) string {
	// 移除非法字符
	invalid := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := filename
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func GetVideoDuration(videoPath string) (float64, error) {
	output, err := RunFFprobe([]string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	})
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(output), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func GetFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func GetDownloadDir() string {
	return config.AppConfig.DownloadDir
}
