package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model
	TagID         int    `json:"tag_id" gorm:"index"`
	Tag           Tag    `json:"tag"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article, err error) {
	if err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func GetArticleTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Article{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

//判断是否存在

//func ExistByNameAt(name string) (bool, error) {
//	var article Article
//	err := db.Select("id").Where("name = ? and deleted_on = ?", name, 0).First(&article).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		return false, err
//	}
//	if article.ID > 0 {
//		return true, nil
//	}
//	return false, nil
//}

func ExistByIdAt(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? and deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}
func AddArticle(data map[string]interface{}) error {
	err := db.Create(&Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		CreatedBy:     data["created_by"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CoverImageUrl: data["cover_image_url"].(string),
		State:         data["state"].(int),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? and deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

func EditArticles(id int, maps interface{}) error {
	err := db.Model(&Article{}).Where("id = ?", id).Update(maps).Error
	if err != nil {
		return err
	}
	return nil
}
func DeleteArticles(id int) error {
	if err := db.Where("id = ?", id).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}

//触发器

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedOn", time.Now().Unix())
}
func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("ModifiedOn", time.Now().Unix())
}

//硬删除

func CleanAllArticles() (bool, error) {
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
