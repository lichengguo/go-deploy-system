package model

import (
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/scryptpwd"

	"gorm.io/gorm"
)

// 用户管理数据模型

type User struct {
	gorm.Model
	Department   Department `gorm:"foreignkey:DepartmentID" validate:"-"` // 外键 部门ID  "-" 跳过验证
	UserName     string     `gorm:"type:varchar(20);not null" json:"user_name" validate:"required,min=4,max=20" label:"用户名"`
	Password     string     `gorm:"type:varchar(100);not null" json:"password" validate:"required,min=6,max=100" label:"密码"`
	Role         int        `gorm:"type:int;not null;comment:'1 管理员 2 普通用户'" json:"role" validate:"required,gte=1" label:"权限"`
	Status       int        `gorm:"type:int;not null;comment:'1 可用 2 禁用'" json:"status" validate:"required,gte=1" label:"状态"`
	DepartmentID int        `gorm:"type:int;not null" json:"department_id" validate:"required" label:"部门ID"`
	UpdatedAt    int
	CreatedAt    int
}

// 接收客户端修改密码的数据
type RceUserPassData struct {
	OldPwd string `json:"old_pwd" validate:"required,min=6,max=100" label:"旧密码"`
	NewPwd string `json:"new_pwd" validate:"required,min=6,max=100" label:"新密码"`
}

// 添加用户
func AddUser(data *User) int64 {
	var (
		user  User
		total int64
	)

	// 添加用户前，先检查用户是否存在
	db.Model(&user).Where("user_name = ?", data.UserName).Count(&total)
	if total > 0 {
		return errmsg.ERROR_USER_EXIST
	}

	// 密码加密
	data.Password = scryptpwd.ScryptPwd(data.Password)

	// 添加用户
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR_USER_CREATED_FAIL
	}

	return errmsg.SUCCESS
}

// 删除用户
func DelUser(id int) int64 {
	var (
		user             User
		deploymentToUser DeploymentToUserRole
	)

	// 如果使用了gorm.Model，会在数据库软删除
	tx := db.Begin()                                  // 开始事务
	err := tx.Where("id = ?", id).Delete(&user).Error // 删除用户表的数据
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}
	err = tx.Where("user_id = ?", id).Delete(&deploymentToUser).Error // 删除项目-用户表的数据
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}
	tx.Commit() // 提交事务

	return errmsg.SUCCESS
}

// 修改用户
func ModUser(id int, data *User) int64 {
	var (
		user     User
		total    int64
		oldPwd   string // 数据库存在的密码
		userMaps = make(map[string]interface{})
	)

	// 检查修改后的用户名是否存在多个
	db.Model(&user).Where("id <> ? AND user_name = ?", id, data.UserName).Count(&total)
	if total > 0 {
		return errmsg.ERROR_USER_EXIST
	}

	// 使用Updates方法更新数据
	// 当通过 struct 更新时，GORM 只会更新非零字段。
	// 如果您想确保指定字段被更新，你应该使用 Select 更新选定字段，或使用 map 来完成更新操作
	userMaps["user_name"] = data.UserName
	userMaps["role"] = data.Role
	userMaps["status"] = data.Status
	db.Model(&user).Select("password").Where("id = ?", id).First(&oldPwd) // 查询数据库中的密码
	if data.Password != oldPwd {
		// 如果传递过来的密码和数据库中的密码不一致则改变
		data.Password = scryptpwd.ScryptPwd(data.Password)
	}
	userMaps["password"] = data.Password
	userMaps["department_id"] = data.DepartmentID

	err := db.Model(&user).Where("id = ?", id).Updates(userMaps).Error
	if err != nil {
		return errmsg.ERROR_USER_UPDATE_FAIL // 更新用户信息失败
	}

	return errmsg.SUCCESS
}

// 用户列表、搜索
func FindUserList(pageSize, page int, userName string) (int64, []User, int64) {
	var (
		userList []User
		total    int64
	)

	// 搜索
	if userName != "" {
		// 预加载 department 表
		db.Preload("Department").Where("user_name LIKE ?", userName+"%").Find(&userList).Count(&total)
		err := db.Preload("Department").Where("user_name LIKE ?", userName+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&userList).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return 0, nil, errmsg.ERROR
		}
		return total, userList, errmsg.SUCCESS
	}

	// 多表查询 预加载
	db.Preload("Department").Find(&userList).Count(&total)
	err := db.Preload("Department").Limit(pageSize).Offset((page - 1) * pageSize).Find(&userList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, nil, errmsg.ERROR
	}

	return total, userList, errmsg.SUCCESS
}

// 查询用户权限
func FindUserRole(userName interface{}) (int64, int) {
	var (
		user User
		role int
	)

	err := db.Model(&user).Select("role").Where("user_name = ?", userName).First(&role).Error
	if err != nil {
		return errmsg.ERROR, 0
	}

	return errmsg.SUCCESS, role
}

// 修改密码
func ChangePwd(data *RceUserPassData, userName string) int64 {
	var (
		oldPwdStr string
		user      User
		userMaps  = make(map[string]interface{})
	)

	err := db.Model(&user).Select("password").Where("user_name = ?", userName).First(&oldPwdStr).Error
	if err != nil {
		return errmsg.ERROR_USER_NOT_EXIST // 用户不不存在
	}

	recOldPwd := scryptpwd.ScryptPwd(data.OldPwd) // 传递过来的明文密码加密后再比较
	if recOldPwd == oldPwdStr {                   // 判断传递过来的旧密码和数据库的密码是否一样
		recNewPwd := scryptpwd.ScryptPwd(data.NewPwd) // 新密码加密处理
		userMaps["password"] = recNewPwd

		err := db.Model(&user).Where("user_name = ?", userName).Updates(userMaps).Error
		if err != nil {
			return errmsg.ERROR_USER_UPDATE_FAIL
		}

		return errmsg.SUCCESS
	}

	return errmsg.ERROR_USER_PASSWORD_FAIL // 旧密码错误
}
