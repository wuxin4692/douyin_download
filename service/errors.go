package service

import "errors"

var (
	ErrInvalidLink    = errors.New("无效的抖音链接")
	ErrParseFailed    = errors.New("解析失败")
	ErrVideoNotFound  = errors.New("视频不存在")
	ErrDownloadFailed = errors.New("下载失败")
	ErrFFmpegNotFound = errors.New("ffmpeg 未安装")
	ErrAPIError       = errors.New("API 调用失败")
)
