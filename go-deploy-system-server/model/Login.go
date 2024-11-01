package model

import (
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/scryptpwd"
)

// 登录数据模型

// 检查登录
func CheckLogin(userName, password string) int64 {
	var user User
	db.Where("user_name = ? ", userName).First(&user)

	// 验证用户名
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST // 用户不存在
	}

	// 验证用户状态
	if user.Status != 1 {
		return errmsg.ERROR_USER_DISABLE // 用户已被禁用
	}

	// 验证密码
	if scryptpwd.ScryptPwd(password) != user.Password {
		return errmsg.ERROR_USER_PASSWORD_FAIL // 用户密码错误
	}

	return errmsg.SUCCESS
}
