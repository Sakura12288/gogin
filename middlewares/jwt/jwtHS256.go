package jwt

import (
	"github.com/gin-gonic/gin"
	"gogin/pkg/logging"
	"gogin/pkg/mistakeMsg"
	"gogin/pkg/util"
	"net/http"
	"time"
)

func JwtHS256() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		code := mistakeMsg.SUCCESS
		if token == "" {
			code = mistakeMsg.INVALID_PARAMS
		} else {
			claim, err := util.ParseToken(token)
			if err != nil { //这里逻辑有问题，超时不会再进行下面的步骤
				logging.Info(err)
				code = mistakeMsg.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if claim.ExpiresAt < time.Now().Unix() {
				code = mistakeMsg.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}
		if code != mistakeMsg.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  mistakeMsg.GetMsgFlags(code),
				"data": make(map[string]string),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
