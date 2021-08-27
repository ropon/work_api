package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ropon/work_api/logics"
	"github.com/ropon/work_api/utils"
	"strconv"
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
	if err := c.Bind(v); err != nil {
		utils.GinErrRsp(c, utils.ErrCodeParamInvalid, "参数有误")
		return false
	}
	return true
}

func CreateService(c *gin.Context) {
	req := new(logics.CUServiceReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.CreateService(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "创建成功")
}

func DeleteService(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	err := logics.DeleteService(id)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, "", "删除成功")
}

func UpdateService(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	req := new(logics.CUServiceReq)
	if !checkData(c, req) {
		return
	}

	req.Init(initExtraKeys(c))
	res, err := logics.UpdateService(id, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "更新成功")
}

func PatchUpdateService(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	req := new(logics.ServiceReq)
	if !checkData(c, req) {
		return
	}

	res, err := logics.PatchUpdateService(id, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "更新成功")
}

func GetServices(c *gin.Context) {
	req := new(logics.ServiceReq)
	if !checkData(c, req) {
		return
	}

	resList, err := logics.GetServices(req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, resList, "获取列表成功")
}

func GetService(c *gin.Context) {
	id, flag := checkParamsId(c)
	if !flag {
		return
	}

	res, err := logics.GetService(id)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, res, "获取成功")
}
