package config

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"model"
)

func ReadyForUser(o orm.Ormer) error {
	user := model.User{
		Id:       2,
		OpenId:   1,
		NameId:   3,
		Nickname: "hello",
		Sex:      1,
		Province: "Sc",
	}
	cnt, err := o.QueryTable("User").Filter("Id__contains", 2).Count()
	if cnt > 0 {
		return errors.New("hell")
	}

	// insert
	id, err := o.Insert(&user)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	// update
	user.Nickname = "astaxie"
	num, err := o.Update(&user)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)

	// read one
	u := model.User{Id: user.Id}
	err = o.Read(&u)
	fmt.Printf("ERR: %v, Nickname: %v\n", err, u.Nickname)

	return errors.New("hello")
}
