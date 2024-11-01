package model

import (
	"fmt"
	"go-deploy-system-server/utils/aespwd"
	"go-deploy-system-server/utils/errmsg"

	"gorm.io/gorm"
)

// 发布项目配置数据模型
type Deployment struct {
	// GitUrlHttp 和 GitUrlSsh 不能互为空，这两个字段的其中一个必须填写
	// 如果填写了 GitUrlHttp 字段，则必须填写 GitUser、GitPasswd 字段
	// 如果填写了 GitUrlSsh 字段，则必须填写 GitKey 字段
	gorm.Model
	DeployName       string `gorm:"type:varchar(100)" json:"deploy_name" validate:"required,max=100" label:"发布项目名称"`
	GitUrlHttp       string `gorm:"type:varchar(256)" json:"git_url_http" validate:"required_without=GitUrlSsh" label:"githttp链接"`
	GitUrlSsh        string `gorm:"type:varchar(256)" json:"git_url_ssh" validate:"required_without=GitUrlHttp" label:"gitssh链接"`
	GitBranch        string `gorm:"type:varchar(50);default:master" json:"git_branch" validate:"-" label:"git分支"`
	GitUser          string `gorm:"type:varchar(100)" json:"git_user" validate:"required_with=GitUrlHttp " label:"git账号"`
	GitPasswd        string `gorm:"type:varchar(100)" json:"git_passwd" validate:"required_with=GitUrlHttp required_without=GitKey" label:"git密码"`
	GitKey           string `gorm:"type:varchar(100)" json:"git_key" validate:"required_with=GitUrlSsh required_without=GitPasswd" label:"git秘钥"`
	DeployServerPath string `gorm:"type:varchar(100)" json:"deploy_server_path" validate:"required,max=100" label:"服务器的发布目录"`
	UpdatedAt        int    // 更新时间
	CreatedAt        int    // 创建时间

}

// 发布项目和服务器对应数据模型
type DeploymentToServer struct {
	gorm.Model
	Deployment   Deployment `gorm:"foreignkey:DeploymentID"`
	Server       Server     `gorm:"foreignkey:ServerID"`
	DeploymentID int        `gorm:"type:int" json:"deployment_id" validate:"required" label:"发布项目ID"`
	ServerID     int        `gorm:"type:int" json:"server_id" validate:"required" label:"发布服务器ID"`
	UpdatedAt    int        // 更新时间
	CreatedAt    int        // 创建时间
}

// 发布项目和发布用户对应表数据模型
type DeploymentToUserRole struct {
	gorm.Model
	Deployment   Deployment `gorm:"foreignkey:DeploymentID"`
	User         User       `gorm:"foreignkey:UserID" validate:"-"`
	DeploymentID int        `gorm:"type:int" json:"deployment_id" validate:"required" label:"发布项目ID"`
	UserID       int        `gorm:"type:int" json:"user_id" validate:"required" label:"发布用户ID"`
	UpdatedAt    int        // 更新时间
	CreatedAt    int        // 创建时间
}

// 接收客户端传递过来的数据
type DeploymentAcceptData struct {
	Deployment
	Server   Server `gorm:"foreignkey:ServerID" validate:"-"`
	ServerID []int  `gorm:"type:int;not null" json:"server_id" validate:"required" label:"发布服务器ID"`
	User     User   `gorm:"foreignkey:UserID" validate:"-"`
	UserID   []int  `gorm:"type:int;not null" json:"user_id" validate:"required" label:"发布用户ID"`
}

// 返回给客户端的数据
type DeploymentSentData struct {
	ID uint
	Deployment
	ServerList []*Server
	UserList   []*User
}

