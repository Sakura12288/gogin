package app

import (
	"github.com/astaxie/beego/validation"
	"gogin/pkg/logging"
)

func MakeErrors(errors []*validation.Error) {
	for _, v := range errors {
		logging.Info(v.Key, v.Message)
	}
	return
}
