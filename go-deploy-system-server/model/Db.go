package model

import (
	"fmt"
	"go-deploy-system-server/utils"
	"go-deploy-system-server/utils/errmsg"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 连接数据库

var (
	db  *gorm.DB
	err error
)

func InitDb() {
	// 拼接数据库连接信息
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.DbUser,
		utils.DbPassWord,
		utils.DbHost,
		utils.DbPort,
		utils.DbName,
	)

	// 连接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false, // 启用外键约束
		SkipDefaultTransaction:                   true,  // 禁用默认事务（提高运行速度）
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 启用该选项时 `User`的表名是`user`
		},
	})
	if err != nil {
		fmt.Println("连接数据库失败,请检查参数")
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(8)                   // 设置连接池中的最大闲置连接数
	sqlDB.SetMaxOpenConns(64)                  // 设置数据库的最大连接数量
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 设置连接的最大可复用时间

	// 迁移表
	_ = db.AutoMigrate(
		&Department{},           // 部门表
		&User{},                 // 用户表
		&Engineroom{},           // 机房表
		&Server{},               // 服务器表
		&Deployment{},           // 发布项目配置表
		&DeploymentToServer{},   // 发布项目和服务器对应表
		&DeploymentToUserRole{}, // 发布项目和发布用户对应表
		&DeploymentLog{},        // 发布项目日志表
	)

	// 创建初始化默认账户
	code := initData()
	if code != errmsg.SUCCESS {
		fmt.Println("程序初始化默认用户失败", errmsg.GetErrMsg(code))
		os.Exit(1)
	}
}

// initData 初始化部门、账号、密码
func initData() int64 {
	var (
		department Department // 部门
		user       User       // 用户
		count      int64
	)

	// 创建部门
	db.Model(&department).Count(&count)
	if count == 0 {
		// 部门表没有任何记录；创建部门
		department.DepartmentName = "管理部"
		code := AddDepartment(&department)
		if code != errmsg.SUCCESS {
			return errmsg.ERROR_DEPARTMENT_CREATED_FAIL // 创建部门失败
		}

		// 创建登录用户
		user.DepartmentID = 1
		user.UserName = "admin"
		user.Password = "123456"
		user.Role = 1
		user.Status = 1
		code = AddUser(&user)
		if code != errmsg.SUCCESS {
			return errmsg.ERROR_USER_CREATED_FAIL // 创建用户失败
		}
	}

	return errmsg.SUCCESS
}
