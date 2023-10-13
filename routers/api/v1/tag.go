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

func GetTags(c *gin.Context) { //似乎没有检查state
	name := c.Query("name")
	t := models.Tag{}
	t.Name = name
	data := map[string]interface{}{}
	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		t.State = state
	}
	code := mistakeMsg.SUCCESS
	data["list"] = models.GetTags(util.GetPage(c), setting.PageSize, t)
	data["total"] = models.GetTagTotal(t)
	c.JSON(code, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": data,
	})
}

//	@Summary	新增文章标签
//	@Produce	json
//	@Param		name		query		string	true	"Name"
//	@Param		state		query		int		false	"State"
//	@Param		created_by	query		int		false	"CreatedBy"
//	@Success	200			{string}	json	"{"code":200,"data":{},"msg":"ok"}"
//	@Router		/api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字节")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字节")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistByName(name) {
			models.AddTag(name, state, createdBy)
			code = mistakeMsg.SUCCESS
		} else {
			code = mistakeMsg.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": make(map[string]string),
	})
}
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Min(id, 1, "id").Message("id 至少为1")
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字节")
	valid.Required(modifiedBy, "modified_by").Message("调整人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("调整人最长为100字节")
	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		code = mistakeMsg.SUCCESS
		if !models.ExistById(id) {
			code = mistakeMsg.ERROR_NOT_EXIST_TAG
		} else {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
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
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id 至少为1")
	code := mistakeMsg.INVALID_PARAMS
	if !valid.HasErrors() {
		code = mistakeMsg.SUCCESS
		if models.ExistById(id) {
			models.DeleteTag(id)
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
