package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workApi/config"
	"workApi/logic"
	"workApi/models"
)

// Auth
func Auth(c *gin.Context) {
	res := models.NewDefaultResult()
	token := c.Request.Header.Get("token")
	if token == "" {
		res.Code = 4001
		res.Msg = config.CodeType[res.Code]
		c.JSON(http.StatusOK, res)
		return
	}
	resMap, err := logic.CheckAuth(token)
	if err != nil {
		res.Code = 4003
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	res.Data = resMap
	c.JSON(http.StatusOK, res)
}
