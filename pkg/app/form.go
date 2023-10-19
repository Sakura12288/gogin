package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"net/http"
)

//检查数据格式等是否符合要求

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		logging.Info(err)
		return http.StatusBadRequest, mistakeMsg.INVALID_PARAMS
	}
	valid := validation.Validation{}
	ok, err := valid.Valid(form)
	if err != nil {
		logging.Info(err)
		return http.StatusInternalServerError, mistakeMsg.ERROR
	}
	if !ok {
		return http.StatusBadRequest, mistakeMsg.INVALID_PARAMS
	}
	return http.StatusOK, mistakeMsg.SUCCESS
}
