package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workApi/models"
)

// Demo1
func Demo1(c *gin.Context) {
	res := models.NewDefaultResult()
	c.JSON(http.StatusOK, res)
}
