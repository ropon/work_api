package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/logger"
	"github.com/ropon/work_api/utils"
)

func Hi(c *gin.Context) {
	logger.Debug("这是一条debug日志")
	logger.Error("这是一条错误日志%s", c.ClientIP())
	s := []string{"1","2"}
	utils.GinOKRsp(c, s[2], "ok")
}
