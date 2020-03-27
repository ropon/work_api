package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workApi/config"
	"workApi/logic"
	"workApi/models"
)

func base(c *gin.Context) (res *models.Result) {
	res = models.NewDefaultResult()
	token := c.Request.Header.Get("token")
	if token == "" {
		res.Code = 4001
		res.Msg = config.CodeType[res.Code]
		return
	}
	resMap, err := logic.CheckAuth(token)
	if err != nil {
		res.Code = 4003
		res.Msg = err.Error()
		return
	}
	res.Data = resMap
	return
}

// Auth
func Auth(c *gin.Context) {
	res := base(c)
	c.JSON(http.StatusOK, res)
}

// Task
func Task(c *gin.Context) {
	res := base(c)
	//授权异常
	if res.Code != 0 {
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	uid := res.Data.(map[string]interface{})["uid"].(uint)
	res.Data = nil
	var taskLog models.TaskLog
	err := c.ShouldBindJSON(&taskLog)
	if err != nil {
		res.Code = 4001
		res.Msg = config.CodeType[res.Code]
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	taskLog.UID = uid
	err = logic.NewTaskLog(&taskLog)
	if err != nil {
		res.Code = 4004
		res.Msg = err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Data = taskLog
	c.JSON(http.StatusOK, res)
}
