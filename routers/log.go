/*
Author:Ropon
Date:  2020-12-22
*/
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/logger"
	"github.com/ropon/work_api/utils"
	"time"
)

func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := logger.ColorForStatus(statusCode)
		methodColor := logger.ColorForMethod(method)
		resetColor := logger.Reset
		userEmail := c.Request.Header.Get("user_email")
		requestData := utils.GetRequestData(c)
		logger.Info("[GIN] %s%s%s %s%s %s%d%s cost:%.03f [ip:%s] [user_email:%s]",
			methodColor, method, resetColor,
			c.Request.Host,
			utils.Cuts(requestData, utils.MaxLogDataLen),
			statusColor, statusCode, resetColor,
			end.Sub(start).Seconds(),
			clientIP,
			userEmail)
	}
}
