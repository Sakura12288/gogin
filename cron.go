package main

import (
	"github.com/robfig/cron/v3"
	"gogin/models"
	"gogin/pkg/logging"
	"log"
)

//定时任务

func main() {
	log.Println("starting")
	c := cron.New()
	var err error
	_, err = c.AddFunc("* * * * *", func() {
		logging.Info("系统正在删除多余的tag")
		models.CleanAllTag()
	})
	_, err = c.AddFunc("* * * * *", func() {
		logging.Info("系统正在删除多余的Articles")
		models.CleanAllArticles()
	})
	if err != nil {
		logging.Info(err)
	}
	c.Start()
	select {}
}
