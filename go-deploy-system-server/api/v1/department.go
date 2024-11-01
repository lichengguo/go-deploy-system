package v1

import (
	"go-deploy-system-server/model"
	"go-deploy-system-server/utils/errmsg"
	"go-deploy-system-server/utils/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 部门管理接口

// 添加部门
func AddDep(c *gin.Context) {
	var (
		data model.Department
		msg  string
		code int64
	)

	_ = c.ShouldBindJSON(&data) // 绑定客户端传递过来的数据

	// 对客户端提交过来的数据进行验证
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return // 字段校验未通过直接结束函数
	}

	code = model.AddDepartment(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除部门
func DelDep(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DelDepartment(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 修改部门
func EditDep(c *gin.Context) {
	var (
		data model.Department
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

	code = model.ModDepartment(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 部门列表、搜索
func GetDeps(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	page, _ := strconv.Atoi(c.Query("page"))
	departmentName := c.Query("department_name") // 搜索

	total, data, code := model.FindDepartmentList(pageSize, page, departmentName)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
