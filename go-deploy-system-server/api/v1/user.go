package v1

import (
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户管理接口

// 添加用户
func AddUser(c *gin.Context) {
	var (
		data model.User
		msg  string
		code int64
	)

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

	code = model.AddUser(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除用户
func DelUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DelUser(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 修改用户
func EditUser(c *gin.Context) {
	var (
		data model.User
		msg  string
		code int64
	)

	id, _ := strconv.Atoi(c.Param("id"))
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

	code = model.ModUser(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 用户列表、搜索
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	page, _ := strconv.Atoi(c.Query("page"))
	userName := c.Query("user_name") // 搜索

	total, data, code := model.FindUserList(pageSize, page, userName)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 修改用户密码
func ChangePwd(c *gin.Context) {
	var (
		data model.RceUserPassData
		msg  string
		code int64
	)

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

	userName, _ := c.Get("name")     // 获取用户名
	userNameStr := userName.(string) // interface 转为 string

	code = model.ChangePwd(&data, userNameStr)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})

}
