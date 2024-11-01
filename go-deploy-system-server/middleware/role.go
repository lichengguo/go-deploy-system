package middleware

import (
	"github.com/gin-gonic/gin"
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"net/http"
)

// 用户权限钩子

// RoleUser 根据token判断用户权限
func RoleUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		name, _ := c.Get("name")

		// 查询用户权限
		code, role := model.FindUserRole(name)
		if code != errmsg.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 1 管理员 2 普通用户
		if role == 1 {
			c.Next()
		} else {
			code = errmsg.ERROR_USER_ROLE_FAIL // 用户权限不够
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
	}
}
