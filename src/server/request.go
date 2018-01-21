package server

import (
	// "log"
	"gamedef"
	"errors"
	"github.com/davyxu/cellnet"
)

func Echo(context *Context, args interface{}) error {
	// log := log.Log
	// log.Infof("sid:%d", ev.Ses.ID())
	
	msg := args.(gamedef.EchoReq)
	ack := gamedef.EchoAck{
		Session: msg.Session,
		Errorcode: 1,
		Content: msg.Content,
	}

	// 广播给所有连接
	context.peer.VisitSession(func(ses cellnet.Session) bool {

		ses.Send(&ack)

		return true
	})

	return errors.New("fail")
}
