package config

// import (
// 	"errors"
// 	"fmt"
// 	"github.com/astaxie/beego/orm"
// 	"model"
// )

// func ReadyForUser(o orm.Ormer) error {
// 	user := model.User{Id: 1}

// 	// insert
// 	id, err := o.Insert(&user)
// 	fmt.Printf("ID: %d, ERR: %v\n", id, err)

// 	// update
// 	user.Nickname = "astaxie"
// 	num, err := o.Update(&user)
// 	fmt.Printf("NUM: %d, ERR: %v\n", num, err)

// 	// read one
// 	u := model.User{Id: user.Id}
// 	err = o.Read(&u)
// 	fmt.Printf("ERR: %v, Nickname: %v\n", err, u.Nickname)

// 	return errors.New("text")
// }
