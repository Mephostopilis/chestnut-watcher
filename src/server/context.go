package server

import (
	"errors"
	"gamedef"
	"github.com/astaxie/beego/orm"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"mylog"
	"syncdb"
)

type Context struct {
	Log     *golog.Logger
	queue   cellnet.EventQueue
	client  *redis.Client
	o       orm.Ormer
	peer    cellnet.Peer
	request map[string]func(args interface{})
}

func NewContext() *Context {

	log := golog.New("server")

	client := syncdb.NewRedisClient()
	o := syncdb.NewOrm()

	queue := cellnet.NewEventQueue()

	peer := socket.NewAcceptor(queue).Start("127.0.0.1:8801")
	peer.SetName("client")

	context := &Context{
		Log:     log,
		client:  client,
		o:       o,
		queue:   queue,
		peer:    peer,
		request: make(map[string]func(args interface{})),
	}
	return context
}

func (context *Context) Queue() cellnet.EventQueue {
	return context.queue
}

func (context *Context) Startup() error {
	// log.Log.Infoln("Startup")

	context.Log.Infof("%d", mylog.I)
	cellnet.RegisterMessage(context.peer, "gamedef.EchoReq", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.EchoReq)
		Echo(context, msg)
	})

	syncdb.Load(context.client, context.o)
	return errors.New("h")
}

func (context *Context) Update() error {
	syncdb.Sync(context.client, context.o)
	context.Log.Infof("%d", mylog.I)
	mylog.Log.Infof("%d", mylog.I)
	return errors.New("h")
}

func (context *Context) Cleanup() error {
	context.peer.Stop()
	return errors.New("h")
}
