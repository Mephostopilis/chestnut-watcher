package model

import (
	// "fmt"
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"mylog"
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
	HeadImg  string `orm:"size(1024)"`
}

func UserLoad(client *redis.Client, o orm.Ormer) error {
	cnt, err := o.QueryTable("User").Count()
	if err != nil {
		mylog.Log.Errorln(err)
	}
	if cnt <= 0 {
		return common.ErrOK
	}

	u := User{Id: 2}
	err = o.Read(&u)
	mylog.Log.Infoln("%s", err)
	value := fmt.Sprintf("Id:%d,OpenId:%d", u.Id, u.OpenId)
	client.Set("User:2", value, 0)
	client.LPush("User", 2)
	value = client.Get("User:2").Val()
	mylog.Log.Infof("user:2 of value is: %s", value)
	return KFail
}

func UserSync(client *redis.Client, o orm.Ormer) error {

	return KFail
}
