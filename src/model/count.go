package model

import (
	"fmt"
	// "strings"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Count struct {
	Id  int
	Uid int
}

func CountLoad(client *redis.Client, o orm.Ormer) error {
	var counts []Count
	num, err := o.Raw("SELECT * FROM count").QueryRows(&counts)
	if err == nil {
		var	i int64
		i = 0
		for ; i < num; i++ {
			count := counts[i]
			key := fmt.Sprintf("%s:%d:%s", "count", count.Id, "uid")
			value := fmt.Sprintf("%d", count.Id)
			err = client.Set(key, value, 0).Err()
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}

func CountSync(client *redis.Client, o orm.Ormer) error {
	// key := fmt.Sprintf("%s:%s", "count", "id")
	// client.ZRang()
	return KFail
}
