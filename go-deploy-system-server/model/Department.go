package model

import (
	"go-deploy-system-server/utils/errmsg"

	"gorm.io/gorm"
)

// 部门管理数据模型

type Department struct {
	gorm.Model
	DepartmentName string `gorm:"type:varchar(20);not null" json:"department_name" validate:"required,max=20" label:"部门名称"`
	UpdatedAt      int
	CreatedAt      int
}

// 添加部门
func AddDepartment(data *Department) int64 {
	var (
		dep   Department
		total int64
	)

	// 添加部门之前，先检查部门名称是否存在
	db.Model(&dep).Where("department_name = ?", data.DepartmentName).Count(&total)
	if total > 0 {
		return errmsg.ERROR_DEPARTMENT_EXIST
	}

	// 添加部门
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR_DEPARTMENT_CREATED_FAIL
	}

	return errmsg.SUCCESS
}

// 删除部门
func DelDepartment(id int) int64 {
	var (
		dep   Department
		user  User
		total int64
	)

	// 如果部门下还存在用户；部门不能删除
	err := db.Model(&user).Where("department_id = ?", id).Count(&total).Error
	if err != nil || total > 0 {
		return errmsg.ERROR_DEPARTMENT_NOT_DEL // 该部门下存在用户，无法删除该部门
	}

	// 如果使用了gorm.Model，会在数据库软删除
	err = db.Where("id = ?", id).Delete(&dep).Error
	if err != nil {
		return errmsg.ERROR_DEPARTMENT_DEL_FAIL
	}

	return errmsg.SUCCESS
}

// 修改部门
func ModDepartment(id int, data *Department) int64 {
	var (
		dep            Department
		total          int64
		departmentMaps = make(map[string]interface{})
	)

	// 检查要修改的部门名称是否在数据库中已经存在
	db.Model(&dep).Where("id <> ? AND department_name = ?", id, data.DepartmentName).Count(&total)
	if total > 0 {
		return errmsg.ERROR_DEPARTMENT_EXIST
	}

	// 使用Updates方法更新数据
	// 当通过 struct 更新时，GORM 只会更新非零字段。
	// 如果您想确保指定字段被更新，你应该使用 Select 更新选定字段，或使用 map 来完成更新操作
	departmentMaps["department_name"] = data.DepartmentName
	err := db.Model(&dep).Where("id = ?", id).Updates(departmentMaps).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errmsg.ERROR_DEPARTMENT_MODIFY_FAIL
	}

	return errmsg.SUCCESS
}

// 部门列表、搜索
func FindDepartmentList(pageSize, page int, departmentName string) (int64, []Department, int64) {
	var (
		dep     Department
		depList []Department
		total   int64
	)

	// 搜索
	if departmentName != "" {
		// 查询部门总数
		db.Model(&dep).Where("department_name LIKE ?", departmentName+"%").Count(&total)
		// 查询部门列表数据
		err := db.Where("department_name LIKE ?", departmentName+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&depList).Error
		if err != nil {
			return 0, nil, errmsg.ERROR
		}
		return total, depList, errmsg.SUCCESS
	}

	db.Model(&dep).Count(&total) // 查询部门总数


	// 查询全部部门
	if pageSize == 0 || page == 0 {
		err := db.Find(&depList).Error
		if err != nil {
			return 0, nil, errmsg.ERROR
		}
		return total, depList, errmsg.SUCCESS
	}

	// 分页查询部门
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&depList).Error
	if err != nil {
		return 0, nil, errmsg.ERROR
	}

	return total, depList, errmsg.SUCCESS
}
