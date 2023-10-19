package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"gogin/pkg/app"
	"gogin/pkg/export"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/service/tag_service"
	"net/http"
)

func GetTags(c *gin.Context) { //似乎没有检查state
	appG := app.Response{C: c}

	valid := validation.Validation{}
	name := c.PostForm("name")
	valid.Required(name, "name").Message("标签名字不能为空")
	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能为0或者1")
	}

	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusBadRequest, mistakeMsg.INVALID_PARAMS, nil)
		return
	}

	tag := tag_service.Tag{Name: name, State: state, PageNum: util.GetPage(c), PageSize: setting.AppSetting.PageSize}
	data := map[string]interface{}{}
	tags, err := tag.GetAll()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	total, err := tag.Count()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_COUNT_TAG_FAIL, nil)
		return
	}
	data["list"] = tags
	data["total"] = total
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, data)
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary	新增文章标签
// @Produce	json
// @Param		name		query		string	true	"Name"
// @Param		state		query		int		false	"State"
// @Param		created_by	query		string	false	"CreatedBy"
// @Success	200			{string}	json	"{"code":200,"data":{},"msg":"ok"}"
// @Router		/api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		appG = app.Response{C: c}
		form = AddTagForm{}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if httpCode != http.StatusOK || errCode != mistakeMsg.SUCCESS {
		appG.ResponseJson(httpCode, errCode, nil)
	}
	tag := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	err := tag.Add()
	if err != nil {
		errCode = mistakeMsg.ERROR_ADD_TAG_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, errCode, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;min(1)"`
	Name       string `form:"name" valid:"MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary	修改文章标签
// @Produce	json
// @Param		id			path		int		true	"TagId"
// @Param		state		query		int		false	"State"
// @Param		modified_by	query		string	true	"ModifiedBy"
// @Success	200			{string}	json	"{"code":200,"data":{},"msg":"ok"}"
// @Failure	400			{string}	json	"{"code":400,"data":{},"msg":"参数错误"}"
// @Router		/api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		appG = app.Response{C: c}
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if httpCode != http.StatusOK || errCode != mistakeMsg.SUCCESS {
		appG.ResponseJson(httpCode, errCode, nil)
	}
	tag := tag_service.Tag{
		ID:   form.ID,
		Name: form.Name,

		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	exists, err := tag.ExistByID()
	if err != nil {
		errCode = mistakeMsg.ERROR_CHECK_EXIST_TAG_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	if !exists {
		errCode = mistakeMsg.ERROR_NOT_EXIST_TAG
		appG.ResponseJson(http.StatusOK, errCode, nil)
		return
	}
	err = tag.Edit()
	if err != nil {
		errCode = mistakeMsg.ERROR_EDIT_TAG_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, errCode, nil)

}
func DeleteTag(c *gin.Context) {
	appG := app.Response{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id 至少为1")

	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusBadRequest, mistakeMsg.INVALID_PARAMS, nil)
		return
	}
	tag := tag_service.Tag{ID: id}
	exists, err := tag.ExistByID()
	if err != nil {
		errCode := mistakeMsg.ERROR_CHECK_EXIST_TAG_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	if !exists {
		errCode := mistakeMsg.ERROR_NOT_EXIST_TAG
		appG.ResponseJson(http.StatusOK, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, nil)
}

func ExportTags(c *gin.Context) {
	appG := app.Response{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	valid := validation.Validation{}
	valid.Required(name, "name").Message("标签名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名字最长为100字节")
	valid.Range(state, 0, 1, "state").Message("状态只能为0或者1")
	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusBadRequest, mistakeMsg.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}
	filename, err := tagService.Export()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

func ImportTags(c *gin.Context) {
	appG := app.Response{C: c}

	file, _, err := c.Request.FormFile("tag_file")
	if err != nil {
		logging.Info(err)
		appG.ResponseJson(http.StatusOK, mistakeMsg.ERROR, nil)
		return
	}
	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, nil)
}
