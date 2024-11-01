package v1

import (
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 发布日志

// GetDeploymentLogs 日志列表、搜索
func GetDeploymentLogs(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	page, _ := strconv.Atoi(c.Query("page"))

	userName, _ := c.Get("name") // 获取用户名
	userNameStr := userName.(string)

	deploymentName := c.Query("deployment_name")          // 搜索 项目名
	deploymentUserName := c.Query("deployment_user_name") // 搜索 发布人名称

	total, data, code := model.FindDeploymentLogList(pageSize, page, userNameStr, deploymentName, deploymentUserName)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
