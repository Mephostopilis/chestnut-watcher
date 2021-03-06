package syncdb

import (
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"model"
	"mylog"
)

// var client *redis.Client

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/default?charset=utf8")
	orm.RegisterModel(new(model.Count))
	orm.RegisterModel(new(model.NameId))
	orm.RegisterModel(new(model.OpenId))
	orm.RegisterModel(new(model.User))
	orm.RunSyncdb("default", false, true)
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	mylog.Log.Infoln(pong, err)
	// fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func NewOrm() orm.Ormer {
	return orm.NewOrm()
}

func Load(client *redis.Client, o orm.Ormer) {
	model.CountLoad(client, o)
	model.NameIdLoad(client, o)
	model.OpenIdLoad(client, o)
	model.UserLoad(client, o)
}

func LoadByUid(client *redis.Client, o orm.Ormer, uid int) {

}

func Sync(client *redis.Client, o orm.Ormer) {

}
