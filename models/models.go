package models

//用于开启数据库
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //注意注册
	"gogin/pkg/setting"
	"log"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {
	var (
		err                                                 error
		dbType, user, password, host, database, tableprefix string
	)
	sec, err := setting.Conf.GetSection("database")
	if err != nil {
		log.Fatalf("参数载入错误 %s", err.Error())
	}
	dbType = sec.Key("TYPE").MustString("mysql")
	user = sec.Key("USER").MustString("root")
	password = sec.Key("PASSWORD").MustString("123456")
	host = sec.Key("HOST").MustString("localhost:9090")
	database = sec.Key("NAME").MustString("blog")
	tableprefix = sec.Key("TABLE_PREFIX").MustString("blog_")
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, database))
	if err != nil {
		log.Fatalf("数据库载入失败 %s", err.Error())
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tableprefix + defaultTableName
	}
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}
