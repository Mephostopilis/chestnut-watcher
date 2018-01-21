package model

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type OpenId struct {
	OpenId string `orm:"pk;size(100)"`
	Uid    int
}

func OpenIdLoad(client *redis.Client, o orm.Ormer) error {
	var counts []Count
	num, err := o.Raw("SELECT * FROM count").QueryRows(&counts)
	if err == nil {
		var i int64 = 0
		for ; i < num; i++ {
			count := counts[i]
			key := fmt.Sprintf("%s:%d:%s", "count", count.Id, "uid")
			value := fmt.Sprintf("%d", count.Id)
			client.Set(key, value, 0)
		}
	}
	return err
}

func OpenIdSync(client *redis.Client, o orm.Ormer) error {
	// key := fmt.Sprintf("%s:%s", "count", "id")
	// client.ZRang()
	return KFail
}
