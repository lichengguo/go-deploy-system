package v1

import (
	"github.com/gin-gonic/gin"
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/validator"
	"net/http"
	"strconv"
)

// 添加发布项目配置
func AddDeployment(c *gin.Context) {
	var (
		data model.DeploymentAcceptData
		msg  string
		code int64
	)

	_ = c.ShouldBindJSON(&data)

	// 对客户端提交过来的数据进行验证
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return // 字段校验未通过直接结束函数
	}

	code = model.AddDeployment(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除发布项目配置
func DelDeployment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DelDeployment(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 修改发布项目配置
func EditDeployment(c *gin.Context) {
	var (
		data model.DeploymentAcceptData
		msg  string
		code int64
	)

	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	// 对客户端提交过来的数据进行验证
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return // 字段校验未通过直接结束函数
	}

	code = model.ModDeployment(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 发布项目配置列表、搜索
func GetDeploymentList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	page, _ := strconv.Atoi(c.Query("page"))
	deploymentName := c.Query("deployment_name") // 搜索

	total, data, code := model.FindDeploymentList(pageSize, page, deploymentName)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
