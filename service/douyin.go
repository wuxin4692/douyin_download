package service

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"douyin-tool/utils"
)

type VideoInfo struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	CoverURL     string `json:"cover_url"`
	DownloadURL  string `json:"download_url"`
	Author       string `json:"author"`
	Duration     int    `json:"duration"`
}

type Task struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"` // pending, processing, completed, failed
	Progress    int       `json:"progress"`
	Message     string    `json:"message"`
	Result      any       `json:"result,omitempty"`
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

var (
	tasks   = make(map[string]*Task)
	taskMux sync.RWMutex
)

func Init() {
}

func CreateTask() *Task {
	return &Task{
		ID:        generateID(),
		Status:    "pending",
		Progress:  0,
		CreatedAt: time.Now(),
	}
}

func GetTask(id string) *Task {
	taskMux.RLock()
	defer taskMux.RUnlock()
	return tasks[id]
}

func UpdateTask(id string, status string, progress int, message string) {
	taskMux.Lock()
	defer taskMux.Unlock()
	if task, ok := tasks[id]; ok {
		task.Status = status
		task.Progress = progress
		task.Message = message
		if status == "completed" || status == "failed" {
			task.CompletedAt = time.Now()
		}
	}
}

func SetTaskResult(id string, result any) {
	taskMux.Lock()
	defer taskMux.Unlock()
	if task, ok := tasks[id]; ok {
		task.Result = result
	}
}

func SetTaskError(id string, err string) {
	taskMux.Lock()
	defer taskMux.Unlock()
	if task, ok := tasks[id]; ok {
		task.Error = err
	}
}

func generateID() string {
	return time.Now().Format("20060102150405") + strings.ToLower(randomString(6))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// ParseShareUrl parses Douyin share link or modal_id
func ParseShareUrl(shareText string) (*VideoInfo, error) {
	url := extractURL(shareText)
	if url == "" {
		url = shareText
	}

	if isModalId(url) {
		return GetVideoInfoByModalId(url)
	}

	modalId := extractModalId(url)
	if modalId == "" {
		return nil, ErrInvalidLink
	}

	return GetVideoInfoByModalId(modalId)
}

func extractURL(text string) string {
	urlPattern := regexp.MustCompile(`https?://[^\s]+`)
	matches := urlPattern.FindString(text)
	return strings.TrimSpace(matches)
}

func isModalId(s string) bool {
	modalPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]{20,}$`)
	return modalPattern.MatchString(s)
}

func extractModalId(url string) string {
	patterns := []string{
		`/video/(\d+)`,
		`modal_id=(\d+)`,
		`(\d{19,})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// GetVideoInfoByModalId gets video info by modal_id
func GetVideoInfoByModalId(modalId string) (*VideoInfo, error) {
	apiURL := "https://www.iesdouyin.com/share/video/" + modalId

	resp, err := utils.HttpRequest(apiURL, map[string]interface{}{
		"headers": map[string][]string{
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		},
	})
	if err != nil {
		return nil, err
	}

	html := string(resp.Body)

	routerDataPattern := regexp.MustCompile(`window\._ROUTER_DATA\s*=\s*(\{.*?\})\s*;?\s*</script>`)
	matches := routerDataPattern.FindStringSubmatch(html)
	if len(matches) < 2 {
		return nil, ErrParseFailed
	}

	var routerData map[string]interface{}
	jsonStr := matches[1]
	if err := json.Unmarshal([]byte(jsonStr), &routerData); err != nil {
		return nil, ErrParseFailed
	}

	videoInfo := extractVideoFromRouterData(routerData)
	if videoInfo == nil {
		return nil, ErrVideoNotFound
	}

	return videoInfo, nil
}

func extractVideoFromRouterData(data map[string]interface{}) *VideoInfo {
	if video, ok := data["video"].(map[string]interface{}); ok {
		videoInfo := &VideoInfo{}

		if v, ok := video["id"].(float64); ok {
			videoInfo.VideoID = strconv.FormatInt(int64(v), 10)
		}
		if v, ok := video["id_str"].(string); ok {
			videoInfo.VideoID = v
		}
		if v, ok := video["desc"].(string); ok {
			videoInfo.Title = v
		}
		if v, ok := video["author"].(map[string]interface{}); ok {
			if nickname, ok := v["nickname"].(string); ok {
				videoInfo.Author = nickname
			}
		}

		videoInfo.DownloadURL = extractDownloadURL(video)

		if cover, ok := video["cover"].(map[string]interface{}); ok {
			if urlList, ok := cover["url_list"].([]interface{}); ok && len(urlList) > 0 {
				if url, ok := urlList[0].(string); ok {
					videoInfo.CoverURL = url
				}
			}
		}

		return videoInfo
	}

	return nil
}

func extractDownloadURL(video map[string]interface{}) string {
	patterns := []string{"play_addr", "download_addr"}

	for _, pattern := range patterns {
		if addr, ok := video[pattern].(map[string]interface{}); ok {
			if urlList, ok := addr["url_list"].([]interface{}); ok && len(urlList) > 0 {
				if url, ok := urlList[0].(string); ok {
					return strings.Replace(url, "watermark=1", "watermark=0", -1)
				}
			}
		}
	}

	return ""
}
