package routers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ropon/logger"
	"github.com/ropon/work_api/conf"
	"github.com/ropon/work_api/controllers"
	"github.com/ropon/work_api/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// setupRouter 初始化路由
func setupRouter() *gin.Engine {
	engine := gin.New()
	engine.Use(ginLogger())
	engine.Use(cors())

	v1 := engine.Group("/work_api/api/v1")
	{
		v1.GET("/hi", controllers.Hi)

		v1.POST("/service", controllers.CreateService)
		v1.DELETE("/service/:id", controllers.DeleteService)
		v1.PUT("/service/:id", controllers.UpdateService)
		v1.PATCH("/service/:id", controllers.PatchUpdateService)
		v1.GET("/service", controllers.GetServices)
		v1.GET("/service/:id", controllers.GetService)
	}

	engine.NoRoute(func(c *gin.Context) {
		utils.GinErrRsp(c, utils.ErrCodeGeneralFail, "Page not found")
	})
	return engine
}

func Run() {
	srv := &http.Server{
		Addr:    conf.Cfg.Listen,
		Handler: setupRouter(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen error: %s\n", err.Error())
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown error: %s", err.Error())
	}
	logger.Info("Server exiting")
}
