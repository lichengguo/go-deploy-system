package model

import (
	"go-deploy-system-server/utils/errmsg"

	"gorm.io/gorm"
)

// 机房

// 机房数据模型
type Engineroom struct {
	gorm.Model
	EngineroomName string `gorm:"type:varchar(100)" json:"engineroom_name" validate:"required,max=100" label:"机房名称"`
	Contact        string `gorm:"type:varchar(100)" json:"contact" validate:"required,max=100" label:"机房联系人"`
	ContactInfo    string `gorm:"type:varchar(100)" json:"contact_info" validate:"required,max=100" label:"联系方式"`
	Address        string `gorm:"type:varchar(255)" json:"address" validate:"required,max=100" label:"机房地址"`
	UpdatedAt      int
	CreatedAt      int
}

// 添加机房
func AddEngineroom(data *Engineroom) int64 {
	var (
		eng   Engineroom
		total int64
	)

	// 添加机房之前，先检查机房是否存在
	db.Model(&eng).Where("engineroom_name = ?", data.EngineroomName).Count(&total)
	if total > 0 {
		return errmsg.ERROR_DEPARTMENT_EXIST
	}

	// 添加机房
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR_ENGINEROOM_CREATED_FAIL
	}

	return errmsg.SUCCESS
}

// 删除机房
func DelEngineroom(id int) int64 {
	var (
		eng    Engineroom
		server Server
		total  int64
	)

	// 如果机房下还存在服务器；则机房不能删除
	err := db.Model(&server).Where("engineroom_id = ?", id).Count(&total).Error
	if err != nil || total > 0 {
		return errmsg.ERROR_ENGINEROOM_NOT_DEL // 该机房下存在服务器,无法删除
	}

	// 如果使用了gorm.model，会在数据库软删除
	err = db.Where("id = ?", id).Delete(&eng).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCESS
}

// 修改机房
func ModEngineroom(id int, data *Engineroom) int64 {
	var (
		eng     Engineroom
		total   int64
		engMaps = make(map[string]interface{})
	)

	// 检查修改后的机房是否存在多个
	db.Where("id != ? AND engineroom_name = ?", id, data.EngineroomName).First(&eng).Count(&total)
	if total >= 1 {
		return errmsg.ERROR_ENGINEROOM_EXIST
	}

	// 使用Updates方法更新数据
	// 当通过 struct 更新时，GORM 只会更新非零字段。
	// 如果您想确保指定字段被更新，你应该使用 Select 更新选定字段，或使用 map 来完成更新操作
	engMaps["engineroom_name"] = data.EngineroomName
	engMaps["contact"] = data.Contact
	engMaps["contact_info"] = data.ContactInfo
	engMaps["address"] = data.Address
	err := db.Model(&eng).Where("id = ?", id).Updates(engMaps).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errmsg.ERROR_ENGINEROOM_UPDATE_FAIL
	}

	return errmsg.SUCCESS
}

// 机房列表、搜索
func FindEngineroomList(pageSize, page int, engineroomName string) (int64, []Engineroom, int64) {
	var (
		engList []Engineroom
		total   int64
	)

	// 搜索
	if engineroomName != "" {
		db.Where("engineroom_name LIKE ?", engineroomName+"%").Find(&engList).Count(&total)
		err := db.Where("engineroom_name LIKE ?", engineroomName+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&engList).Error
		if err != nil {
			return 0, nil, errmsg.ERROR
		}
		return total, engList, errmsg.SUCCESS
	}

	// 查询全部
	if pageSize == 0 || page == 0 {
		err := db.Find(&engList).Error
		if err != nil {
			return 0, nil, errmsg.ERROR
		}
		return total, engList, errmsg.SUCCESS
	}

	// 分页查询机房
	db.Find(&engList).Count(&total) // 查询机房总数
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&engList).Error
	if err != nil {
		return 0, nil, errmsg.ERROR
	}

	return total, engList, errmsg.SUCCESS
}
