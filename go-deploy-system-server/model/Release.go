package model

import (
	"fmt"
	"go-deploy-system-server/utils"
	"go-deploy-system-server/utils/aespwd"
	"go-deploy-system-server/utils/errmsg"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// 发布代码

// 拉取即将发布的项目到发布机 返回给前端的数据
type ReleaseSendData struct {
	DeployFileList []string // 本次git拉取的文件信息
	GitUrl         string   // git链接
	GitHead        string   // git指针
	CommitEmail    string   // 提交者邮箱
	GitInfo        string   // git提交注释
	DeployIP       []string // 服务器IP
	DeployPath     string   // 服务器发布路径
}

// 获取当前登录用户的所有项目 返回给前端的数据
type resDeploymentData struct {
	ID         uint   `json:"id"`
	DeployName string `json:"deploy_name"` // 发布项目名称
}

// 发布和回滚项目代码，用来接收客户端传递数据
type RecDeploymentLog struct {
	DeploymentID       int    `json:"deployment_id" validate:"required" label:"发布项目ID"`
	DeploymentFileList string `json:"deployment_file_list" validate:"required" label:"发布文件列表"`
	DeploymentCommit   string `json:"deployment_commit" validate:"required,max=100"  label:"备注"`
	GitHead            string `json:"git_head" validate:"required,max=100" label:"GitHead指针"`
	Status             int    // 发布状态
}

// GetUserAllDeployment 获取当前登录用户的所有项目
func GetUserAllDeployment(userName string) (int64, []resDeploymentData) {
	var (
		user                 User
		deploymentToUserRole DeploymentToUserRole
		deployment           Deployment
		userID               int                 // 用户ID
		deploymentID         []int               // 项目ID列表
		resData              []resDeploymentData // 返回给前端的数据列表
	)

	// 1. 根据用户名称找到当前用户的ID
	err = db.Model(&user).Select("ID").Where("user_name = ?", userName).Find(&userID).Error
	if err != nil {
		return errmsg.ERROR_USER_NOT_EXIST, nil
	}

	// 2. 找到用户ID对应的项目ID 是一个列表
	err = db.Model(&deploymentToUserRole).Select("deployment_id").Where("user_id = ?", userID).Find(&deploymentID).Error
	if err != nil {
		return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, nil
	}

	// 3. 根据项目ID找到项目信息
	err = db.Model(&deployment).Select("ID, deploy_name").Where("ID IN ?", deploymentID).Find(&resData).Error
	if err != nil {
		return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, nil
	}

	return errmsg.SUCCESS, resData
}

// GitPullDeployment 拉取即将发布的项目代码到本地
func GitPullDeployment(deploymentID int, userName string) (int64, ReleaseSendData) {
	var (
		resData ReleaseSendData
	)

	// 锁文件
	code := WriteDelCheckLock("check", deploymentID) // 检查锁文件是否存在，用来判断是否有其他人正在发布该项目
	if code == errmsg.ERROR_LOCKFILE_EXIST {
		return errmsg.ERROR_LOCKFILE_OTHER_INFO, resData // 有其他人正在发布该项目
	}
	_ = WriteDelCheckLock("write", deploymentID) // 生成锁文件
	defer WriteDelCheckLock("del", deploymentID) // 延迟删除锁文件

	// 权限判断
	code = CheckUserDeploymentRole(deploymentID, userName)
	if code == errmsg.SUCCESS {
		// 具有权限 拉取代码到本地发布机目录
		code, resData := GitCodeLocal(deploymentID)
		return code, resData
	}

	return errmsg.ERROR_DEPLOYMENT_NOT_ROLE, resData
}

// 发布代码
func ReleaseToServer(data *RecDeploymentLog, userName string) int64 {
	var (
		deployment         Deployment
		deploymentToServer DeploymentToServer
		codePath           string   // 项目代码本地存储目录
		serverReleasePath  string   // 目标服务器的代码发布目录
		serverID           []int    // 服务器ID列表
		serverList         []Server // 服务器对象
		wg                 sync.WaitGroup
		codeStatusChan     = make(chan int64, 256) // goroutine里的发布状态
	)

	codePath = fmt.Sprintf("%s/%d/", utils.GitCodePath, data.DeploymentID)

	// 0. 检查锁文件是否存在，用来判断是否有其他人正在发布该项目
	code := WriteDelCheckLock("check", data.DeploymentID)
	if code == errmsg.ERROR_LOCKFILE_EXIST {
		AddDeploymentLog(data, userName, 2, errmsg.GetErrMsg(errmsg.ERROR_LOCKFILE_OTHER_INFO))
		return errmsg.ERROR_LOCKFILE_OTHER_INFO // 有其他人正在发布该项目
	}
	_ = WriteDelCheckLock("write", data.DeploymentID) // 生成锁文件
	defer WriteDelCheckLock("del", data.DeploymentID) // 删除锁文件

	// 1. 权限判断
	code = CheckUserDeploymentRole(data.DeploymentID, userName)
	if code != errmsg.SUCCESS {
		AddDeploymentLog(data, userName, 2, errmsg.GetErrMsg(code))
		return code // 没有该项目权限
	}

	// 2. 找到目标服务器的代码发布目录
	err = db.Model(&deployment).Select("deploy_server_path").Where("id = ?", data.DeploymentID).First(&serverReleasePath).Error
	if err != nil {
		AddDeploymentLog(data, userName, 2, errmsg.GetErrMsg(errmsg.ERROR_SERVER_DIR_NOT_FOUND))
		return errmsg.ERROR_SERVER_DIR_NOT_FOUND // 目标服务器代码发布目录未找到
	}

	// 3. 根据项目ID 找到对应的服务器ID 是一个列表
	err := db.Model(&deploymentToServer).Select("server_id").Where("deployment_id = ?", data.DeploymentID).Find(&serverID).Error
	if err != nil {
		AddDeploymentLog(data, userName, 2, errmsg.GetErrMsg(errmsg.ERROR_SERVER_NOT_EXIST))
		return errmsg.ERROR_SERVER_NOT_EXIST // 目标服务器未找到
	}

	// 4 发布代码
	// 4.1 找到对应的服务器对象信息 是一个列表
	err = db.Where("id IN ? AND server_status = 1", serverID).Find(&serverList).Error
	if err != nil {
		AddDeploymentLog(data, userName, 2, errmsg.GetErrMsg(errmsg.ERROR_SERVER_NOT_EXIST))
		return errmsg.ERROR_SERVER_NOT_EXIST // 服务器信息未找到
	}

	// 4.2 循环服务器对象信息
	for _, server := range serverList {
		wg.Add(1) // goroutine计数器加1

		go func(server Server) {
			defer wg.Done() // goroutine计数器减1

			serverIP := server.ServerIP     // 服务器IP
			port := server.ServerPort       // 服务器端口
			serverUser := server.ServerUser // 服务器账号

			// 同步代码到服务器
			if server.ServerKey != "" {
				// 使用秘钥方式 推送代码到服务器
				code = rsyncFileToServer(data, server.ServerKey, "", serverIP, serverUser, port, codePath, serverReleasePath, userName)
				codeStatusChan <- code
			} else if server.ServerPwd != "" {
				// 使用密码方式 推送代码到服务器
				code = rsyncFileToServer(data, "", server.ServerPwd, serverIP, serverUser, port, codePath, serverReleasePath, userName)
				codeStatusChan <- code
			} else {
				codeStatusChan <- errmsg.ERROR_SERVER_PWD_NULL // 秘钥和密码同时为空
			}
		}(server)
	}

	wg.Wait() // 等待全部goroutine执行完成

	// 读取通道中的状态码
	// 一台服务器只有一个状态码，所以循环服务器数量, 读取通道即可
	for i := 0; i < len(serverList); i++ {
		s := <-codeStatusChan
		if s != errmsg.SUCCESS {
			// 有服务器发布失败
			AddDeploymentLog(data, userName, 2, errmsg.GetErrMsg(s))
			return s
		}
	}

	AddDeploymentLog(data, userName, 1, "")
	return errmsg.SUCCESS
}

// 同步代码文件到远程服务器
func rsyncFileToServer(data *RecDeploymentLog, keyPath, passwd, ip, user, port, codePath, serverReleasePath, userName string) int64 {
	/*
		data:				前端传递过来的数据
		keyPath: 			秘钥路径
		passwd: 			密码
		ip: 				服务器IP
		user: 				服务器账号
		port: 				服务器端口
		codePath: 			本地代码路径
		serverReleasePath: 	远程服务器代码路径
		userName: 			发布人名称
	*/

	var (
		deployFileStr string // 需要发布的文件拼接字符串
	)

	// 获取客户端传递过来的文件列表
	deploymentFile := strings.Split(data.DeploymentFileList, "\n")

	if keyPath != "" {
		// 使用秘钥方式 推送代码到服务器
		key := utils.ExecBaseDir + "/" + keyPath // 秘钥路径

		// 测试服务器是否能连接
		conn, err := KeyOrPwdConnectLinuxServer(ip, user, "", key, port)
		if err != nil {
			return errmsg.ERROR_DEPLOYMENT_CONNECT_SERVER_FAIL // 连接服务器失败
		}
		_ = conn.Close()

		// 如果发布文件中有*号则整个项目发布；否则按照发布文件单个或者多个文件发布
		for _, v := range deploymentFile {
			// 发布文件不能以 / 开头
			if strings.HasPrefix(v, "/") {
				return errmsg.ERROR_DEPLOYMENT_FILE_NOT_GEN
			}

			if v == "*" {
				// 发布整个项目
				cmdStr := fmt.Sprintf(`cd %s;rsync -aR -e "ssh -p %s -i %s -o PubkeyAuthentication=yes -o stricthostkeychecking=no" --delete ./ %s@%s:%s`,
					codePath,
					port,
					key,
					user,
					ip,
					serverReleasePath,
				)
				cmd := exec.Command("/bin/sh", "-c", cmdStr)
				_, err = cmd.Output()
				if err != nil {
					return errmsg.ERROR_DEPLOYMENT_RSYNC_FAIL // 同步代码到服务器失败
				}

				return errmsg.SUCCESS
			}

			// 校验每个文件在本地git目录是否存在
			_, err := os.Stat(fmt.Sprintf("%s/%s", codePath, strings.Trim(v, " ")))
			if err != nil {
				return errmsg.ERROR_DEPLOYMENT_FILE_NOT_FOUND // 发布文件未找到
			}

			// 拼接发布文件
			deployFileStr += v + " "
		}

		// 单个文件发布
		cmdStr := fmt.Sprintf(`cd %s;rsync -aR -e "ssh -p %s -i %s -o PubkeyAuthentication=yes -o stricthostkeychecking=no" --delete %s %s@%s:%s/`,
			codePath,
			port,
			key,
			deployFileStr,
			user,
			ip,
			serverReleasePath,
		)
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		_, err = cmd.Output()
		if err != nil {
			return errmsg.ERROR_DEPLOYMENT_RSYNC_FAIL // 同步代码到服务器失败
		}

		return errmsg.SUCCESS
	}

	if passwd != "" {
		// 使用密码同步文件
		password, _ := aespwd.DePwdCode(passwd) // 密码解密

		// 测试服务器是否能连接
		conn, err := KeyOrPwdConnectLinuxServer(ip, user, password, "", port)
		if err != nil {
			return errmsg.ERROR_DEPLOYMENT_CONNECT_SERVER_FAIL // 连接服务器失败
		}
		_ = conn.Close()

		// 如果发布文件中有*号则整个项目发布；否则按照发布文件单个或者多个文件发布
		for _, v := range deploymentFile {
			// 发布文件不能以 / 开头
			if strings.HasPrefix(v, "/") {
				return errmsg.ERROR_DEPLOYMENT_FILE_NOT_GEN
			}

			if v == "*" {
				// 发布整个项目
				cmdStr := fmt.Sprintf(`cd %s;rsync -aR -e "sshpass -p '%s' ssh -p %s -o PubkeyAuthentication=yes -o stricthostkeychecking=no" --delete ./ %s@%s:%s`,
					codePath,
					password,
					port,
					user,
					ip,
					serverReleasePath,
				)
				cmd := exec.Command("/bin/sh", "-c", cmdStr)
				_, err = cmd.Output()
				if err != nil {
					return errmsg.ERROR_DEPLOYMENT_RSYNC_FAIL // 同步代码到服务器失败
				}

				return errmsg.SUCCESS
			}

			// 校验每个文件在本地git目录是否存在
			_, err := os.Stat(fmt.Sprintf("%s/%s", codePath, strings.Trim(v, " ")))
			if err != nil {
				return errmsg.ERROR_DEPLOYMENT_FILE_NOT_FOUND // 发布文件未找到
			}

			// 拼接发布文件
			deployFileStr += v + " "
		}

		// 发布单个或者多个文件
		cmdStr := fmt.Sprintf(`cd %s;rsync -aR -e "sshpass -p '%s' ssh -p %s -o PubkeyAuthentication=yes -o stricthostkeychecking=no" --delete %s %s@%s:%s/`,
			codePath,
			password,
			port,
			deployFileStr,
			user,
			ip,
			serverReleasePath,
		)
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		_, err = cmd.Output()
		if err != nil {
			return errmsg.ERROR_DEPLOYMENT_RSYNC_FAIL // 同步代码到服务器失败
		}

		return errmsg.SUCCESS
	}

	// 同步文件失败
	return errmsg.ERROR_DEPLOYMENT_RSYNC_FAIL // 同步代码到服务器失败
}

// 回滚代码
func RollBackCode(deploymentLogID int, userName string) int64 {
	var (
		codePath      string        // 项目代码本地存储目录
		deploymentLog DeploymentLog // 日志表对象
		data          RecDeploymentLog
	)

	// 找到日志表中的记录
	err := db.Where("id = ?", deploymentLogID).First(&deploymentLog).Error
	if err != nil || deploymentLog.DeploymentStatus != 1 {
		return errmsg.ERROR_DEPLOYMENT_STATUS_FAIL // 发布状态为失败, 不能回滚
	}

	data = RecDeploymentLog{
		DeploymentID:       deploymentLog.DeploymentID,
		DeploymentFileList: deploymentLog.DeploymentFileList,
		DeploymentCommit:   deploymentLog.DeploymentCommit,
		GitHead:            deploymentLog.GitHead,
	}

	// 检查锁文件是否存在，用来判断是否有其他人正在发布该项目
	code := WriteDelCheckLock("check", data.DeploymentID)
	if code == errmsg.ERROR_LOCKFILE_EXIST {
		return errmsg.ERROR_LOCKFILE_OTHER_INFO // 有其他人正在发布该项目
	}
	_ = WriteDelCheckLock("write", data.DeploymentID) // 生成锁文件
	defer WriteDelCheckLock("del", data.DeploymentID) // 删除锁文件

	// 权限判断
	code = CheckUserDeploymentRole(data.DeploymentID, userName)
	if code != errmsg.SUCCESS {
		return code // 没有该项目权限
	}

	// 本地Git回滚代码
	codePath = fmt.Sprintf("%s/%d/", utils.GitCodePath, deploymentLog.DeploymentID)
	cmdStr := fmt.Sprintf(`cd %s;git reset --hard %s`, codePath, data.GitHead)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	_, err = cmd.Output()
	if err != nil {
		return errmsg.ERROR_GIT_ROLLBACK_FAIL // 代码回滚失败
	}

	// 发布代码
	WriteDelCheckLock("del", data.DeploymentID) // 删除锁文件
	code = ReleaseToServer(&data, userName)
	if code != errmsg.SUCCESS {
		return code
	}

	return errmsg.SUCCESS
}
