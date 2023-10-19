package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"gogin/pkg/app"
	"gogin/pkg/mistakeMsg"
	"gogin/service/auth_service"
	"net/http"
)

//生成HS256token

type Auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuthUsingHS256(c *gin.Context) {
	var (
		appG = app.Response{C: c}
	)
	username := c.PostForm("username")
	password := c.PostForm("password")
	valid := validation.Validation{}
	a := Auth{
		Username: username,
		Password: password,
	}
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MakeErrors(valid.Errors)
		appG.ResponseJson(http.StatusBadRequest, mistakeMsg.INVALID_PARAMS, nil)
		return
	}
	auth := auth_service.Auth{
		Username: a.Username,
		Password: a.Password,
	}
	exists, err := auth.Exist()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if !exists {
		appG.ResponseJson(http.StatusOK, mistakeMsg.ERROR_AUTH, nil)
		return
	}
	token, err := auth.GenerateHS256Token()
	if err != nil {
		appG.ResponseJson(http.StatusInternalServerError, mistakeMsg.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.ResponseJson(http.StatusOK, mistakeMsg.SUCCESS, map[string]interface{}{
		"token": token,
	})
}
