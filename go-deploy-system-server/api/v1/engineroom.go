package v1

import (
	"github.com/gin-gonic/gin"
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/validator"
	"net/http"
	"strconv"
)

// 机房管理模块

// 添加机房
func AddEngineroom(c *gin.Context) {
	var (
		data model.Engineroom
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

	code = model.AddEngineroom(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除机房
func DelEngineroom(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DelEngineroom(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 修改机房
func EditEngineroom(c *gin.Context) {
	var (
		data model.Engineroom
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

	code = model.ModEngineroom(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 机房列表、搜索
func GetEnginerooms(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	page, _ := strconv.Atoi(c.Query("page"))
	engineroomName := c.Query("engineroom_name") // 搜索

	total, data, code := model.FindEngineroomList(pageSize, page, engineroomName)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
