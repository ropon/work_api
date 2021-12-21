package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ropon/work_api/utils"
)

func Hi(c *gin.Context) {
	utils.GinOKRsp(c, "hi ropon!", "ok")
}
