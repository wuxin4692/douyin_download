package handler

import (
	"net/http"

	"douyin-tool/utils"

	"github.com/gin-gonic/gin"
)

func FFmpegStatus(c *gin.Context) {
	status := utils.GetFFmpegStatus()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": status,
	})
}
