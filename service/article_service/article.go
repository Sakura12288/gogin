package article_service

import (
	"encoding/json"
	"gogin/models"
	"gogin/pkg/gredis"
	"gogin/pkg/logging"
	"gogin/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	CreatedBy     string
	ModifiedBy    string
	State         int

	PageNum  int
	PageSize int
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistByIdAt(a.ID)
}

func (a *Article) Add() error {
	data := make(map[string]interface{})
	data["tag_id"] = a.TagID
	data["title"] = a.Title
	data["desc"] = a.Desc
	data["content"] = a.Content
	data["cover_image_url"] = a.CoverImageUrl
	data["created_by"] = a.CreatedBy
	data["state"] = a.State
	return models.AddArticle(data)
}
func (a *Article) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = a.ModifiedBy
	if a.Desc != "" {
		data["desc"] = a.Desc
	}
	if a.Content != "" {
		data["content"] = a.Content
	}
	if a.State >= 0 {
		data["state"] = a.State
	}
	if a.Title != "" {
		data["title"] = a.Title
	}
	if a.CoverImageUrl != "" {
		data["cover_image_url"] = a.CoverImageUrl
	}
	if a.TagID > 0 {
		data["tag_id"] = a.TagID
	}
	return models.EditArticles(a.ID, data)
}

func (a *Article) Delete() error {
	return models.DeleteArticles(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.GetMaps())
}
func (a *Article) GetAll() ([]models.Article, error) {
	var (
		Articles, cacheArticles []models.Article
		err                     error
	)
	cache := cache_service.Article{
		State:    a.State,
		TagID:    a.TagID,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			if err := json.Unmarshal(data, &cacheArticles); err != nil {
				logging.Info(err)
			}
			return cacheArticles, nil
		}
	}
	Articles, err = models.GetArticles(a.PageNum, a.PageSize, a.GetMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, Articles, 3600)
	return Articles, nil
}

func (a *Article) GetArticle() (*models.Article, error) {
	return models.GetArticle(a.ID)
}
func (a *Article) GetMaps() map[string]interface{} {
	data := make(map[string]interface{})
	data["deleted_on"] = 0
	data["tag_id"] = a.TagID
	if a.State >= 0 {
		data["state"] = a.State
	}
	return data
}
