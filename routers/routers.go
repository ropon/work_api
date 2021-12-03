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
	gin.SetMode(gin.ReleaseMode)
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	<-quit
	logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown error: %s", err.Error())
	}
	logger.Info("Server exiting")
}
