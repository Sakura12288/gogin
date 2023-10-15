package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"time"
)

//用来读取参数信息

type App struct {
	PageSize        int
	JwtSecret       string
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxsize   int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var MysqlDatabaseSetting = &Database{}

func Setup() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("读取配置信息出错 %v", err)
	}
	err = cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("读取app配置信息出错 %v", err)
	}
	AppSetting.ImageMaxsize = AppSetting.ImageMaxsize * 1024 * 1024
	err = cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("读取server配置信息出错 %v", err)
	}
	fmt.Println(ServerSetting)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	err = cfg.Section("database").MapTo(MysqlDatabaseSetting)
	if err != nil {
		log.Fatalf("读取database配置信息出错 %v", err)
	}
	fmt.Println(MysqlDatabaseSetting)
}
