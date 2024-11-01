package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-deploy-system-server/utils"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/md5"
	"net/http"
	"os/exec"
)

// 处理秘钥文件上传

func UpLoad(c *gin.Context) {
	dir := utils.ExecBaseDir // 当前程序执行的目录

	// 获取file文件
	file, err := c.FormFile("keyfile")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"url":     "",
		})
		c.Abort() // 不调用后续的函数
		return
	}

	newFileName := md5.StringToMD5(file.Filename)                        // md5化文件名；避免重名
	fileSavePath := fmt.Sprintf("%s/%s", utils.KeyFilePath, newFileName) // 拼接key存储路径

	// 存储key文件
	err = c.SaveUploadedFile(file, fileSavePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"url":     "",
		})
		c.Abort() // 不调用后续的函数
		return
	}

	// 修改key秘钥权限
	cmdStr := fmt.Sprintf(`chmod 600 %s`, dir+"/"+fileSavePath)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	_, err = cmd.Output()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"url":     "",
		})
		c.Abort() // 不调用后续的函数
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
		"url":     fileSavePath,
	})
}
