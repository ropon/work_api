package main

import (
	"workApi/config"
	"workApi/dblayer"
	"workApi/models"
	"workApi/routers"
)

func main() {
	//初始化配置文件
	err := config.InitConf("./config/config.conf")
	if err != nil {
		return
	}
	//连接redis
	err = dblayer.InitRedis()
	if err != nil {
		return
	}
	//连接mysql
	err = dblayer.InitMySQL()
	if err != nil {
		return
	}
	defer dblayer.DB.Close()
	//迁移数据库
	dblayer.DB = dblayer.DB.AutoMigrate(&models.User{})
	//初始化路由
	r := routers.SetupRouter()
	// 启动server
	_ = r.Run(config.GHostPort())
}
