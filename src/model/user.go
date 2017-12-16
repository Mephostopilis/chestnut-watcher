package model

import (
	"fmt"
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
