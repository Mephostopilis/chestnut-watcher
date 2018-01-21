package model

import (
	// "fmt"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int64
	OpenId   int
	NameId   int
	Nickname string `orm:"size(32)"`
	Sex      int
	Province string `orm:"size(32)"`
	City     string `orm:"size(32)"`
	Country  string `orm:"size(32)"`
	Gold     int
	Diamond  int
	RCard    int
	HeadImg  string
}

func UserLoad(client *redis.Client, o orm.Ormer) error {
	return KFail
}


func UserSync(client *redis.Client, o orm.Ormer) error {
	return KFail
}