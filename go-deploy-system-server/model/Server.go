package model

import (
	"go-deploy-system-server/utils/aespwd"
	"go-deploy-system-server/utils/errmsg"
	"gorm.io/gorm"
)

// 服务器

// 服务器数据模型
type Server struct {
	gorm.Model
	EngineroomID int        `gorm:"type:int" json:"engineroom_id" validate:"required" label:"机房ID"`
	Engineroom   Engineroom `gorm:"foreignkey:EngineroomID" validate:"-"`
	ServerName   string     `gorm:"type:varchar(100)" json:"server_name" validate:"required,max=100" label:"服务器名称"`
	ServerIP     string     `gorm:"type:varchar(100)" json:"server_ip" validate:"required,max=100" label:"服务器登录IP"`
	ServerPort   string     `gorm:"type:varchar(100)" json:"server_port" validate:"required,max=100" label:"服务器登录端口"`
	ServerUser   string     `gorm:"type:varchar(100)" json:"server_user" validate:"required,max=100" label:"服务器登录账号"`
	ServerPwd    string     `gorm:"type:varchar(100)" json:"server_pwd"  label:"服务器登录密码"`
	ServerKey    string     `gorm:"type:varchar(100)" json:"server_key"  label:"服务器登录秘钥存储目录"`
	ServerStatus int        `gorm:"type:int;defalut:2;comment:'1:可使用 2:冻结'" json:"server_status" label:"服务器状态"`
	UpdatedAt    int        // 更新时间
	CreatedAt    int        // 创建时间
}

// 添加服务器
func AddServer(data *Server) int64 {
	var (
		server Server
		total  int64
	)

	// 添加服务器之前，查询服务器是否已经添加
	// 同一个机房里 服务器名称 服务器IP 服务器用户不能同时相同
	db.Model(&server).Where("(server_name= ? AND server_ip = ? AND server_user = ?) AND engineroom_id = ?",
		data.ServerName, data.ServerIP, data.ServerUser, data.EngineroomID).Count(&total)
	if total > 0 {
		return errmsg.ERROR_SERVER_EXIST
	}

	// 存储密码要加密，调用密码去连接服务器的时候再解密
	if data.ServerPwd != "" {
		data.ServerPwd, _ = aespwd.EnPwdCode(data.ServerPwd)
	}

	// 添加服务器
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR_SERVER_CREATED_FAIL
	}

	return errmsg.SUCCESS
}

// 删除服务器
func DelServer(id int) int64 {
	var server Server

	// 如果使用了 gorm.model，会在数据库软删除
	err := db.Where("id = ?", id).Delete(&server).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCESS
}

// 修改服务器
func ModServer(id int, data *Server) int64 {
	var (
		server      Server
		oldPwd      Server
		serverTotal int64
		serverMaps  = make(map[string]interface{})
	)

	// 修改前，检查是否存在同样的服务器
	// 同一个机房里 服务器名称 服务器IP 服务器用户不能同时相同
	db.Model(&server).Where("id <> ? AND (server_name= ? AND server_ip = ? AND server_user = ?) AND engineroom_id = ?",
		id, data.ServerName, data.ServerIP, data.ServerUser, data.EngineroomID).Count(&serverTotal)
	if serverTotal > 0 {
		return errmsg.ERROR_SERVER_EXIST
	}

	// 使用Updates方法更新数据
	// 当通过 struct 更新时，GORM 只会更新非零字段。
	// 如果您想确保指定字段被更新，你应该使用 Select 更新选定字段，或使用 map 来完成更新操作
	serverMaps["engineroom_id"] = data.EngineroomID
	serverMaps["server_name"] = data.ServerName
	serverMaps["server_ip"] = data.ServerIP
	serverMaps["server_port"] = data.ServerPort
	serverMaps["server_user"] = data.ServerUser
	serverMaps["server_key"] = data.ServerKey
	serverMaps["server_status"] = data.ServerStatus
	// 如果传递过来的密码不为空，且和数据库中的密码不一致则改变
	db.Select("server_pwd").Where("id = ?", id).First(&oldPwd) // 查询数据库中的密码
	if data.ServerPwd != "" {
		if data.ServerPwd != oldPwd.ServerPwd {
			data.ServerPwd, _ = aespwd.EnPwdCode(data.ServerPwd)
		}
	}
	serverMaps["server_pwd"] = data.ServerPwd

	err := db.Model(&server).Where("id = ?", id).Updates(serverMaps).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCESS
}

// 服务器列表、搜索
func FindServerList(pageSize, page int, serverName string) (int64, []Server, int64) {
	var (
		serverList []Server
		total      int64
	)

	// 搜索
	if serverName != "" {
		db.Preload("Engineroom").Where("server_name LIKE ?", serverName+"%").Find(&serverList).Count(&total)
		err := db.Preload("Engineroom").Where("server_name LIKE ?", serverName+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&serverList).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return 0, nil, errmsg.ERROR
		}
		return total, serverList, errmsg.SUCCESS
	}

	// 多表查询 预加载
	db.Preload("Engineroom").Find(&serverList).Count(&total)
	err := db.Preload("Engineroom").Limit(pageSize).Offset((page - 1) * pageSize).Find(&serverList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, nil, errmsg.ERROR
	}

	return total, serverList, errmsg.SUCCESS
}

// 服务器连接测试
func ConnectServer(id int) int64 {
	var (
		server Server
	)

	err := db.Where("id = ?", id).First(&server).Error
	if err != nil {
		return errmsg.ERROR
	}

	if server.ServerStatus != 1 {
		// 服务器冻结
		return errmsg.ERROR
	}

	if server.ServerKey != "" {
		// 秘钥
		conn, err := KeyOrPwdConnectLinuxServer(server.ServerIP, server.ServerUser, "", server.ServerKey, server.ServerPort)
		if err != nil {
			return errmsg.ERROR
		}
		_ = conn.Close()

		return errmsg.SUCCESS
	}

	if server.ServerPwd != "" {
		// 密码
		password, _ := aespwd.DePwdCode(server.ServerPwd) // 密码解密

		conn, err := KeyOrPwdConnectLinuxServer(server.ServerIP, server.ServerUser, password, "", server.ServerPort)
		if err != nil {
			return errmsg.ERROR
		}
		_ = conn.Close()

		return errmsg.SUCCESS
	}

	return errmsg.ERROR
}
