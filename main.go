package main

import (
	"fmt"
	"github.com/ropon/work_api/conf"
	"github.com/ropon/work_api/routers"
)

// @title work_api
// @version 1.0
// @description 后端快速Api脚手架

// @contact.name Ropon
// @contact.url https://www.ropon.top
// @contact.email ropon@xxx.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host work-api.xxx.com
// @BasePath /
func main() {
	err := conf.Init()
	if err != nil {
		fmt.Printf("init failed, err: %v\n", err)
		return
	}

	routers.Run(conf.Cfg.Listen)
}
