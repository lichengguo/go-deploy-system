package v1

import (
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 发布代码

// GetUserAllDeployment 获取当前登录用户的所属项目列表
func GetUserAllDeployment(c *gin.Context) {
	// 获取用户名
	userName, _ := c.Get("name")
	userNameStr := userName.(string) // interface 转为 string

	code, data := model.GetUserAllDeployment(userNameStr)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GitPullDeployment 拉取即将发布的项目代码到本地
func GitPullDeployment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// 获取用户名
	userName, _ := c.Get("name")
	userNameStr := userName.(string) // interface 转为 string

	code, data := model.GitPullDeployment(id, userNameStr)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"data":    data,
	})
}

// ReleaseToServer 发布代码到远程服务器
func ReleaseToServer(c *gin.Context) {
	var (
		data model.RecDeploymentLog
		msg  string
		code int64
	)

	// 获取用户名
	userName, _ := c.Get("name")
	userNameStr := userName.(string) // interface 转为 string

	_ = c.ShouldBindJSON(&data)

	// 对客户端提交过来的数据进行验证
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return // 字段校验未通过直接结束函数
	}

	code = model.ReleaseToServer(&data, userNameStr)

	// todo 增加邮件提示

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// RollBackCode 回滚代码
func RollBackCode(c *gin.Context) {
	deployLogID, _ := strconv.Atoi(c.Param("id"))

	// 获取用户名
	userName, _ := c.Get("name")
	userNameStr := userName.(string) // interface 转为 string

	code := model.RollBackCode(deployLogID, userNameStr)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
