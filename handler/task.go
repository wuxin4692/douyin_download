package handler

import (
	"net/http"

	"douyin-tool/service"

	"github.com/gin-gonic/gin"
)

type TranscribeRequest struct {
	AudioPath string `json:"audio_path" binding:"required"`
	APIKey    string `json:"api_key"`
}

func Transcribe(c *gin.Context) {
	var req TranscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供音频路径",
		})
		return
	}

	task := service.CreateTask()

	go func() {
		result, err := service.TranscribeAudio(req.AudioPath, task.ID, req.APIKey)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"task_id": task.ID,
		"data":    task,
	})
}

type SegmentRequest struct {
	Text   string `json:"text" binding:"required"`
	APIKey string `json:"api_key"`
}

func Segment(c *gin.Context) {
	var req SegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供文本内容",
		})
		return
	}

	task := service.CreateTask()

	go func() {
		result, err := service.SemanticSegment(req.Text, task.ID, req.APIKey)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, result)
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"task_id": task.ID,
		"data":    task,
	})
}

func GetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	task := service.GetTask(taskID)

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": task,
	})
}

type SaveFileRequest struct {
	Content    string `json:"content" binding:"required"`
	Filename   string `json:"filename" binding:"required"`
}

func SaveFile(c *gin.Context) {
	var req SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供文件内容和文件名",
		})
		return
	}

	filepath, err := service.SaveToFile(req.Content, req.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "保存文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "文件已保存",
		"data": gin.H{
			"path": filepath,
		},
	})
}
