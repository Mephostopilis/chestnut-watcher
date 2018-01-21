package gamedef

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/json"
	// "github.com/davyxu/cellnet/util"
	"github.com/davyxu/goobjfmt"
	"reflect"
)

type EchoReq struct {
	Session int
	Content string
}

func (m *EchoReq) String() string {
	return goobjfmt.CompactTextString(m)
}

type EchoAck struct {
	Session int
	Errorcode int
	Content string
}

func (m *EchoAck) String() string {
	return goobjfmt.CompactTextString(m)	
}

func init() {

	// coredef.proto
	cellnet.RegisterMessageMeta("json", "gamedef.EchoReq", reflect.TypeOf((*EchoReq)(nil)).Elem(), 1)
	cellnet.RegisterMessageMeta("json", "gamedef.EchoAck", reflect.TypeOf((*EchoAck)(nil)).Elem(), 2)
}
