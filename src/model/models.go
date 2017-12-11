package model

import (
	"github.com/astaxie/beego/orm"
)

type Tag struct {
	Id   int
	Name string
}

func init() {
	orm.RegisterModel(new(Tag))
}
