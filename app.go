package main

import (
	"context"
	"fmt"
	"log"

	"douyin-tool/config"
	"douyin-tool/handler"
	"douyin-tool/service"
	"douyin-tool/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 结构体 - Wails 应用主类
type App struct {
	ctx context.Context
}

// NewApp 创建一个新的 App 实例
func NewApp() *App {
	return &App{}
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化配置
	if err := config.Init(); err != nil {
		log.Printf("配置初始化失败: %v", err)
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

	log.Println("应用启动完成")
}

// ============================================
// FFmpeg 相关
// ============================================

// GetFFmpegStatus 获取 FFmpeg 状态
func (a *App) GetFFmpegStatus() map[string]interface{} {
	return utils.GetFFmpegStatus()
}

// ============================================
// 设置相关
// ============================================

// GetSettings 获取设置
func (a *App) GetSettings() map[string]interface{} {
	return map[string]interface{}{
		"silicon_flow_key": config.AppConfig.SiliconFlowKey,
		"minimax_key":      config.AppConfig.MiniMaxKey,
		"download_dir":     config.AppConfig.DownloadDir,
		"ffmpeg_dir":       config.AppConfig.FFmpegDir,
		"ai_provider":      config.AppConfig.AIProvider,
		"ai_base_url":     config.AppConfig.AIBaseURL,
		"ai_model":        config.AppConfig.AIModel,
	}
}

// SaveSettings 保存设置
func (a *App) SaveSettings(settings map[string]interface{}) error {
	if key, ok := settings["silicon_flow_key"].(string); ok {
		config.AppConfig.SiliconFlowKey = key
	}
	if key, ok := settings["minimax_key"].(string); ok {
		config.AppConfig.MiniMaxKey = key
	}
	if dir, ok := settings["ffmpeg_dir"].(string); ok {
		config.AppConfig.FFmpegDir = dir
	}
	if provider, ok := settings["ai_provider"].(string); ok {
		config.AppConfig.AIProvider = provider
	}
	if baseURL, ok := settings["ai_base_url"].(string); ok {
		config.AppConfig.AIBaseURL = baseURL
	}
	if model, ok := settings["ai_model"].(string); ok {
		config.AppConfig.AIModel = model
	}
	return handler.SaveSettingsToFile(config.AppConfig)
}

// ============================================
// 视频处理相关
// ============================================

// ParseVideo 解析视频
func (a *App) ParseVideo(url string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := handler.ParseVideoSync(url, task.ID)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// DownloadVideo 下载视频
func (a *App) DownloadVideo(url string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := handler.DownloadVideoSync(url, task.ID)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// ExtractAudio 提取音频
func (a *App) ExtractAudio(videoPath string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := handler.ExtractAudioSync(videoPath, task.ID)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// ============================================
// 任务相关
// ============================================

// GetTaskStatus 获取任务状态
func (a *App) GetTaskStatus(taskID string) map[string]interface{} {
	task := service.GetTask(taskID)
	if task == nil {
		return map[string]interface{}{
			"error": "任务不存在",
		}
	}

	result := map[string]interface{}{
		"id":       task.ID,
		"status":   task.Status,
		"progress": task.Progress,
		"message":  task.Message,
	}

	if task.Error != "" {
		result["error"] = task.Error
	}
	if task.Result != nil {
		result["result"] = task.Result
	}

	return result
}

// ============================================
// 语音转写与分段
// ============================================

// TranscribeAudio 语音转写
func (a *App) TranscribeAudio(audioPath string, apiKey string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := service.TranscribeAudio(audioPath, task.ID, apiKey)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// SemanticSegment 语义分段
func (a *App) SemanticSegment(text string, apiKey string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := service.SemanticSegment(text, task.ID, apiKey)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// ============================================
// LLM 分析与改写
// ============================================

// AnalyzeText LLM 分析文案
func (a *App) AnalyzeText(text string, apiKey string, provider string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := service.AnalyzeText(text, task.ID, apiKey, provider)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// RewriteText 一键改写文案
func (a *App) RewriteText(text string, apiKey string, provider string, style string, customInstruction string) (map[string]interface{}, error) {
	task := service.CreateTask()

	go func() {
		result, err := service.RewriteText(text, task.ID, apiKey, provider, style, customInstruction)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	return map[string]interface{}{
		"task_id": task.ID,
	}, nil
}

// ============================================
// 文件操作
// ============================================

// SaveFile 保存文案到文件
func (a *App) SaveFile(content string, filename string) (map[string]interface{}, error) {
	filepath, err := service.SaveToFile(content, filename)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"path": filepath,
	}, nil
}

// OpenFileDialog 打开文件选择对话框
func (a *App) OpenFileDialog(title string, filter string) string {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: filter,
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return ""
	}
	return result
}

// OpenDirectoryDialog 打开目录选择对话框
func (a *App) OpenDirectoryDialog(title string) string {
	result, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return ""
	}
	return result
}

// ============================================
// 系统相关
// ============================================

// GetAppVersion 获取应用版本
func (a *App) GetAppVersion() string {
	return "1.0.0"
}

// Log 日志输出到控制台
func (a *App) Log(message string) {
	fmt.Println("[Go]", message)
}
