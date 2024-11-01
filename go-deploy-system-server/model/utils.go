package model

import (
	"errors"
	"fmt"
	"go-deploy-system-server/utils"
	"go-deploy-system-server/utils/aespwd"
	"go-deploy-system-server/utils/errmsg"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

// 锁文件函数
func WriteDelCheckLock(action string, ID int) int64 {
	// write 生成锁文件
	// del   删除锁文件
	// Check 检查锁文件是否存在
	lockFile := fmt.Sprintf("%s/%d", utils.LockDirPath, ID) // 锁文件

	if action == "write" {
		_, err = os.Create(lockFile)
		if err != nil {
			return errmsg.ERROR_LOCKFILE_WRITE_FAIL // 锁文件创建失败
		}
		return errmsg.SUCCESS
	} else if action == "del" {
		err := os.Remove(lockFile)
		if err != nil {
			return errmsg.ERROR_LOCKFILE_DEL_FAIL // 锁文件删除失败
		}
		return errmsg.SUCCESS
	} else if action == "check" {
		_, err = os.Lstat(lockFile)
		if err != nil {
			return errmsg.ERROR_LOCKFILE_NOT_EXIST // 锁文件不存在
		}
		return errmsg.ERROR_LOCKFILE_EXIST // 锁文件存在
	} else {
		return errmsg.ERROR
	}
}

// 检查用户对于该项目是否拥有权限
func CheckUserDeploymentRole(deploymentID int, userName interface{}) int64 {
	var (
		user                 User
		dep                  Deployment
		userID               int // 用户的ID
		total                int64
		deploymentToUserRole DeploymentToUserRole
	)

	// 0. 查询项目ID是否存在
	err := db.Where("id = ?", deploymentID).First(&dep).Error
	if err != nil || dep.ID == 0 {
		return errmsg.ERROR_DEPLOYMENT_NOT_EXIST
	}

	// 1. 根据用户名称找到当前用户的ID
	err = db.Model(&user).Select("id").Where("user_name = ?", userName).First(&userID).Error
	if err != nil {
		return errmsg.ERROR_USER_NOT_EXIST
	}

	// 2. 判断用户对于该项目是否具有权限
	err = db.Model(&deploymentToUserRole).Where("deployment_id = ? AND user_id = ?", deploymentID, userID).Count(&total).Error
	if err == nil && total > 0 {
		// 具有权限
		return errmsg.SUCCESS
	}

	return errmsg.ERROR_DEPLOYMENT_NOT_ROLE
}

// 拉取git代码到本地
func GitCodeLocal(deployID int) (int64, ReleaseSendData) {
	var (
		releaseSendData ReleaseSendData
		deployment      Deployment
	)

	// 1. 获取发布项目的配置
	err := db.Where("ID = ?", deployID).First(&deployment).Error
	if err != nil {
		return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, releaseSendData
	}

	// 2. 判断是否是第一次拉取git代码
	// 第一次拉取代码 git clone; 否则 git pull
	codePath := fmt.Sprintf("%s/%d/", utils.GitCodePath, deployID) // 具体某个项目代码存储目录

	// 判断项目代码目录是否存在
	_, err = os.Stat(codePath)
	if err != nil {
		// 本地具体项目代码目录不存在 git clone
		if deployment.GitUrlHttp != "" {
			// 密码拉取代码
			// git clone https://用户名:密码@git_http的链接 -b 分支 代码存放目录
			gitUser := url.QueryEscape(deployment.GitUser)     // 如果用户名和密码里有@等特殊字符，要转成urlcode码 比如@就转成%40
			pwd, err := aespwd.DePwdCode(deployment.GitPasswd) // 密码解密
			if err != nil {
				return errmsg.ERROR_GIT_PULL_PWD_FAIL, releaseSendData
			}
			gitPasswd := url.QueryEscape(pwd) // 如果用户名和密码里有@等特殊字符，要转成urlcode码 比如@就转成%40
			branch := deployment.GitBranch
			urlList := strings.Split(deployment.GitUrlHttp, "//")
			urlHTTP := fmt.Sprintf("%s//%s:%s@%s", urlList[0], gitUser, gitPasswd, urlList[1])
			cmdStr := fmt.Sprintf("git clone %s -b %s %s", urlHTTP, branch, codePath)
			cmd := exec.Command("/bin/sh", "-c", cmdStr) // 有些系统可能没有/bin/bash
			_, err = cmd.Output()
			if err != nil {
				return errmsg.ERROR_GIT_PULL_PWD_FAIL, releaseSendData
			}

			// 获取本次git的信息
			code, data := GitInfo(deployID, codePath, deployment)
			if code != errmsg.SUCCESS {
				return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, releaseSendData
			}
			return errmsg.SUCCESS, data

		} else if deployment.GitKey != "" {
			// 秘钥拉取代码
			// ssh-agent sh -c 'ssh-add 秘钥路径; git clone git_ssh链接 -b 分支 存放目录'
			urlSSH := fmt.Sprintf("git clone %s -b %s %s", deployment.GitUrlSsh, deployment.GitBranch, codePath)
			key := fmt.Sprintf("%s/%s", utils.ExecBaseDir, deployment.GitKey)
			cmdStr := fmt.Sprintf(`ssh-agent sh -c 'ssh-add %s;%s'`, key, urlSSH)
			cmd := exec.Command("/bin/sh", "-c", cmdStr)
			_, err = cmd.Output()
			if err != nil {
				return errmsg.ERROR_GIT_PULL_KEY_FAIL, releaseSendData
			}

			// 获取本次git的信息
			code, data := GitInfo(deployID, codePath, deployment)
			if code != errmsg.SUCCESS {
				return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, releaseSendData
			}

			return errmsg.SUCCESS, data

		} else {
			// 拉取代码失败；git的SSH链接和HTTP链接同时为空
			return errmsg.ERROR_GIT_CONFIG_FAIL, releaseSendData
		}
	}

	// 本地具体项目代码目录存在 git pull
	if deployment.GitUrlHttp != "" {
		// https链接 git pull origin 分支
		//cmdStr := fmt.Sprintf(`git pull origin %s`, deployment.GitBranch)
		cmdStr := fmt.Sprintf(`cd %s;git pull origin %s`, codePath, deployment.GitBranch)

		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		_, err = cmd.Output()
		if err != nil {
			return errmsg.ERROR_GIT_PULL_PWD_FAIL, releaseSendData
		}

		// 获取本次git的信息
		code, data := GitInfo(deployID, codePath, deployment)
		if code != errmsg.SUCCESS {
			return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, releaseSendData
		}

		return errmsg.SUCCESS, data

	} else if deployment.GitKey != "" {
		// ssh链接 ssh-agent sh -c 'ssh-add 秘钥路径; git pull origin 分支'
		urlStr := fmt.Sprintf("git pull origin %s", deployment.GitBranch)
		cmdStr := fmt.Sprintf(`cd %s;ssh-agent sh -c 'ssh-add %s;%s'`, codePath, utils.ExecBaseDir+"/"+deployment.GitKey, urlStr)
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		_, err = cmd.Output()
		if err != nil {
			return errmsg.ERROR_GIT_PULL_KEY_FAIL, releaseSendData
		}

		// 获取本次git的信息
		code, data := GitInfo(deployID, codePath, deployment)
		if code != errmsg.SUCCESS {
			return errmsg.ERROR_DEPLOYMENT_NOT_EXIST, releaseSendData
		}

		return errmsg.SUCCESS, data

	} else {
		// 拉取代码失败；git的SSH链接和HTTP链接同时为空
		return errmsg.ERROR_GIT_CONFIG_FAIL, releaseSendData
	}
}

// 获取git clone以后的相关信息
func GitInfo(deployID int, codePath string, deployment Deployment) (int64, ReleaseSendData) {
	var (
		server             Server
		releaseSendData    ReleaseSendData
		deploymentToServer DeploymentToServer
		serverID           []int    // 服务器ID列表
		serverIP           []string // 服务器IP列表
	)

	// 1. 获取本次更新的相关信息
	cmsStr := fmt.Sprintf(`cd %s;git log --pretty=format:"%%h@YHD@%%ce@YHD@%%s" -1`, codePath)
	cmd := exec.Command("/bin/sh", "-c", cmsStr)
	res, err := cmd.Output()
	if err != nil {
		return errmsg.ERROR_GIT_CONFIG_FAIL, releaseSendData
	}
	resStrList := strings.Split(string(res), "@YHD@")

	releaseSendData.GitHead = resStrList[0]                  // git指针
	releaseSendData.CommitEmail = resStrList[1]              // 提交者邮箱
	releaseSendData.GitInfo = resStrList[2]                  // git提交注释
	releaseSendData.DeployPath = deployment.DeployServerPath // 服务器发布路径
	if deployment.GitUrlHttp != "" {
		releaseSendData.GitUrl = deployment.GitUrlHttp // git http链接
	} else {
		releaseSendData.GitUrl = deployment.GitUrlSsh // git ssh链接
	}

	// 2. 获取本次更新的文件
	// git log --pretty=oneline --name-status -1 |egrep '^M|^A' | awk '{print $2}'
	// git config  core.quotepath false 解决中文乱码问题
	cmsStr = fmt.Sprintf(`cd %s;git config core.quotepath false;git log --pretty=oneline --name-status -1 |egrep '^M|^A' | awk '{print $2}'`, codePath)
	cmd = exec.Command("/bin/sh", "-c", cmsStr)
	res, err = cmd.Output()
	if err != nil {
		return errmsg.ERROR_GIT_CONFIG_FAIL, releaseSendData
	}
	fileList := strings.Split(string(res), "\n")
	releaseSendData.DeployFileList = fileList[:len(fileList)-1]

	// 3. 获取服务器的IP
	// 3.1 根据项目的ID，找到对应的服务器的ID 是一个列表
	db.Model(&deploymentToServer).Select("server_id").Where("deployment_id = ?", deployID).Find(&serverID)
	// 3.2 根据服务器的ID，找到对应的服务器IP 是一个列表
	db.Model(&server).Select("server_ip").Where("ID in ?", serverID).Find(&serverIP)
	releaseSendData.DeployIP = serverIP

	return errmsg.SUCCESS, releaseSendData
}

// 使用密码或者秘钥测试连接发布服务器
func KeyOrPwdConnectLinuxServer(sshHost, sshUser, sshPassword, sshKey string, sshPort string) (*ssh.Client, error) {
	// 创建ssh登录配置
	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,             // 超时时间
		User:            sshUser,                     // 登录账号
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 连接方式
	}

	// 秘钥登录
	if sshKey != "" {
		// 读取秘钥
		key, err := ioutil.ReadFile(sshKey)
		if err != nil {
			fmt.Println("秘钥读取失败")
			return nil, err
		}
		// 创建秘钥签名
		// 会拿着秘钥去登录服务器
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			fmt.Println("秘钥签名失败")
			return nil, err
		}
		// 配置秘钥登录
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else if sshPassword != "" {
		// 密码登录
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		err := errors.New("秘钥或者密码登录")
		return nil, err
	}

	// dial连接服务器
	addr := fmt.Sprintf("%s:%s", sshHost, sshPort)
	Client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("连接服务器失败")
		return nil, err
	}

	return Client, nil
}
