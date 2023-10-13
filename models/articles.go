package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model
	TagID      int    `json:"tag_id" gorm:"index"`
	Tag        Tag    `json:"tag"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Tag) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

//判断是否存在

func ExistByNameAt(name string) bool {
	var article Article
	db.Select("id").Where("name = ?", name).First(&article)
	return article.ID > 0
}

func ExistByIdAt(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)
	return article.ID > 0
}
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		CreatedBy: data["created_by"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		State:     data["state"].(int),
	})
	return true
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).Find(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

func EditArticles(id int, maps interface{}) error {
	db.Model(&Article{}).Where("id = ?", id).Update(maps)
	return nil
}
func DeleteArticles(id int) error {
	db.Where("id = ?", id).Delete(&Article{})
	return nil
}

//触发器

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedOn", time.Now().Unix())
}
func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("ModifiedOn", time.Now().Unix())
}
