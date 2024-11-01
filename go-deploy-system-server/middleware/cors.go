package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// 跨域
func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowAllOrigins:  true,                                                        // 允许所有的跨域请求
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},         // 允许跨域请求使用的方法
			AllowHeaders:     []string{"*"},                                               // 允许客户端使用简单的头部发起跨域请求
			ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type"}, // 允许的请求头
			AllowCredentials: true,                                                        // cookie ssl 认证
			MaxAge:           12 * time.Hour,                                              // 缓存时间
		},
	)
}
