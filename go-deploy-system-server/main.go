/*
Go在MacOS上打包成Win及Linux运行文件
Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o yf_deploy_system_go_Linux

Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o yf_deploy_system_go_Win

Mac
go build -o yf_deploy_system_go_Mac
*/
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
