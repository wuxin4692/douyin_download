package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"douyin-tool/config"
)

type SegmentResult struct {
	Segments []Segment `json:"segments"`
	FullText string    `json:"full_text"`
}

type Segment struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// SemanticSegment performs semantic segmentation on text
func SemanticSegment(text string, taskID string, apiKey string) (*SegmentResult, error) {
	UpdateTask(taskID, "processing", 0, "开始语义分段...")

	if apiKey == "" {
		apiKey = config.AppConfig.MiniMaxKey
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未设置 MiniMax API Key")
	}

	// 构建请求
	prompt := `请对以下文本进行语义分段。每个段落需要添加一个简洁的小标题（用##开头）。

要求：
1. 保持原文语义完整性
2. 小标题简洁明了，概括该段主题
3. 如果内容较短（少于100字），可以不分段
4. 只返回分段结果，不要其他解释

文本内容：
` + text

	requestBody := map[string]interface{}{
		"model": "MiniMax-Text-01",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":   4096,
		"temperature":  0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	UpdateTask(taskID, "processing", 30, "调用 MiniMax API...")

	req, err := http.NewRequest("POST", "https://api.minimax.chat/v1/text/chatcompletion_pro", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
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

	UpdateTask(taskID, "processing", 70, "处理分段结果...")

	// 解析响应
	var apiResp map[string]interface{}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取文本内容
	choices, ok := apiResp["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("无效的 API 响应")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("无效的响应格式")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("无效的消息格式")
	}

	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("无效的内容格式")
	}

	// 解析分段结果
	result := parseSegments(content)
	result.FullText = text

	UpdateTask(taskID, "completed", 100, "语义分段完成")

	return result, nil
}

func parseSegments(text string) *SegmentResult {
	result := &SegmentResult{
		Segments: []Segment{},
	}

	lines := strings.Split(text, "\n")
	var currentTitle string
	var currentContent strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 检查是否是标题行
		if strings.HasPrefix(trimmed, "##") {
			// 保存之前的段落
			if currentContent.Len() > 0 {
				result.Segments = append(result.Segments, Segment{
					Title:   currentTitle,
					Content: strings.TrimSpace(currentContent.String()),
				})
				currentContent.Reset()
			}
			currentTitle = strings.TrimPrefix(trimmed, "##")
			currentTitle = strings.TrimSpace(currentTitle)
		} else if trimmed != "" {
			currentContent.WriteString(trimmed)
			currentContent.WriteString("\n")
		}
	}

	// 保存最后一个段落
	if currentContent.Len() > 0 {
		title := currentTitle
		if title == "" {
			title = "全文"
		}
		result.Segments = append(result.Segments, Segment{
			Title:   title,
			Content: strings.TrimSpace(currentContent.String()),
		})
	}

	return result
}
