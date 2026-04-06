package handler

import (
	"net/http"

	"douyin-tool/service"

	"github.com/gin-gonic/gin"
)

type ParseRequest struct {
	URL string `json:"url" binding:"required"`
}

func ParseVideo(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供有效的抖音链接",
		})
		return
	}

	// 创建任务
	task := service.CreateTask()

	// 异步解析
	go func() {
		videoInfo, err := service.ParseShareUrl(req.URL)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}
		service.SetTaskResult(task.ID, videoInfo)
		service.UpdateTask(task.ID, "completed", 100, "解析完成")
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"task_id": task.ID,
		"data":    task,
	})
}

func DownloadVideo(c *gin.Context) {
	var req struct {
		URL    string `json:"url"`
		TaskID string `json:"task_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "参数错误",
		})
		return
	}

	task := service.CreateTask()

	go func() {
		var videoInfo *service.VideoInfo
		var err error

		if req.TaskID != "" {
			// 从已有任务获取视频信息
			parentTask := service.GetTask(req.TaskID)
			if parentTask != nil && parentTask.Result != nil {
				videoInfo, _ = parentTask.Result.(*service.VideoInfo)
			}
		}

		if videoInfo == nil {
			// 解析链接
			videoInfo, err = service.ParseShareUrl(req.URL)
			if err != nil {
				service.SetTaskError(task.ID, err.Error())
				service.UpdateTask(task.ID, "failed", 0, err.Error())
				return
			}
		}

		// 下载视频
		result, err := service.DownloadVideo(videoInfo, task.ID)
		if err != nil {
			service.SetTaskError(task.ID, err.Error())
			service.UpdateTask(task.ID, "failed", 0, err.Error())
			return
		}

		service.SetTaskResult(task.ID, result)
		service.UpdateTask(task.ID, "completed", 100, "下载完成")
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"task_id": task.ID,
		"data":    task,
	})
}

func ExtractAudio(c *gin.Context) {
	var req struct {
		VideoPath string `json:"video_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请提供视频路径",
		})
		return
	}

	task := service.CreateTask()

	go func() {
		result, err := service.ExtractAudio(req.VideoPath, task.ID)
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
