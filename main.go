package main

import (
	"fmt"
	"workApi/config"
	"workApi/dblayer"
	"workApi/models"
	"workApi/routers"
)

func main() {
	//初始化配置文件
	err := config.InitConf("./config.conf")
	if err != nil {
		fmt.Println("解析配置文件异常")
		return
	}
	//连接redis
	err = dblayer.InitRedis()
	if err != nil {
		fmt.Println("连接redis数据库异常")
		return
	}
	//连接mysql
	err = dblayer.InitMySQL()
	if err != nil {
		fmt.Println("连接mysql数据库异常")
		return
	}
	defer dblayer.DB.Close()
	//迁移数据库
	dblayer.DB = dblayer.DB.AutoMigrate(&models.User{}, &models.TaskLog{}, &models.Order{})
	//初始化路由
	r := routers.SetupRouter()
	// 启动server
	_ = r.Run(config.GHostPort())
}
