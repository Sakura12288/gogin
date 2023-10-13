package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"gogin/models"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"gogin/pkg/util"
	"net/http"
)

//生成HS256token

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuthUsingHS256(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	a := auth{Username: username, Password: password}
	valid := validation.Validation{}
	ok, _ := valid.Valid(&a)
	code := mistakeMsg.INVALID_PARAMS
	data := make(map[string]interface{})
	if ok {
		code = mistakeMsg.SUCCESS
		if models.CheckAuth(username, password) {
			token, err := util.GenerateTokenUsingHs256(username, password)
			if err != nil {
				code = mistakeMsg.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
			}
		} else {
			code = mistakeMsg.ERROR_AUTH
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
