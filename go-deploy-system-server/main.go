package main

import (
	"go-deploy-system-server/model"
	"go-deploy-system-server/routes"
)

func main() {
	// 初始化数据库
	model.InitDb()

	// 引用路由
	routes.InitRouter()

}
