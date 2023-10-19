package app

import (
	"github.com/gin-gonic/gin"
	"gogin/pkg/mistakeMsg"
)

type Response struct {
	C *gin.Context
}

func (r *Response) ResponseJson(httpCode, code int, data interface{}) {
	r.C.JSON(httpCode, gin.H{
		"code": code,
		"msg":  mistakeMsg.GetMsgFlags(code),
		"data": data,
	})
}
