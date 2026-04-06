package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"douyin-tool/config"
	"douyin-tool/utils"
)

type TranscriptionResult struct {
	Text      string  `json:"text"`
	Language  string  `json:"language"`
	Duration  float64 `json:"duration"`
}

// TranscribeAudio transcribes audio to text using Silicon Flow API
func TranscribeAudio(audioPath string, taskID string, apiKey string) (*TranscriptionResult, error) {
	UpdateTask(taskID, "processing", 0, "开始语音转文字...")

	if apiKey == "" {
		apiKey = config.AppConfig.SiliconFlowKey
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未设置 Silicon Flow API Key")
	}

	// 构建 multipart form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("打开音频文件失败: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 添加其他字段
	writer.WriteField("model", "FunAudioLLM/SenseVoiceSmall")
	writer.WriteField("response_format", "json")

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭 writer 失败: %w", err)
	}

	UpdateTask(taskID, "processing", 30, "上传音频文件...")

	// 发送请求
	req, err := http.NewRequest("POST", "https://api.siliconflow.cn/v1/audio/transcriptions", body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 300 * 1e9} // 5分钟超时
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败: %s", string(respBody))
	}

	UpdateTask(taskID, "processing", 80, "处理转写结果...")

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	transcription := &TranscriptionResult{}

	if text, ok := result["text"].(string); ok {
		transcription.Text = text
	}

	// 获取音频时长
	duration, _ := utils.GetVideoDuration(audioPath)
	transcription.Duration = duration

	UpdateTask(taskID, "completed", 100, "语音转文字完成")

	return transcription, nil
}