// 添加发布项目配置
func AddDeployment(data *DeploymentAcceptData) int64 {
	// 添加发布项目配置前，先检查发布项目配置是否存在
	var deployment Deployment
	db.Select("id").Where("deploy_name = ?", data.DeployName).First(&deployment)
	if deployment.ID > 0 {
		return errmsg.ERROR_DEPLOYMENT_EXIST // 发布项目配置已经存在
	}

	// 数据写入到数据库
	// 这里应该作为一组事务提交
	tx := db.Begin() // 开始事务
	// 发布项目配置表
	deployment.DeployName = data.DeployName             // 发布项目名称
	deployment.GitUrlHttp = data.GitUrlHttp             // githttp链接
	deployment.GitUrlSsh = data.GitUrlSsh               // gitssh链接
	deployment.GitBranch = data.GitBranch               // git分支
	deployment.GitUser = data.GitUser                   // git账号
	deployment.GitKey = data.GitKey                     // git秘钥
	deployment.DeployServerPath = data.DeployServerPath // 服务器发布目标目录
	if data.GitPasswd != "" {
		pwd, _ := aespwd.EnPwdCode(data.GitPasswd) // 密码加密存储，使用时再解密
		data.GitPasswd = pwd
	}
	deployment.GitPasswd = data.GitPasswd // git密码
	err := tx.Create(&deployment).Error   // 入库
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR_DEPLOYMENT_CREATED_FAIL
	}
	// 发布项目和服务器对应表
	// 一个发布项目可以有多台服务器
	for _, id := range data.ServerID {
		var deploymenttoserver DeploymentToServer
		deploymenttoserver.ServerID = id                     // 服务器ID
		deploymenttoserver.DeploymentID = int(deployment.ID) // 发布项目ID
		//deploymenttoserver.DeployServerPath = data.DeployServerPath // 服务器发布路径
		err := tx.Create(&deploymenttoserver).Error
		if err != nil {
			tx.Rollback()                               // 遇到错误时回滚事务
			return errmsg.ERROR_DEPLOYMENT_CREATED_FAIL // 发布项目已经存在
		}
	}
	// 发布项目和发布用户对应表
	// 一个项目可以有多个发布人
	for _, id := range data.UserID {
		var deploymenttouserrole DeploymentToUserRole
		deploymenttouserrole.UserID = id                       // 发布用户ID
		deploymenttouserrole.DeploymentID = int(deployment.ID) // 发布项目配置ID
		err := tx.Create(&deploymenttouserrole).Error
		if err != nil {
			tx.Rollback()                               // 遇到错误时回滚事务
			return errmsg.ERROR_DEPLOYMENT_CREATED_FAIL // 发布项目已经存在
		}
	}
	tx.Commit() // 提交事务
	return errmsg.SUCCESS
}

// 删除发布项目配置
func DelDeployment(id int) int64 {
	var deployment Deployment
	var deploymenttouserrole DeploymentToUserRole
	var deploymenttoserver DeploymentToServer

	tx := db.Begin() // 开始事务
	// 删除发布项目配置表的数据
	err := tx.Where("id = ?", id).Delete(&deployment).Error
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}

	// 删除发布项目和发布用户对应表中的数据
	err = tx.Where("deployment_id = ?", id).Delete(&deploymenttouserrole).Error
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}
	// 删除发布项目和服务器对应表中的数据
	err = tx.Where("deployment_id = ?", id).Delete(&deploymenttoserver).Error
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}

	tx.Commit() // 提交事务
	return errmsg.SUCCESS
}

// 修改发布项目配置
func ModDeployment(id int, data *DeploymentAcceptData) int64 {
	var deployment Deployment
	var deploymentTotal int64
	var oldPwd Deployment
	var deploymenttouserrole DeploymentToUserRole
	var deploymenttoserver DeploymentToServer

	// 修改前，检查发布项目的名称是否相同
	db.Where("id != ? AND deploy_name = ?", id, data.DeployName).First(&deployment).Count(&deploymentTotal)
	if deploymentTotal >= 1 {
		return errmsg.ERROR_DEPLOYMENT_EXIST
	}

	tx := db.Begin() // 开始事务
	// 更新发布项目配置表
	var maps = make(map[string]interface{})
	maps["deploy_name"] = data.DeployName              // 发布项目名称
	maps["git_url_http"] = data.GitUrlHttp             // githttp链接
	maps["git_url_ssh"] = data.GitUrlSsh               // gitssh链接
	maps["git_branch"] = data.GitBranch                // git分支
	maps["git_user"] = data.GitUser                    // git账号
	maps["git_key"] = data.GitKey                      // git秘钥
	maps["deploy_server_path"] = data.DeployServerPath // 服务器发布目标目录
	// 密码
	db.Select("git_passwd").Where("id = ?", id).First(&oldPwd) // 查询数据库中的密码
	// 如果传递过来的密码和数据库中的密码不一致则改变
	if data.GitPasswd != oldPwd.GitPasswd {
		data.GitPasswd, _ = aespwd.EnPwdCode(data.GitPasswd)
	}
	maps["git_passwd"] = data.GitPasswd // git密码
	err := tx.Model(&deployment).Where("id = ?", id).Updates(maps).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}

	// 更新发布项目和服务器对应表
	// 先删除之前项目对应的服务器、用户数据，然后在添加
	// 删除发布项目和发布用户对应表中的数据
	err = tx.Where("deployment_id = ?", id).Delete(&deploymenttouserrole).Error
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}
	// 删除发布项目和服务器对应表中的数据
	err = tx.Where("deployment_id = ?", id).Delete(&deploymenttoserver).Error
	if err != nil {
		tx.Rollback() // 遇到错误时回滚事务
		return errmsg.ERROR
	}
	// 发布项目和服务器对应表
	// 一个发布项目可以有多台服务器
	for _, idx := range data.ServerID {
		var deploymenttoserver DeploymentToServer
		deploymenttoserver.ServerID = idx    // 服务器ID
		deploymenttoserver.DeploymentID = id // 发布项目ID
		//deploymenttoserver.DeployServerPath = data.DeployServerPath // 服务器发布路径
		err := tx.Create(&deploymenttoserver).Error
		if err != nil {
			fmt.Println("err:", err)
			tx.Rollback() // 遇到错误时回滚事务
			return errmsg.ERROR_DEPLOYMENT_UPDATE_FAIL
		}
	}
	// 发布项目和发布用户对应表
	// 一个项目可以有多个发布人
	for _, idx := range data.UserID {
		var deploymenttouserrole DeploymentToUserRole
		deploymenttouserrole.UserID = idx      // 发布用户ID
		deploymenttouserrole.DeploymentID = id // 发布项目配置ID
		err := tx.Create(&deploymenttouserrole).Error
		if err != nil {
			tx.Rollback() // 遇到错误时回滚事务
			return errmsg.ERROR_DEPLOYMENT_UPDATE_FAIL
		}
	}
	tx.Commit() // 提交事务
	return errmsg.SUCCESS
}

