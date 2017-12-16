package model

import (
	"fmt"
)

type OpenId struct {
	OpenId string `orm:"size(100)"`
	Uid    int
}

func load(client *redis.Client, o orm.orm) error {
	var counts []Count
	num, err := o.Raw("SELECT * FROM count").QueryRows(&counts)
	if err == nil {
		for i := 0; i < num; i++ {
			count := counts[i]
			key := fmt.Sprintf("%s:%d:%s", "count", count.Id, "uid")
			value := fmt.Sprintf("%d", count.Id)
			client.Set(key, value)
		}
	}
	return err
}

func sync(client *redis.Client, o orm.orm) error {
	key := fmt.Sprintf("%s:%s", "count", "id")
	client.ZRang()
}
