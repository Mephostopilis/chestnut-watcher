package server

import (
	"errors"
	"gamedef"
	"github.com/davyxu/cellnet"
	"mylog"
	"strings"
)

func Echo(context *Context, args interface{}) error {
	mylog.Log.Infof("sid:%d request.", ev.Ses.ID())
	msg := args.(gamedef.EchoReq)

	c := strings.Split(msg.Content, ":")
	if c[0] == "load" {
	}

	ack := gamedef.EchoAck{
		Session:   msg.Session,
		Errorcode: 1,
		Content:   msg.Content,
	}

	// 广播给所有连接
	context.peer.VisitSession(func(ses cellnet.Session) bool {

		ses.Send(&ack)

		return true
	})

	return errors.New("fail")
}
