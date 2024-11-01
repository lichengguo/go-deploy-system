package routes

import (
	"github.com/gin-gonic/gin"
	v1 "go-deploy-system-server/api/v1"
	"go-deploy-system-server/middleware"
	"go-deploy-system-server/utils"
)

// InitRouter 初始化路由
func InitRouter() {
	gin.SetMode(utils.AppMode) // 设置gin启动模式，有debug、release两种
	//r := gin.Default()         // 创建路由
	//r.Use(middleware.Cors())   // 解决跨域
	r := gin.New()           // 创建路由
	r.Use(middleware.Log())  // 日志钩子
	r.Use(gin.Recovery())    // Gin框架默认的函数
	r.Use(middleware.Cors()) // 跨域

	// 任何用户都可以访问的路由
	router := r.Group("api/v1") // 创建路由组
	{
		router.POST("login/", v1.Login) // 登录
	}

	// 登录用户可以访问的路由
	release := r.Group("api/v1")       // 创建路由组
	release.Use(middleware.JwtToken()) // jwt登录中间件钩子
	{
		// 发布项目模块
		release.GET("releases", v1.GetUserAllDeployment)         // 获取当前登录用户的所拥有的项目列表
		release.PUT("release/gitpull/:id", v1.GitPullDeployment) // 拉取即将发布的项目代码到本地
		release.POST("release/add", v1.ReleaseToServer)          // 发布代码
		release.POST("release/rollback/:id", v1.RollBackCode)    // 回滚项目

		// 发布日志模块
		release.GET("deploymentlogs", v1.GetDeploymentLogs) // 发布日志列表、搜索

		release.POST("user/changepassword", v1.ChangePwd) // 修改密码
	}

	// 登录用户并且拥有管理权限才可以访问的路由
	auth := r.Group("api/v1")       // 创建路由组
	auth.GET("health", v1.Health)   // 心跳健康检测
	auth.Use(middleware.JwtToken()) // jwt登录中间件钩子
	auth.Use(middleware.RoleUser()) // 用户权限钩子
	{
		// 部门管理模块
		auth.POST("department", v1.AddDep)
		auth.DELETE("department/:id", v1.DelDep)
		auth.PUT("department/:id", v1.EditDep)
		auth.GET("department", v1.GetDeps)

		// 用户管理模块
		auth.POST("user", v1.AddUser)
		auth.DELETE("user/:id", v1.DelUser)
		auth.PUT("user/:id", v1.EditUser)
		auth.GET("user", v1.GetUsers)

		// 机房模块
		auth.POST("engineroom", v1.AddEngineroom)
		auth.DELETE("engineroom/:id", v1.DelEngineroom)
		auth.PUT("engineroom/:id", v1.EditEngineroom)
		auth.GET("engineroom", v1.GetEnginerooms)

		// 服务器模块
		auth.POST("server", v1.AddServer)
		auth.DELETE("server/:id", v1.DelServer)
		auth.PUT("server/:id", v1.EditServer)
		auth.GET("server", v1.GetServers)
		auth.GET("server/connect/:id", v1.ConnectServer)

		// 项目配置模块
		auth.POST("deployment", v1.AddDeployment)
		auth.DELETE("deployment/:id", v1.DelDeployment)
		auth.PUT("deployment/:id", v1.EditDeployment)
		auth.GET("deployment", v1.GetDeploymentList)

		auth.POST("upload", v1.UpLoad) // 秘钥上传
	}

	_ = r.Run(utils.HttpPort) // 启动程序
}
