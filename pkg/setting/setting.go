package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

//用来读取参数信息

var (
	Conf    *ini.File
	RunMode string

	PageSize  int
	JwtSecret string

	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Type        string
	User        string
	Password    string
	Host        string
	Database    string
	TablePrefix string
)

func init() {
	var err error
	Conf, err = ini.Load("conf/conf.ini")
	if err != nil {
		log.Fatal("配置信息读取失败", err)
	}
	LoadBase()
	LoadApp()
	LoadServer()
	LoadDatabase()
}
func LoadBase() {
	RunMode = Conf.Section("").Key("RUNMODE").MustString("debug")
}

func LoadServer() {
	sec, err := Conf.GetSection("server")
	if err != nil {
		log.Fatalf("Parse Server err : %s", err.Error())
	}
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
	JwtSecret = sec.Key("JWTSECRET").MustString("sibensb")
}
func LoadApp() {
	sec, err := Conf.GetSection("app")
	if err != nil {
		log.Fatalf("Parse app err : %s", err.Error())
	}
	HttpPort = sec.Key("HTTP_PORT").MustInt(9090)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadDatabase() {
	sec, err := Conf.GetSection("database")
	if err != nil {
		log.Fatalf("Parse database err : %s", err.Error())
	}
	Type = sec.Key("TYPE").MustString("mysql")
	User = sec.Key("USER").MustString("root")
	Password = sec.Key("PASSWORD").MustString("123456")
	Host = sec.Key("HOST").MustString("localhost:9090")
	Database = sec.Key("NAME").MustString("blog")
	TablePrefix = sec.Key("TABLE_PREFIX").MustString("blog_")
}
