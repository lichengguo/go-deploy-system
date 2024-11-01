package model

import (
	"go-deploy-system-server/utils/errmsg"
	"gorm.io/gorm"
)

// 发布项目日志

// DeploymentLog 日志数据模型
type DeploymentLog struct {
	gorm.Model
	DeploymentID       int    `gorm:"type:int;" json:"deployment_id" validate:"required" label:"发布项目ID"`
	DeploymentName     string `gorm:"type:varchar(100)" json:"deployment_name" validate:"required,max=100" label:"发布项目名称"`
	DeploymentUserName string `gorm:"type:varchar(100)" json:"deployment_user_name" validate:"required,max=100" label:"发布用户"`
	DeploymentFileList string `gorm:"type:text" json:"deployment_file_list" validate:"required" label:"发布文件列表"`
	DeploymentCommit   string `gorm:"type:varchar(100)" json:"deployment_commit" validate:"required,max=100"  label:"备注"`
	DeploymentStatus   int    `gorm:"type:int;default:2;comment:'1:成功 2:失败'" json:"deployment_status" label:"发布的状态"`
	DeploymentFailInfo string `gorm:"type:varchar(255)" json:"deployment_fail_info" validate:"max=255" label:"发布失败说明"`
	GitHead            string `gorm:"type:varchar(100)" json:"git_head" validate:"required,max=100" label:"GitHead指针"`
	UpdatedAt          int
	CreatedAt          int
}

// 发布日志列表、搜索
func FindDeploymentLogList(pageSize, page int, userName, deploymentName, deploymentUserName string) (int64, []DeploymentLog, int64) {
	var (
		depLogList []DeploymentLog
		depLog     DeploymentLog
		total      int64
	)

	// 分页功能
	if page == 0 {
		page = 1 // page 页码
	}
	switch {
	case pageSize > 100:
		pageSize = 100 // pageSize 每页多少条
	case pageSize <= 0:
		pageSize = 10
	}

	// 如果是管理员可以看到所有的发布日志
	// 普通用户只能看到自己的发布日志
	// 权限判断
	code, role := FindUserRole(userName)
	if code != errmsg.SUCCESS {
		return total, depLogList, errmsg.ERROR
	}

	// 管理员
	if role == 1 {
		// 日志搜索 项目名&&发布人为条件搜索
		if deploymentName != "" && deploymentUserName != "" {
			db.Find(&depLog).
				Where("deployment_name LIKE ? AND deployment_user_name LIKE ?", deploymentName+"%", deploymentUserName+"%").
				Count(&total)
			err := db.Order("ID desc").
				Where("deployment_name LIKE ? AND deployment_user_name LIKE ?", deploymentName+"%", deploymentUserName+"%").
				Limit(pageSize).Offset((page - 1) * pageSize).Find(&depLogList).Error
			if err != nil {
				return 0, nil, errmsg.ERROR
			}
			return total, depLogList, errmsg.SUCCESS
		}

		// 日志搜索 项目名为条件
		if deploymentName != "" {
			db.Find(&depLog).
				Where("deployment_name LIKE ?", deploymentName+"%").Count(&total)
			err := db.Order("ID desc").
				Where("deployment_name LIKE ? ", deploymentName+"%").
				Limit(pageSize).Offset((page - 1) * pageSize).Find(&depLogList).Error
			if err != nil {
				return 0, nil, errmsg.ERROR
			}
			return total, depLogList, errmsg.SUCCESS
		}

		// 日志搜索 发布人为条件
		if deploymentUserName != "" {
			db.Find(&depLog).Where("deployment_user_name LIKE ?", deploymentUserName+"%").Count(&total)
			err := db.Order("ID desc").
				Where("deployment_user_name LIKE ?", deploymentUserName+"%").
				Limit(pageSize).Offset((page - 1) * pageSize).Find(&depLogList).Error
			if err != nil {
				return 0, nil, errmsg.ERROR
			}
			return total, depLogList, errmsg.SUCCESS
		}

		// 日志列表
		db.Find(&depLog).Count(&total) // 日志总数
		err := db.Order("ID desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&depLogList).Error
		if err != nil {
			return 0, nil, errmsg.ERROR
		}

		return total, depLogList, errmsg.SUCCESS
	}

	// 普通用户
	// 日志搜索 项目名为条件
	if deploymentName != "" {
		db.Find(&depLog).
			Where("deployment_user_name = ? AND deployment_name LIKE ?", userName, deploymentName+"%").Count(&total)
		err := db.Order("ID desc").
			Where("deployment_user_name = ? AND deployment_name LIKE ?", userName, deploymentName+"%").
			Limit(pageSize).Offset((page - 1) * pageSize).Find(&depLogList).Error
		if err != nil {
			return 0, nil, errmsg.ERROR
		}
		return total, depLogList, errmsg.SUCCESS
	}

	// 日志列表
	db.Where("deployment_user_name = ?", userName).Find(&depLog).Count(&total) // 日志总数
	err := db.Order("ID desc").
		Where("deployment_user_name = ?", userName).
		Limit(pageSize).Offset((page - 1) * pageSize).Find(&depLogList).Error
	if err != nil {
		return 0, nil, errmsg.ERROR
	}

	return total, depLogList, errmsg.SUCCESS
}

// 添加发布日志
func AddDeploymentLog(recData *RecDeploymentLog, userName string, deploymentStatus int, deploymentFailInfo string) int64 {
	var (
		data           DeploymentLog
		deployment     Deployment
		deploymentName string // 发布项目名称
	)

	// 发布项目名称
	err := db.Model(&deployment).Select("deploy_name").Where("id = ?", recData.DeploymentID).Find(&deploymentName).Error
	if err != nil {
		return errmsg.ERROR
	}

	data.DeploymentID = recData.DeploymentID             // 项目ID
	data.DeploymentName = deploymentName                 // 项目名称
	data.DeploymentUserName = userName                   // 发布用户名称
	data.DeploymentFileList = recData.DeploymentFileList // 发布文件列表字符串
	data.DeploymentCommit = recData.DeploymentCommit     // 备注信息
	data.DeploymentStatus = deploymentStatus             // 发布状态
	data.DeploymentFailInfo = deploymentFailInfo         // 发布失败详情
	data.GitHead = recData.GitHead                       // git指针 用于回滚

	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCESS
}
