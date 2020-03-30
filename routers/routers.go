package routers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"workApi/controller"
)

func SetupRouter() *gin.Engine {
	//记录日志
	gin.DisableConsoleColor()
	f, _ := os.Create("./workApi.log")
	gin.DefaultWriter = io.MultiWriter(f)
	//生产模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//跨域 频率控制
	r.Use(cors(), limitRate())
	v1 := r.Group("api/v1")
	{
		v1.POST("auth", controller.Auth)
		v1.POST("task", controller.Task)
		v1.POST("order", controller.Order)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 4404, "msg": "Page not found"})
	})
	return r
}
