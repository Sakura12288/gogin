package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"gogin/pkg/app"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"gogin/pkg/qrcode"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/service/article_service"
	"gogin/service/tag_service"
	"net/http"
)

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;min(1)"`
	Title         string `form:"title" valid:"MaxSize(100)"`
	Desc          string `form:"desc" valid:"MaxSize(255)"`
	Content       string `form:"content" valid:"MaxSize(65535)"`
	CoverImageUrl string `form:"cover_image_url" valid:"MaxSize(255)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary	编辑文章内容
// @Produce	json
// @Param		id			path		int		true	"文章ID"
// @Param		tag_id		body		int		false	"tag_ID"
// @Param		desc		body		string	false	"简述"
// @Param		content		body		string	false	"文章内容"
// @Param		modified_by	body		string	true	"调整人"
// @Success	200			{object}	string	"成功"
// @Failure	400			{object}	string	"参数错误"
// @Failure	10002		{object}	string	"tag不存在"
// @Failure	10003		{object}	string	"文章不存在"
// @Router		/api/v1/articles/{id} [put]
func EditArticles(c *gin.Context) { //注意修改可以为部分修改，故部分内容可以为空
	var (
		appG = app.Response{C: c}
		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if httpCode != http.StatusOK || errCode != mistakeMsg.SUCCESS {
		appG.ResponseJson(httpCode, errCode, nil)
	}
	article := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exists, err := article.ExistByID()
	if err != nil {
		errCode = mistakeMsg.ERROR_CHECK_EXIST_ARTICLE_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	if !exists {
		errCode = mistakeMsg.ERROR_NOT_EXIST_ARTICLE
		appG.ResponseJson(http.StatusOK, errCode, nil)
		return
	}
	tag := tag_service.Tag{ID: form.TagID}
	exists, err = tag.ExistByID()
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
	err = article.Edit()
	if err != nil {
		errCode = mistakeMsg.ERROR_EDIT_ARTICLE_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, errCode, nil)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

func AddArticles(c *gin.Context) { //感觉是不是可以把valid加在字段上
	var (
		appG = app.Response{C: c}
		form = AddArticleForm{}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if httpCode != http.StatusOK || errCode != mistakeMsg.SUCCESS {
		appG.ResponseJson(httpCode, errCode, nil)
	}
	article := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		CreatedBy:     form.CreatedBy,
		State:         form.State,
	}
	tag := tag_service.Tag{ID: form.TagID}
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
	err = article.Add()
	if err != nil {
		errCode = mistakeMsg.ERROR_ADD_ARTICLE_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, errCode, nil)
}
func GetArticles(c *gin.Context) { //通过标签id获取大量文章
	appG := app.Response{C: c}

	valid := validation.Validation{}
	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	}
	var tagId int = -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
	}
	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusBadRequest, mistakeMsg.INVALID_PARAMS, nil)
		return
	}
	article := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	articles, err := article.GetAll()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	total, err := article.Count()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, data)
}
func GetArticle(c *gin.Context) {
	appG := app.Response{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Required(id, "id").Message("id 不能为空")
	valid.Min(id, 1, "id").Message("id 为正整数")
	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusBadRequest, mistakeMsg.INVALID_PARAMS, nil)
		return
	}
	article := article_service.Article{ID: id}
	exists, err := article.ExistByID()
	if err != nil {
		errCode := mistakeMsg.ERROR_CHECK_EXIST_ARTICLE_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	if !exists {
		errCode := mistakeMsg.ERROR_NOT_EXIST_ARTICLE
		appG.ResponseJson(http.StatusOK, errCode, nil)
		return
	}
	var data interface{}
	data, err = article.GetArticle()
	if err != nil {
		errCode := mistakeMsg.ERROR_GET_ARTICLE_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, data)
}

func DeleteArticles(c *gin.Context) {
	appG := app.Response{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")
	code := mistakeMsg.INVALID_PARAMS
	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusOK, code, nil)
		return
	}
	article := article_service.Article{ID: id}
	exists, err := article.ExistByID()
	if err != nil {
		errCode := mistakeMsg.ERROR_CHECK_EXIST_ARTICLE_FAIL
		appG.ResponseJson(http.StatusInternalServerError, errCode, nil)
		return
	}
	if !exists {
		errCode := mistakeMsg.ERROR_NOT_EXIST_ARTICLE
		appG.ResponseJson(http.StatusOK, errCode, nil)
		return
	}
	appG.ResponseJson(http.StatusOK, code, nil)
}

const (
	PosterUrl = "https://pan.baidu.com/s/10N0SsXLrfrmSYhfA1yACRw"
)

func GenerateArticlePoster(c *gin.Context) {
	appG := app.Response{C: c}

	article := &article_service.Article{}
	q := qrcode.NewQrCode(PosterUrl, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, q)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Point{
			X: 125,
			Y: 298,
		})
	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.ResponseJson(http.StatusOK, mistakeMsg.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		logging.Info(err)
		return
	}
	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, map[string]string{
		"poster_save_url": qrcode.GetQrCodeFullUrl(posterName),
		"poster_url":      filePath + posterName,
	})
}