// 发布项目配置列表、搜索
func FindDeploymentList(pageSize, page int, deployName string) (int64, []DeploymentSentData, int64) {
	var deploymentSentData DeploymentSentData
	var deploymentSentDataList []DeploymentSentData
	var deploymentList []Deployment
	var Total int64
	var serverID []uint // 服务器的ID切片
	var userID []uint   // 用户的ID切片

	// 搜索
	if deployName != "" {
		db.Where("deploy_name LIKE ?", deployName+"%").Find(&deploymentList).Count(&Total) // 查找总数
		err := db.Where("deploy_name LIKE ?", deployName+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&deploymentList).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return 0, nil, 0
		}
		// 循环项目配置的切片 找到对应的服务器和用户资源
		for _, deploymen := range deploymentList {
			deploymentSentData.ID = deploymen.ID
			deploymentSentData.Deployment = deploymen
			// 1 找服务器
			// 1.1 找到服务器的ID 是一个切片
			db.Model(&DeploymentToServer{}).Select("server_id").Where("deployment_id = ?", deploymen.ID).Find(&serverID)
			// 1.2 循环服务器的ID 找到对应的服务器
			var serverList []*Server
			for _, indexID := range serverID {
				var server Server
				db.Preload("Engineroom").Where("ID = ?", indexID).Find(&server)
				serverList = append(serverList, &server)
			}
			deploymentSentData.ServerList = serverList

			// 2 找用户
			// 2.1 找到用户的ID 是一个切片
			db.Model(&DeploymentToUserRole{}).Select("user_id").Where("deployment_id = ?", deploymen.ID).Find(&userID)
			// 2.2 循环用户的ID 找到对应的用户
			var userList []*User
			for _, indexID := range userID {
				var user User
				db.Where("ID = ?", indexID).Find(&user)
				userList = append(userList, &user)
			}
			deploymentSentData.UserList = userList

			deploymentSentDataList = append(deploymentSentDataList, deploymentSentData)
		}

		return Total, deploymentSentDataList, errmsg.SUCCESS
	}

	// 查询所有的项目配置 结果是一个切片
	db.Find(&deploymentList).Count(&Total)
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&deploymentList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, nil, 0
	}
	// 循环项目配置的切片 找到对应的服务器和用户资源
	for _, deploymen := range deploymentList {
		deploymentSentData.ID = deploymen.ID
		deploymentSentData.Deployment = deploymen
		// 1 找服务器
		// 1.1 找到服务器的ID 是一个切片
		db.Model(&DeploymentToServer{}).Select("server_id").Where("deployment_id = ?", deploymen.ID).Find(&serverID)
		// 1.2 循环服务器的ID 找到对应的服务器
		var serverList []*Server
		for _, indexID := range serverID {
			var server Server
			db.Preload("Engineroom").Where("ID = ?", indexID).Find(&server)
			serverList = append(serverList, &server)
		}
		deploymentSentData.ServerList = serverList

		// 2 找用户
		// 2.1 找到用户的ID 是一个切片
		db.Model(&DeploymentToUserRole{}).Select("user_id").Where("deployment_id = ?", deploymen.ID).Find(&userID)
		// 2.2 循环用户的ID 找到对应的用户
		var userList []*User
		for _, indexID := range userID {
			var user User
			db.Where("ID = ?", indexID).Find(&user)
			userList = append(userList, &user)
		}
		deploymentSentData.UserList = userList

		deploymentSentDataList = append(deploymentSentDataList, deploymentSentData)
	}
	return Total, deploymentSentDataList, errmsg.SUCCESS
}
