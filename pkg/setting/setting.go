package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

//用来读取参数信息

type App struct {
	PageSize        int
	JwtSecret       string
	RuntimeRootPath string

	PrefixUrl      string
	ImageSavePath  string
	ImageMaxsize   int
	ImageAllowExts []string

	ExportExcelSavePath string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	QrCodeSavePath string
	PosterSavePath string
	FontSavePath   string
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

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

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
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	err = cfg.Section("database").MapTo(MysqlDatabaseSetting)
	if err != nil {
		log.Fatalf("读取database配置信息出错 %v", err)
	}
	err = cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("读取Redis配置出错 %v", err)
	}
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}
