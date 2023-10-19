package tag_service

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gogin/models"
	"gogin/pkg/export"
	"gogin/pkg/gredis"
	"gogin/pkg/logging"
	"gogin/service/cache_service"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	State      int
	ModifiedBy string

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistById(t.ID)
}
func (t *Tag) ExistByName() (bool, error) {
	return models.ExistByName(t.Name)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}
func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.GetMaps())
}
func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
		err             error
	)
	cache := cache_service.Tag{
		State:    t.State,
		Name:     t.Name,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			if err := json.Unmarshal(data, &cacheTags); err != nil {
				logging.Info(err)
			}
			return cacheTags, nil
		}
	}
	tags, err = models.GetTags(t.PageNum, t.PageSize, t.GetMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tags, 3600)
	return tags, nil
}
func (t *Tag) GetMaps() map[string]interface{} {
	data := make(map[string]interface{})
	data["deleted_on"] = 0
	if t.Name != "" {
		data["name"] = t.Name
	}
	if t.State >= 0 {
		data["state"] = t.State
	}
	return data
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}
	f := excelize.NewFile()
	_, err = f.NewSheet("标签信息")
	if err != nil {
		return "", fmt.Errorf("建立标签信息表格错误 :%v", err)
	}
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间", "状态"}
	if err := f.SetSheetRow("标签信息", "A1", &titles); err != nil {
		return "", err
	}
	for i, v := range tags {
		cell := string('A') + strconv.Itoa(2+i)
		tag := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
			strconv.Itoa(v.State),
		}
		f.SetSheetRow("标签信息", cell, &tag)
	}
	timeNow := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags" + timeNow + ".xlsx"
	if err := f.SaveAs(export.GetExcelFullSavePath() + filename); err != nil {
		return "", err
	}
	return filename, nil
}

func (t *Tag) Import(r io.Reader) error { //按照导出的顺序导入
	f, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	rows, err := f.GetRows("标签信息")
	if err != nil {
		return nil
	}
	for i, row := range rows {
		if i > 0 {
			var data []string
			for _, v := range row {
				data = append(data, v)
			}
			if err := models.AddTag(data[1], 1, data[2]); err != nil {
				return err
			}
		}
	}
	return nil
}
