package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"gogin/models"
	"gogin/pkg/gredis"
	"gogin/pkg/logging"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/routers"
	"log"
	"syscall"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	util.Setup()
	gredis.Setup()
}
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	//routersInit := routers.InitRouter()
	//readTimeout := setting.ServerSetting.ReadTimeout
	//writeTimeout := setting.ServerSetting.WriteTimeout
	//endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttPPort)
	//maxHeaderBytes := 1 << 20
	//
	//server := &http.Server{
	//	Addr:           endPoint,
	//	Handler:        routersInit,
	//	ReadTimeout:    readTimeout,
	//	WriteTimeout:   writeTimeout,
	//	MaxHeaderBytes: maxHeaderBytes,
	//}
	//
	//log.Printf("[info] start http server listening %s", endPoint)
	//
	//server.ListenAndServe()
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20

	endpoint := fmt.Sprintf(":%d", setting.ServerSetting.HttPPort)

	server := endless.NewServer(endpoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server err : %v", err)
	}
}
