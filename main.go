package main

import (
	"fmt"
	"github.com/ropon/work_api/conf"
	"github.com/ropon/work_api/routers"
)

func main() {
	err := conf.Init()
	if err != nil {
		fmt.Printf("init failed, err: %v\n", err)
		return
	}

	routers.Run()
}
