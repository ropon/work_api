package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ropon/work_api/utils"
	"strconv"
	"time"
)

func initExtraKeys(c *gin.Context) (string, string) {
	return c.Request.Header.Get("user_email"), c.Request.Header.Get("ops_admin")
}

//检查params id
func checkParamsId(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		utils.GinErrRsp(c, utils.ErrCodeParamInvalid, fmt.Sprintf("传入id[%s]不合法", idStr))
		return 0, false
	}
	return uint(id), true
}

//检查请求体
func checkData(c *gin.Context, v interface{}) bool {
	if err := c.ShouldBind(v); err != nil {
		utils.GinErrRsp(c, utils.ErrCodeParamInvalid, "参数有误:"+err.Error())
		return false
	}
	return true
}

func Hi(c *gin.Context) {
	time.Sleep(time.Second * 10)
	utils.GinOKRsp(c, "hi ropon", "ok")
}
