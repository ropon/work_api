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

// CreateService 创建服务接口
// @Summary 创建服务接口
// @Description 创建服务接口
// @Tags 服务相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUServiceReq true "请求参数"
// @Success 200 {object} models.Service "创建成功返回结果"
// @Router /work_api/api/v1/service [post]
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

// DeleteService 删除服务接口
// @Summary 删除服务接口
// @Description 删除服务接口
// @Tags 服务相关接口
// @Produce application/json
// @Param id path uint true "id"
// @Success 200
// @Router /work_api/api/v1/service/{id} [delete]
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

// UpdateService 更新服务全部参数接口
// @Summary 更新服务全部参数接口
// @Description 更新服务全部参数接口
// @Tags 服务相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.CUServiceReq true "请求参数"
// @Success 200 {object} models.Service "更新成功返回结果"
// @Router /work_api/api/v1/service [put]
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

// PatchUpdateService 更新服务部分参数接口
// @Summary 更新服务部分参数接口
// @Description 更新服务部分参数接口
// @Tags 服务相关接口
// @Accept application/json
// @Produce application/json
// @Param data body logics.ServiceReq true "请求参数"
// @Success 200 {object} models.Service "更新成功返回结果"
// @Router /work_api/api/v1/service [patch]
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

// GetServices 获取服务列表接口
// @Summary 获取服务列表接口
// @Description 获取服务列表接口
// @Tags 服务相关接口
// @Produce application/json
// @Param data query logics.ServiceReq true "请求参数"
// @Success 200 {object} logics.ServiceRes "服务列表返回结果"
// @Router /work_api/api/v1/services [get]
func GetServices(c *gin.Context) {
	req := new(logics.ServiceReq)
	if !checkData(c, req) {
		return
	}

	ctx := utils.ExtractStdContext(nil, c.Request.Header)
	resList, err := logics.GetServices(ctx, req)
	if err != nil {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, err.Error())
		return
	}
	utils.GinOKRsp(c, resList, "获取列表成功")
}

// GetService 获取单个服务接口
// @Summary 获取单个服务接口
// @Description 获取单个服务接口
// @Tags 服务相关接口
// @Produce application/json
// @Param id path uint true "id"
// @Success 200 {object} models.Service "服务返回结果"
// @Router /work_api/api/v1/service/{id} [get]
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
