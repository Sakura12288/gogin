package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"gogin/models"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"net/http"
)

//	@Summary	编辑文章内容
//	@Produce	json
//	@Param		id			path		int		true	"文章ID"
//	@Param		tag_id		body		int		false	"tag_ID"
//	@Param		desc		body		string	false	"简述"
//	@Param		content		body		string	false	"文章内容"
//	@Param		modified_by	body		string	true	"调整人"
//	@Success	200			{object}	string	"成功"
//	@Failure	400			{object}	string	"参数错误"
//	@Failure	10002		{object}	string	"tag不存在"
//	@Failure	10003		{object}	string	"文章不存在"
//	@Router		/api/v1/articles/{id} [put]
func EditArticles(c *gin.Context) { //注意修改可以为部分修改，故部分内容可以为空
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	var state int = -1
	data := make(map[string]interface{})
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		data["state"] = state
		valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	}
	valid.Min(id, 1, "id").Message("文章id必须大于0")
	//valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
	//valid.Required(title, "title").Message("title不能为空")
	valid.MaxSize(title, 100, "title").Message("title最长为100字节")
	valid.Required(modifiedBy, "modified_by").Message("ar修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "created_by").Message("ar修改人最长为100字节")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字节")
	//valid.Required(desc, "desc").Message("desc不能为空")
	//valid.Required(content, "title").Message("content不能为空")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字节")

	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistByIdAt(id) {
			if models.ExistById(tagId) {
				data := make(map[string]interface{})
				if title != "" {
					data["title"] = title
				}
				data["tag_id"] = tagId
				if content != "" {
					data["content"] = content
				}
				if desc != "" {
					data["desc"] = desc
				}
				data["modified_by"] = modifiedBy
				models.EditArticles(id, data)
				code = mistakeMsg.SUCCESS
			} else {
				code = mistakeMsg.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = mistakeMsg.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": make(map[string]interface{}),
	})
}
func AddArticles(c *gin.Context) { //感觉是不是可以把valid加在字段上
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
	valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	valid.Required(title, "title").Message("title不能为空")
	valid.MaxSize(title, 100, "title").Message("title最长为100字节")
	valid.Required(createdBy, "created_by").Message("ar创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("ar创建人最长为100字节")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字节")
	valid.Required(desc, "desc").Message("desc不能为空")
	valid.Required(content, "title").Message("content不能为空")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字节")

	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistById(tagId) {
			data := make(map[string]interface{})
			data["title"] = title
			data["tag_id"] = tagId
			data["content"] = content
			data["desc"] = desc
			data["created_by"] = createdBy
			data["state"] = state
			models.AddArticle(data)
			code = mistakeMsg.SUCCESS
		} else {
			code = mistakeMsg.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": make(map[string]interface{}),
	})
}
func GetArticles(c *gin.Context) { //通过标签id获取大量文章
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	}
	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
	}
	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		code = mistakeMsg.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": data,
	})
}
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Required(id, "id").Message("id 不能为空")
	valid.Min(id, 1, "id").Message("id 为正整数")
	code := mistakeMsg.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistByIdAt(id) {
			data = models.GetArticle(id)
			code = mistakeMsg.SUCCESS
		} else {
			code = mistakeMsg.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": data,
	})
}

func DeleteArticles(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")
	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistByIdAt(id) {
			code = mistakeMsg.SUCCESS
			models.DeleteArticles(id)
		} else {
			code = mistakeMsg.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": make(map[string]interface{}),
	})
}
