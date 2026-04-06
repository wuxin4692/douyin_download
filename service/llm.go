package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"douyin-tool/config"
)

type LLMResult struct {
	Analysis string `json:"analysis"`
}

type RewriteResult struct {
	Original string `json:"original"`
	Rewritten string `json:"rewritten"`
	Style    string `json:"style"`
}

// AnalyzeText 使用 LLM 分析文案
func AnalyzeText(text string, taskID string, apiKey string, provider string) (*LLMResult, error) {
	UpdateTask(taskID, "processing", 10, "开始分析文案...")

	if apiKey == "" {
		apiKey = config.AppConfig.MiniMaxKey
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未设置 API Key")
	}

	prompt := buildAnalysisPrompt(text)
	
	result, err := callLLM(prompt, apiKey, provider)
	if err != nil {
		UpdateTask(taskID, "failed", 0, err.Error())
		return nil, err
	}

	UpdateTask(taskID, "completed", 100, "分析完成")
	
	return &LLMResult{
		Analysis: result,
	}, nil
}

// RewriteText 一键改写文案
func RewriteText(text string, taskID string, apiKey string, provider string, style string, customInstruction string) (*RewriteResult, error) {
	UpdateTask(taskID, "processing", 10, "开始改写文案...")

	if apiKey == "" {
		apiKey = config.AppConfig.MiniMaxKey
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未设置 API Key")
	}

	prompt := buildRewritePrompt(text, style, customInstruction)
	
	rewritten, err := callLLM(prompt, apiKey, provider)
	if err != nil {
		UpdateTask(taskID, "failed", 0, err.Error())
		return nil, err
	}

	UpdateTask(taskID, "completed", 100, "改写完成")
	
	return &RewriteResult{
		Original: text,
		Rewritten: rewritten,
		Style:    style,
	}, nil
}

func buildAnalysisPrompt(text string) string {
	return fmt.Sprintf(`请分析以下抖音文案，从内容结构、情感基调、受众定位、爆款元素等维度进行解读：

---
%s
---

请给出专业的分析报告，帮助理解这篇文案为何有效或需要改进的地方。`, text)
}

func buildRewritePrompt(text string, style string, customInstruction string) string {
	styleDesc := map[string]string{
		"inspiring": "激励人心、充满正能量的风格",
		"humorous":  "幽默风趣、轻松活泼的风格",
		"professional": "专业严谨、权威可靠的风格",
		"casual":    "日常随意、亲近友好的风格",
		"emotional": "情感丰富、打动人心的风格",
	}
	
	desc, ok := styleDesc[style]
	if !ok {
		desc = "更适合短视频传播的风格"
	}
	
	prompt := fmt.Sprintf(`请将以下文案改写成"%s"风格，保持核心信息不变，但调整表达方式使其更符合目标风格：

---
%s
---

`, desc, text)
	
	// 添加自定义指令
	if customInstruction != "" {
		prompt += fmt.Sprintf(`额外要求：
%s

`, customInstruction)
	}
	
	prompt += `请直接输出改写后的文案，不需要额外说明。`
	
	return prompt
}

func callLLM(prompt string, apiKey string, provider string) (string, error) {
	// 默认使用配置的 provider
	if provider == "" {
		provider = config.AppConfig.AIProvider
	}

	// 获取模型名称，优先使用配置的模型
	model := config.AppConfig.AIModel
	if model == "" {
		model = "gpt-3.5-turbo" // OpenAI 默认
	}

	switch provider {
	case "openai":
		return callOpenAICompatible(prompt, apiKey, "https://api.openai.com/v1/chat/completions", model)
	case "minimax":
		if model == "gpt-3.5-turbo" || model == "" {
			model = "MiniMax-M2.5" // MiniMax 默认
		}
		return callMiniMax(prompt, apiKey, model)
	case "compatible":
		baseURL := config.AppConfig.AIBaseURL
		if model == "" {
			model = "gpt-3.5-turbo"
		}
		if baseURL == "" {
			return "", fmt.Errorf("请先配置 API 接口地址")
		}
		return callOpenAICompatible(prompt, apiKey, baseURL, model)
	default:
		return callOpenAICompatible(prompt, apiKey, "https://api.openai.com/v1/chat/completions", model)
	}
}

func callMiniMax(prompt string, apiKey string, model string) (string, error) {
	url := "https://api.minimax.chat/v1/text/chatcompletion_v2"
	
	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 4096,
	}
	
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("请求构建失败: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 120 * 1e9}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 请求失败: %s", string(respBody))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}
	
	// 提取响应内容
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["messages"].([]interface{}); ok && len(msg) > 0 {
				if message, ok := msg[len(msg)-1].(map[string]interface{}); ok {
					if content, ok := message["text"].(string); ok {
						return strings.TrimSpace(content), nil
					}
				}
			}
		}
	}
	
	return "", fmt.Errorf("无法解析响应")
}

// callOpenAICompatible 调用 OpenAI 兼容接口
func callOpenAICompatible(prompt string, apiKey string, baseURL string, model string) (string, error) {
	// 确保 baseURL 格式正确
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	url := baseURL + "chat/completions"
	
	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 4096,
	}
	
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("请求构建失败: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 120 * 1e9}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 请求失败: %s", string(respBody))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}
	
	// 提取响应内容 (OpenAI 兼容格式)
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := msg["content"].(string); ok {
					return strings.TrimSpace(content), nil
				}
			}
		}
	}
	
	return "", fmt.Errorf("无法解析响应")
}
