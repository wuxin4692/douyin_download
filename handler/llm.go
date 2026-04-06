package handler

import (
	"net/http"

	"douyin-tool/service"

	"github.com/gin-gonic/gin"
)

type AnalyzeRequest struct {
	Text     string `json:"text" binding:"required"`
	APIKey   string `json:"api_key"`
	Provider string `json:"provider"`
}

func AnalyzeText(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供文本内容",
		})
		return
	}

	task := service.CreateTask()

	go func() {
		result, err := service.AnalyzeText(req.Text, task.ID, req.APIKey, req.Provider)
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
	})
}

type RewriteRequest struct {
	Text             string `json:"text" binding:"required"`
	APIKey           string `json:"api_key"`
	Provider         string `json:"provider"`
	Style            string `json:"style"`
	CustomInstruction string `json:"custom_instruction"`
}

func RewriteText(c *gin.Context) {
	var req RewriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供文本内容",
		})
		return
	}

	task := service.CreateTask()

	go func() {
		result, err := service.RewriteText(req.Text, task.ID, req.APIKey, req.Provider, req.Style, req.CustomInstruction)
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
	})
}
