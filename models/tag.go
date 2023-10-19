package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	if pageNum > 0 && pageSize > 0 {
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

func GetTagTotal(maps interface{}) (count int, err error) {
	err = db.Where(maps).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, nil
}

//判断是否存在

func ExistByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ? and deleted_on = ?", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistById(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ? and deleted_on = ?", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}
func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func EditTag(id int, maps interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ?", id).Update(maps).Error; err != nil {
		return err
	}
	return nil
}
func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}
	return nil
}

//触发器

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedOn", time.Now().Unix())
}
func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("ModifiedOn", time.Now().Unix())
}

//硬删除

func CleanAllTag() (bool, error) {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
