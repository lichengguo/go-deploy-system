package v1

import (
	"go-deploy-system-server/middleware"
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登录接口

// 接收客户端登录的数据
type RcvLoginData struct {
	UserName string `json:"user_name" validate:"required,min=4,max=20" label:"登录用户名"`
	Password string `json:"password" validate:"required,min=6,max=20" label:"登录密码"`
}

// Login 登录函数
func Login(c *gin.Context) {
	var (
		data  RcvLoginData
		token string
		code  int64
		role  int // 用户权限 1 管理员; 2 普通用户
		msg   string
	)

	_ = c.ShouldBindJSON(&data) // 绑定客户端传递过来的数据到结构体

	// 对客户端提交过来的数据进行验证
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":    code,
			"message":   msg,
			"user_name": data.UserName,
			"token":     token,
			"role":      role,
		})
		return // 字段校验未通过直接结束函数
	}

	code = model.CheckLogin(data.UserName, data.Password)
	if code == errmsg.SUCCESS {
		// 登录成功
		token, code = middleware.SetToken(data.UserName) // 生成token
		_, role = model.FindUserRole(data.UserName)      // 获取当前登录用户权限
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    code,
		"message":   errmsg.GetErrMsg(code),
		"user_name": data.UserName,
		"token":     token,
		"role":      role,
	})
}
